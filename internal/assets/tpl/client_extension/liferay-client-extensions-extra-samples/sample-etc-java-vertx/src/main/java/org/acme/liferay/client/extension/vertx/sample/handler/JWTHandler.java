package org.acme.liferay.client.extension.vertx.sample.handler;

import io.vertx.core.Vertx;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.jwt.JWTAuth;
import io.vertx.ext.auth.jwt.JWTAuthOptions;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.client.WebClient;
import java.util.List;
import java.util.Map;

public class JWTHandler extends BaseHandler {

  public JWTHandler(Vertx vertx, WebClient webClient, Map<String, String> configMap) {
    super(vertx, webClient, configMap);
  }

  @Override
  public void handle(RoutingContext ctx) {
    final String token = ctx.request().headers().get("Authorization").split("Bearer ")[1];
    final String externalReferenceCode =
        _configMap
            .get("liferay.oauth.application.external.reference.codes")
            .split(",")[0];
    final String liferayDomain =
        _configMap
            .getOrDefault("com.liferay.lxc.dxp.mainDomain", "localhost:8080");
    final String liferayHost = liferayDomain.split(":")[0];
    final int liferayPort = Integer.parseInt(liferayDomain.split(":")[1]);

    _webClient
        .get(liferayPort, liferayHost, "/o/oauth2/jwks")
        .send()
        .onSuccess(response -> {
          JsonObject jwksResponse = response.bodyAsJsonObject();
          List<JsonObject> jwks =
              jwksResponse
                  .getJsonArray("keys")
                  .stream()
                  .map(JsonObject.class::cast)
                  .toList();

          JWTAuthOptions jwtOptions = new JWTAuthOptions();
          jwtOptions.setJwks(jwks);
          JWTAuth jwtAuth = JWTAuth.create(_vertx, jwtOptions);

          jwtAuth.authenticate(
              new JsonObject().put("token", token)
          ).onSuccess(user -> {
            final JsonObject decodedToken = user.principal();
            final String tokenClientId = decodedToken.getString("client_id");
            final String oauth2ApplicationURI =
                "/o/oauth2/application?externalReferenceCode=" + externalReferenceCode;

            _log.info("JWT Claims: " + decodedToken);
            _log.info("JWT ID: " + decodedToken.getString("jti"));
            _log.info("JWT Subject: " + decodedToken.getString("sub"));

            validateClientId(ctx, tokenClientId);

          }).onFailure(error -> {
            _log.error(error.getMessage());
            ctx.fail(error);
          });

        })
        .onFailure(error -> {
          _log.error(error.getMessage());
          ctx.fail(error);
        });
  }

  private void validateClientId(RoutingContext ctx, String tokenClientId) {
    final String externalReferenceCode =
        _configMap
            .get("liferay.oauth.application.external.reference.codes")
            .split(",")[0];
    final String liferayDomain =
        _configMap
            .getOrDefault("com.liferay.lxc.dxp.mainDomain", "localhost:8080");
    final String liferayHost = liferayDomain.split(":")[0];
    final int liferayPort = Integer.parseInt(liferayDomain.split(":")[1]);
    final String oauth2ApplicationURI =
        "/o/oauth2/application?externalReferenceCode=" + externalReferenceCode;

    _webClient
        .get(liferayPort, liferayHost, oauth2ApplicationURI)
        .send()
        .onSuccess(res -> {
          JsonObject oauth2Application = res.bodyAsJsonObject();
          if (tokenClientId.equals(oauth2Application.getString("client_id"))) {
            ctx.next();
          } else {
            String message = "Client id from token and oauth application matched";
            _log.error(message);
            ctx.fail(new Exception(message));
          }
        })
        .onFailure(error -> {
          _log.error(error.getMessage());
          ctx.fail(error);
        });
  }

  private static final Logger _log = LoggerFactory.getLogger(JWTHandler.class);
}
