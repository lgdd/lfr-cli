package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var (
	Config              = make(map[string]string)
	ConfigTreePathsEnvs = []string{"LIFERAY_ROUTES_DXP", "LIFERAY_ROUTES_CLIENT_EXTENSION"}
)

func getConfigTreePaths() []string {
	var configTreePaths []string

	for _, env := range ConfigTreePathsEnvs {
		envPath, envPathExists := os.LookupEnv(env)
		if envPathExists {
			slog.Info(fmt.Sprintf("%s: %s", env, envPath))
			configTreePaths = append(configTreePaths, envPath)
		}
	}

	if len(configTreePaths) == 0 {
		slog.Warn(fmt.Sprintf("No environment variable found for config %s", ConfigTreePathsEnvs))
		slog.Warn("Default config path to './dxp-metadata'")
		configTreePaths = append(configTreePaths, "dxp-metadata")
	}

	return configTreePaths
}

func initConfig() error {
	configTreePaths := getConfigTreePaths()

	slog.Info("Loading config:")

	for _, configTreePath := range configTreePaths {
		err := filepath.Walk(configTreePath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				fileContentBytes, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				slog.Info(fmt.Sprintf("- %s=%s", info.Name(), string(fileContentBytes)))
				Config[info.Name()] = string(fileContentBytes)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

type JWTClaims struct {
	Subject    string    `json:"sub"`
	Issuer     string    `json:"iss"`
	ClientId   string    `json:"client_id"`
	Audience   []string  `json:"aud"`
	GrantType  string    `json:"grant_type"`
	Scope      string    `json:"scope"`
	Expiration time.Time `json:"exp"`
	IssuedAt   time.Time `json:"iat"`
	ID         string    `json:"jti"`
	Username   string    `json:"username"`
}

func logDecodedToken(token jwt.Token) {
	var clientId string
	var grantType string
	var scope string
	var username string

	jti, _ := token.JwtID()
	sub, _ := token.Subject()
	iss, _ := token.Issuer()
	aud, _ := token.Audience()
	iat, _ := token.IssuedAt()
	exp, _ := token.Expiration()
	_ = token.Get("scope", &scope)
	_ = token.Get("username", &username)
	_ = token.Get("grant_type", &grantType)
	_ = token.Get("client_id", &clientId)

	claims := &JWTClaims{
		Subject:    sub,
		Issuer:     iss,
		ClientId:   clientId,
		Audience:   aud,
		GrantType:  grantType,
		Scope:      scope,
		Expiration: exp,
		IssuedAt:   iat,
		ID:         jti,
		Username:   username,
	}

	claimsJson, _ := json.Marshal(claims)

	slog.Info(fmt.Sprintf("JWT Claims: %s", string(claimsJson)))
	slog.Info(fmt.Sprintf("JWT ID: %s", jti))
	slog.Info(fmt.Sprintf("JWT Subject: %s", sub))
}

func validateJWT(tokenString string) (jwt.Token, error) {
	var oauth2JWKSURIBuilder strings.Builder

	protocol := Config["com.liferay.lxc.dxp.server.protocol"]
	host := Config["com.liferay.lxc.dxp.mainDomain"]

	oauth2JWKSURIBuilder.WriteString(protocol)
	oauth2JWKSURIBuilder.WriteString("://")
	oauth2JWKSURIBuilder.WriteString(host)
	oauth2JWKSURIBuilder.WriteString("/o/oauth2/jwks")

	response, err := http.Get(oauth2JWKSURIBuilder.String())

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	jwks, err := jwk.Parse(body)

	if err != nil {
		return nil, err
	}

	key, found := jwks.Key(0)

	if found {
		alg, found := key.Algorithm()

		if found {
			token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(alg, key), jwt.WithValidate(true))

			if err != nil {
				return nil, err
			}

			logDecodedToken(token)

			return token, nil

		}
	}

	return nil, errors.New("no json web key found")
}

func validateClientId(token jwt.Token) error {
	var oauth2ApplicationURIBuilder strings.Builder

	protocol := Config["com.liferay.lxc.dxp.server.protocol"]
	host := Config["com.liferay.lxc.dxp.mainDomain"]
	externalReferenceCodes := Config["liferay.oauth.application.external.reference.codes"]
	externalReferenceCode := strings.Split(externalReferenceCodes, ",")[0]

	oauth2ApplicationURIBuilder.WriteString(protocol)
	oauth2ApplicationURIBuilder.WriteString("://")
	oauth2ApplicationURIBuilder.WriteString(host)
	oauth2ApplicationURIBuilder.WriteString("/o/oauth2/application?externalReferenceCode=")
	oauth2ApplicationURIBuilder.WriteString(externalReferenceCode)

	var clientId string
	err := token.Get("client_id", &clientId)

	if err != nil {
		return err
	}

	response, err := http.Get(oauth2ApplicationURIBuilder.String())

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	jsonResponse := make(map[string]string)

	err = json.Unmarshal(body, &jsonResponse)

	if err != nil {
		return err
	}

	if clientId == jsonResponse["client_id"] {
		return nil
	}

	return errors.New("client id from token and oauth application don't match")
}

func jwtHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		token, err := validateJWT(tokenString)

		if err != nil {
			log.Fatal(err)
		}

		err = validateClientId(token)

		if err != nil {
			log.Fatal(err)
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	err := initConfig()

	if err != nil {
		log.Fatal(err)
	}

	homeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseBuilder strings.Builder

		responseBuilder.WriteString("Endpoints available:\n\n")
		responseBuilder.WriteString("- /ready\n")
		responseBuilder.WriteString("- /object/action/1\n")

		_, err := fmt.Fprintf(w, responseBuilder.String())

		if err != nil {
			log.Fatalf("/ready failed with error: %s", err.Error())
		}
	})

	readyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "ready")

		if err != nil {
			log.Fatalf("/ready failed with error: %s", err.Error())
		}
	})

	objectAction1Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("execute /object/action/1")

		var objectEntry interface{}
		err := json.NewDecoder(r.Body).Decode(&objectEntry)

		if err != nil {
			log.Fatal(err)
		}

		objectEntryJsonBytes, err := json.MarshalIndent(objectEntry, "", "\t")

		if err != nil {
			log.Fatal(err)
		}

		objectEntryJsonString := string(objectEntryJsonBytes)
		slog.Info(objectEntryJsonString)
	})

	http.Handle("/", homeHandler)
	http.Handle("/ready", readyHandler)
	http.Handle("/object/action/1", jwtHandler(objectAction1Handler))

	slog.Info("Server started at http://localhost:8126")

	err = http.ListenAndServe(":8126", nil)

	if err != nil {
		log.Fatal(err)
	}
}
