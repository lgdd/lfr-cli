FROM gradle:8.5-jdk21-alpine AS builder

WORKDIR /sample-etc-java-vertx

COPY . .

RUN sed -i -e 's/localhost/liferay/g' src/main/resources/conf/dxp-metadata.json
RUN sed -i -e 's/localhost/liferay/g' src/main/resources/conf/dxp-metadata.json

RUN gradle clean assemble --no-daemon

FROM eclipse-temurin:21-alpine

RUN mkdir /opt/app

COPY --from=builder /sample-etc-java-vertx/build/libs/*-fat.jar /opt/app/vertx.jar

EXPOSE 8082

HEALTHCHECK --interval=5s --timeout=3s --retries=5 --start-period=2s \
  CMD curl --silent --fail 127.0.0.1:8082/ready || exit 1

CMD ["java", "-jar", "/opt/app/vertx.jar"]