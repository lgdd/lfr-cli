package org.acme.liferay.client.extension.quarkus.sample;

import com.auth0.jwk.Jwk;
import com.auth0.jwk.JwkException;
import com.auth0.jwk.JwkProvider;
import com.auth0.jwk.UrlJwkProvider;
import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.exceptions.JWTVerificationException;
import com.auth0.jwt.interfaces.DecodedJWT;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.container.ContainerRequestFilter;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.Provider;
import java.io.IOException;
import java.net.URI;
import java.security.interfaces.RSAPublicKey;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Provider
public class JWTValidationFilter implements ContainerRequestFilter {

  @Override
  public void filter(ContainerRequestContext containerRequestContext) throws IOException {
    String authorization = containerRequestContext.getHeaders().getFirst("Authorization");
    if (authorization != null && authorization.startsWith("Bearer ")) {
      try {
        String token = authorization.split("Bearer ")[1];
        DecodedJWT jwt = JWT.decode(token);

        JwkProvider jwkProvider = new UrlJwkProvider(
            URI.create(liferayProtocol + "://" + liferayDomain + "/o/oauth2/jwks").toURL()
        );

        Jwk jwk = jwkProvider.get(jwt.getKeyId());

        Algorithm algorithm = Algorithm.RSA256((RSAPublicKey) jwk.getPublicKey(), null);

        JWT.require(algorithm).build().verify(token);

        if (_log.isInfoEnabled()) {
          _log.info("JWT Claims: {}", jwt.getClaims());
          _log.info("JWT ID: {}", jwt.getId());
          _log.info("JWT Subject: {}", jwt.getSubject());
        }

      } catch (JwkException | JWTVerificationException e) {
        if(_log.isErrorEnabled()) {
          _log.error(e.getMessage(), e);
        }
        containerRequestContext.abortWith(Response.status(Response.Status.UNAUTHORIZED).build());
      }
    } else {
      if(_log.isErrorEnabled()) {
        if(authorization == null) {
          _log.error("Authorization header is missing");
        } else {
          _log.error("Bearer JWT is missing");
        }
      }
      containerRequestContext.abortWith(Response.status(Response.Status.UNAUTHORIZED).build());
    }
  }

  @ConfigProperty(name = "liferay.oauth.application.external.reference.codes")
  String externalReferenceCodes;

  @ConfigProperty(name = "com.liferay.lxc.dxp.mainDomain")
  String liferayDomain;

  @ConfigProperty(name = "com.liferay.lxc.dxp.server.protocol")
  String liferayProtocol;

  private static final Logger _log = LoggerFactory.getLogger(JWTValidationFilter.class);

}
