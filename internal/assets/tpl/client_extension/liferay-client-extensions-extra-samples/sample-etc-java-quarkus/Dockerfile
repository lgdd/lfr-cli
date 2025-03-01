FROM openjdk:21-bullseye AS builder

WORKDIR /sample-etc-java-quarkus

COPY pom.xml .
COPY mvnw .
COPY .mvn .mvn

RUN ./mvnw dependency:go-offline

COPY src src

RUN sed -i -e 's/localhost/liferay/g' src/main/resources/application.properties
RUN sed -i -e 's/localhost/liferay/g' src/main/resources/application.properties

RUN ./mvnw package

FROM registry.access.redhat.com/ubi9/openjdk-21:1.21

ENV LANGUAGE='en_US:en'

COPY --from=builder --chown=185 /sample-etc-java-quarkus/target/quarkus-app/lib/ /deployments/lib/
COPY --from=builder --chown=185 /sample-etc-java-quarkus/target/quarkus-app/*.jar /deployments/
COPY --from=builder --chown=185 /sample-etc-java-quarkus/target/quarkus-app/app/ /deployments/app/
COPY --from=builder --chown=185 /sample-etc-java-quarkus/target/quarkus-app/quarkus/ /deployments/quarkus/

EXPOSE 8083
USER 185
ENV JAVA_OPTS_APPEND="-Dquarkus.http.host=0.0.0.0 -Djava.util.logging.manager=org.jboss.logmanager.LogManager"
ENV JAVA_APP_JAR="/deployments/quarkus-run.jar"

HEALTHCHECK --interval=5s --timeout=3s --retries=5 --start-period=2s \
  CMD curl --silent --fail 127.0.0.1:8083/q/health/ready || exit 1

ENTRYPOINT [ "/opt/jboss/container/java/run/run-java.sh" ]