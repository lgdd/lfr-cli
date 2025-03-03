package org.acme.liferay.client.extension.vertx.sample;

import io.vertx.config.ConfigRetriever;
import io.vertx.config.ConfigRetrieverOptions;
import io.vertx.config.ConfigStoreOptions;
import io.vertx.core.DeploymentOptions;
import io.vertx.core.Launcher;
import io.vertx.core.Vertx;
import io.vertx.core.json.JsonObject;

public class VertxApplication extends Launcher {

  public static void main(String[] args) {
    Vertx vertx = Vertx.vertx();

    ConfigStoreOptions properties = new ConfigStoreOptions()
        .setType("file")
        .setFormat("properties")
        .setConfig(
            new JsonObject()
                .put("path", "application.properties")
        );

    ConfigStoreOptions env = new ConfigStoreOptions()
        .setType("env");

    ConfigRetrieverOptions options = new ConfigRetrieverOptions()
        .addStore(properties)
        .addStore(env);

    ConfigRetriever retriever = ConfigRetriever.create(vertx, options);

    retriever.getConfig().onComplete(json -> {
      vertx.deployVerticle(new MainVerticle(), new DeploymentOptions().setConfig(json.result()));
    });
  }

}
