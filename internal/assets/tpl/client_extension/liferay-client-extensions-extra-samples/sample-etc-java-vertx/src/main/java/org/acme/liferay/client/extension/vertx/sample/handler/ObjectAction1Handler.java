package org.acme.liferay.client.extension.vertx.sample.handler;

import io.vertx.core.Handler;
import io.vertx.ext.web.RoutingContext;
import java.util.logging.Logger;

public class ObjectAction1Handler implements Handler<RoutingContext> {

  @Override
  public void handle(RoutingContext ctx) {
    _log.info("execute object action 1");
    _log.info(ctx.body().asString());
    ctx.response().putHeader("content-type", "text/plain").end("executed object action 1");
  }

  private static final Logger _log = Logger.getLogger(ObjectAction1Handler.class.getName());
}
