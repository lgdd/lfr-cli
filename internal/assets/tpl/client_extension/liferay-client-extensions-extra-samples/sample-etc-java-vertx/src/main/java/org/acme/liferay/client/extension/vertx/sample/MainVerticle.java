package org.acme.liferay.client.extension.vertx.sample;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.Handler;
import io.vertx.core.Promise;
import io.vertx.core.http.HttpServer;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.jwt.JWTAuth;
import io.vertx.ext.auth.jwt.JWTAuthOptions;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.client.WebClient;
import io.vertx.ext.web.handler.BodyHandler;
import java.util.List;

public class MainVerticle extends AbstractVerticle {

  @Override
  public void start(Promise<Void> startPromise) throws Exception {
    HttpServer server = vertx.createHttpServer();
    WebClient webClient = WebClient.create(vertx);
    Router router = Router.router(vertx);

    router.route().handler(BodyHandler.create());

    Handler<RoutingContext> jwtHandler = ctx -> {
      validateJWT(ctx, webClient);
    };

    router
        .get("/")
        .handler(ctx ->
            ctx.response()
                .putHeader("content-type", "application/json")
                .end(
                    new JsonArray()
                        .add("/ready")
                        .add("/object/action/1")
                        .encode()
                )
        );

    router
        .get("/ready")
        .handler(ctx ->
            ctx.response()
                .putHeader("content-type", "text/plain")
                .end("ready")
        );

    router
        .post("/object/action/1")
        .consumes("application/json")
        .handler(jwtHandler)
        .handler(ctx -> {
          _log.info("execute object action 1");
          _log.info(ctx.body().asString());
          ctx.response().putHeader("content-type", "text/plain").end("executed object action 1");
        });

    server
        .requestHandler(router)
        .listen(8082)
        .onComplete(http -> {
          if (http.succeeded()) {
            startPromise.complete();
            _log.info("HTTP server started on port 8082");
          } else {
            startPromise.fail(http.cause());
          }
        });
  }

  private void validateJWT(RoutingContext ctx, WebClient webClient) {
    final String token = ctx.request().headers().get("Authorization").split("Bearer ")[1];
    final String externalReferenceCode =
        config()
            .getString("liferay.oauth.application.external.reference.codes")
            .split(",")[0];
    final String liferayDomain =
        config()
            .getString("com.liferay.lxc.dxp.mainDomain", "localhost:8080");
    final String liferayHost = liferayDomain.split(":")[0];
    final int liferayPort = Integer.parseInt(liferayDomain.split(":")[1]);

    webClient
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
          JWTAuth jwtAuth = JWTAuth.create(vertx, jwtOptions);

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

            validateClientId(ctx, webClient, tokenClientId);

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

  private void validateClientId(RoutingContext ctx, WebClient webClient, String tokenClientId) {
    final String externalReferenceCode =
        config()
            .getString("liferay.oauth.application.external.reference.codes")
            .split(",")[0];
    final String liferayDomain =
        config()
            .getString("com.liferay.lxc.dxp.mainDomain", "localhost:8080");
    final String liferayHost = liferayDomain.split(":")[0];
    final int liferayPort = Integer.parseInt(liferayDomain.split(":")[1]);
    final String oauth2ApplicationURI =
        "/o/oauth2/application?externalReferenceCode=" + externalReferenceCode;

    webClient
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

  private static final Logger _log = LoggerFactory.getLogger(MainVerticle.class);

}
