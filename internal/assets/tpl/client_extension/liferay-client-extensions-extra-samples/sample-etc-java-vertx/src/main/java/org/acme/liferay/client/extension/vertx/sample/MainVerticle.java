package org.acme.liferay.client.extension.vertx.sample;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.Promise;
import io.vertx.core.http.HttpServer;
import io.vertx.core.impl.logging.Logger;
import io.vertx.core.impl.logging.LoggerFactory;
import io.vertx.core.json.JsonArray;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.client.WebClient;
import io.vertx.ext.web.handler.BodyHandler;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;
import org.acme.liferay.client.extension.vertx.sample.handler.JWTHandler;
import org.acme.liferay.client.extension.vertx.sample.handler.ObjectAction1Handler;

public class MainVerticle extends AbstractVerticle {

  @Override
  public void start(Promise<Void> startPromise) throws Exception {
    HttpServer server = vertx.createHttpServer();
    WebClient webClient = WebClient.create(vertx);
    Router router = Router.router(vertx);

    _configMap = getConfigMap();

    JWTHandler jwtHandler = new JWTHandler(vertx, webClient, _configMap);
    ObjectAction1Handler objectAction1Handler = new ObjectAction1Handler(vertx, webClient, _configMap);

    router.route().handler(BodyHandler.create());

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
        .handler(objectAction1Handler);

    int serverPort = config().getInteger("server.port", 8082);

    server
        .requestHandler(router)
        .listen(serverPort)
        .onComplete(http -> {
          if (http.succeeded()) {
            startPromise.complete();
            _log.info("HTTP server started on port " + serverPort);
          } else {
            startPromise.fail(http.cause());
          }
        });
  }

  private Map<String, String> getConfigMap() {
    final Map<String, String> configMap = new HashMap<>();
    final List<String> configTreePaths = new ArrayList<>();

    _log.info("Start loading config...");

    String liferayRoutesDXP = config().getString("LIFERAY_ROUTES_DXP", null);
    String liferayRoutesClientExtension = config().getString("LIFERAY_ROUTES_CLIENT_EXTENSION",
        null);

    if (liferayRoutesDXP != null) {
      configTreePaths.add(liferayRoutesDXP);
    }

    if (liferayRoutesClientExtension != null) {
      configTreePaths.add(liferayRoutesClientExtension);
    }

    config().stream()
        .filter(keyValue -> keyValue.getKey().contains("liferay"))
        .forEach(keyValue -> {
          configMap.put(keyValue.getKey(), (String) keyValue.getValue());
        });

    for (String configTreePath : configTreePaths) {
      Path start = Paths.get(configTreePath);
      try (Stream<Path> walk = Files.walk(start)) {
        walk
            .filter(Files::isRegularFile)
            .forEach(path -> {
              try {
                String content = Files.readString(path, StandardCharsets.UTF_8);
                configMap.put(path.getFileName().toString(), content);
              } catch (IOException e) {
                _log.error(e);
              }
            });
      } catch (IOException e) {
        _log.error(e);
      }
    }

    configMap.forEach((key, value) -> {
      _log.info(key + ": " + value);
    });

    _log.info("Loading config done!");

    return configMap;
  }


  private static Map<String, String> _configMap;

  private static final Logger _log = LoggerFactory.getLogger(MainVerticle.class);

}
