import {serve} from 'bun';
import {log} from "./log.ts";
import {type JSONWebKeySet, type JWK, type JWTPayload, jwtVerify} from "jose";
import {getConfigMap} from "./config.ts";

const configMap = await getConfigMap()
const externalReferenceCode = configMap.get("liferay.oauth.application.external.reference.codes")?.split(",")[0]
const liferayHost = `${configMap.get("com.liferay.lxc.dxp.server.protocol")}://${configMap.get("com.liferay.lxc.dxp.mainDomain")}`

const AuthError = (message: string) => {
  const error = new Error(message)
  error.name = "AuthError"
  return error;
}

const getJWT = async (req: Request) => {
  const authorization = req.headers.get("Authorization")

  if (authorization == null) {
    throw AuthError("Authorization header is missing")
  }

  if (!authorization.startsWith("Bearer")) {
    throw AuthError("Bearer token is missing")
  }

  return authorization.slice("Bearer ".length)
}

const validateJWT = async (jwt: string) => {
  const response = await fetch(`${liferayHost}/o/oauth2/jwks`)
  const jwks: JSONWebKeySet = await response.json()
  const jwk: JWK = jwks.keys[0]
  const {payload} = await jwtVerify(jwt, jwk)
  log.info(`Decoded JWT: ${JSON.stringify(payload)}`)
  return payload
}

const validateClientId = async (req: Request, decodedToken: JWTPayload) => {
  const response = await fetch(`${liferayHost}/o/oauth2/application?externalReferenceCode=${externalReferenceCode}`)
  const jsonResponse = await response.json()
  const clientId = jsonResponse["client_id"]
  if (clientId !== decodedToken["client_id"]) {
    throw AuthError("Client id from token and OAuth application don't match")
  }
}

const objectAction1 = async (req: Request) => {
  const jwt = await getJWT(req)
  const decodedToken = await validateJWT(jwt)
  await validateClientId(req, decodedToken)

  const data = await req.json()
  log.info(JSON.stringify(data))

  const authorUserId = data.objectEntry.userId

  const authorUserInfoURL = `${liferayHost}/o/headless-admin-user/v1.0/user-accounts/${authorUserId}`
  const headers = new Headers();
  headers.set('Authorization', `Bearer ${jwt}`);
  headers.set('Content-Type', 'application/json');

  const requestOptions = {
    method: 'GET',
    headers: headers
  };

  log.info(`Fetching author user information at ${authorUserInfoURL}`)
  const userInfoResponse = await fetch(`${authorUserInfoURL}`, requestOptions)

  if (userInfoResponse.status / 100 != 2) {
    const errorMessage = `Could not fetch author user information: ${userInfoResponse.status} error`
    log.error(errorMessage)
    return new Response('', {
      status: 500,
      statusText: errorMessage
    })
  }

  const authorUserInfo = await userInfoResponse.json()

  log.info(JSON.stringify(authorUserInfo));

  return new Response('', {
    status: 202
  })
}

const home = async () => {
  return new Response(
      'Endpoints available are:\n' +
      '- /ready\n' +
      '- /object/action/1\n'
  )
}

const ready = async () => {
  return new Response('ready')
}

log.info("Listening to http://localhost:8228")

serve({
  port: 8228,
  routes: {
    "/": home,
    "/ready": ready,
    "/object/action/1": objectAction1,
  },
  error(error) {
    log.error(error.message);
    if (error.name == "AuthError") {
      return new Response(`${error.message}`, {
        status: 401
      })
    }
    return new Response(`Internal Error: ${error.message}`, {
      status: 500,
      headers: {
        "Content-Type": "text/plain",
      },
    });
  },
})