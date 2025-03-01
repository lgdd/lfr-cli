package org.acme.liferay.client.extension.quarkus.sample;

import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.core.Response;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Path("/object/action/1")
public class ObjectAction1 {

  @POST
  public Response objectAction1(String objectEntryString) {
    if(_log.isInfoEnabled()) {
      _log.info("Object Entry: {}", objectEntryString);
    }
    return Response.accepted().build();
  }

  private static final Logger _log = LoggerFactory.getLogger(ObjectAction1.class);

}
