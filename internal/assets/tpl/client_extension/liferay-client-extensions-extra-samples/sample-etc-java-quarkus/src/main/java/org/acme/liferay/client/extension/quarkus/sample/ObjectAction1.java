package org.acme.liferay.client.extension.quarkus.sample;

import jakarta.json.Json;
import jakarta.json.JsonObject;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.client.Client;
import jakarta.ws.rs.client.ClientBuilder;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.core.Response.Status;
import java.net.URI;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.RestHeader;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Path("/object/action/1")
public class ObjectAction1 {

  @POST
  public Response objectAction1(@RestHeader("Authorization") String authorization,
      JsonObject objectEntry) {

    URI authorUserInfoURI = URI.create(
        liferayProtocol + "://" + liferayDomain
            + "/o/headless-admin-user/v1.0/user-accounts/"
            + objectEntry.getJsonObject("objectEntry").getInt("userId")
    );

    if (_log.isInfoEnabled()) {
      _log.info("Object Entry: {}", objectEntry.toString());
      _log.info("Fetching author user info from {}", authorUserInfoURI);
    }

    try (Client client = ClientBuilder.newClient()) {
      JsonObject response = client.target(authorUserInfoURI)
          .request(MediaType.APPLICATION_JSON_TYPE)
          .header("Authorization", authorization)
          .get(JsonObject.class);

      if (_log.isInfoEnabled()) {
        _log.info(response.toString());
      }
    }

    return Response.accepted().build();
  }

  @ConfigProperty(name = "com.liferay.lxc.dxp.mainDomain")
  String liferayDomain;

  @ConfigProperty(name = "com.liferay.lxc.dxp.server.protocol")
  String liferayProtocol;

  private static final Logger _log = LoggerFactory.getLogger(ObjectAction1.class);

}
