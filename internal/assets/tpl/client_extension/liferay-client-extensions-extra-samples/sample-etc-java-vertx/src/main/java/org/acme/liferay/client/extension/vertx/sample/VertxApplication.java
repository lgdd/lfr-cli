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

    ConfigRetriever retriever = ConfigRetriever.create(vertx, new ConfigRetrieverOptions()
        .addStore(new ConfigStoreOptions().setType("file").setConfig(new JsonObject().put("path", "conf/dxp-metadata.json"))));

    retriever.getConfig().onComplete(json -> {
      vertx.deployVerticle(new MainVerticle(), new DeploymentOptions().setConfig(json.result()));
    });
  }

}
