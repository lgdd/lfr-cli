package org.acme.liferay.client.extension.vertx.sample.handler;

import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.client.WebClient;
import java.util.Map;

public abstract class BaseHandler implements Handler<RoutingContext> {

  public BaseHandler(Vertx vertx, WebClient webClient, Map<String, String> configMap) {
    _vertx = vertx;
    _webClient = webClient;
    _configMap = configMap;
  }

  protected final Map<String, String> _configMap;
  protected final WebClient _webClient;
  protected final Vertx _vertx;
}
