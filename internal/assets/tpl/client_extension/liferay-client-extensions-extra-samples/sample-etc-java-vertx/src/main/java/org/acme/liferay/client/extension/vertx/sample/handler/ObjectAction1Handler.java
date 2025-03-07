package org.acme.liferay.client.extension.vertx.sample.handler;

import io.vertx.core.Vertx;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.client.WebClient;
import java.util.Map;

public class ObjectAction1Handler extends BaseHandler {

  public ObjectAction1Handler(Vertx vertx, WebClient webClient, Map<String, String> configMap) {
    super(vertx, webClient, configMap);
  }

  @Override
  public void handle(RoutingContext ctx) {
    final String token = ctx.request().headers().get("Authorization");
    final String liferayDomain =
        _configMap
            .getOrDefault("com.liferay.lxc.dxp.mainDomain", "localhost:8080");
    final String liferayHost = liferayDomain.split(":")[0];
    final int liferayPort = Integer.parseInt(liferayDomain.split(":")[1]);

    JsonObject data = ctx.body().asJsonObject();

    _log.info("Object Entry: " + data.toString());

    int authorUserId = data.getJsonObject("objectEntry").getInteger("userId");
    String authorUserInfoURI = "/o/headless-admin-user/v1.0/user-accounts/" + authorUserId;

    _webClient.get(liferayPort, liferayHost, authorUserInfoURI)
        .putHeader("Authorization", token)
        .putHeader("Content-Type", "application/json")
        .send()
        .onSuccess(response -> {
          if (_log.isInfoEnabled()) {
            _log.info(response.bodyAsString());
          }
        }).onFailure(error -> {
          _log.error(error.getMessage());
          ctx.fail(error);
        });

    ctx.response().putHeader("content-type", "text/plain").end("executed object action 1");
  }

  private static final Logger _log = LoggerFactory.getLogger(ObjectAction1Handler.class);
}
