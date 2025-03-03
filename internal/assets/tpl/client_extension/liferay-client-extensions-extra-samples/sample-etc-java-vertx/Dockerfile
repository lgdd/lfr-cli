FROM gradle:8.5-jdk21-alpine AS builder

WORKDIR /sample-etc-java-vertx

COPY . .

RUN sed -i -e 's/localhost/liferay/g' src/main/resources/application.properties
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.domains
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.mainDomain

RUN gradle clean assemble --no-daemon

FROM eclipse-temurin:21-alpine

RUN mkdir /opt/app

COPY --from=builder /sample-etc-java-vertx/build/libs/*-fat.jar /opt/app/vertx.jar
COPY --from=builder /sample-etc-java-vertx/dxp-metadata /opt/app/dxp-metadata

EXPOSE 8082

# add path to Liferay client extension metadata
ENV LIFERAY_ROUTES_DXP=/opt/app/dxp-metadata

HEALTHCHECK --interval=5s --timeout=3s --retries=5 --start-period=2s \
  CMD curl --silent --fail 127.0.0.1:8082/ready || exit 1

CMD ["java", "-jar", "/opt/app/vertx.jar"]