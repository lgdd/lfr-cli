FROM {{.Image}}

ENV JAVA_VERSION=zulu{{.JavaVersion}}

COPY --chown=liferay:liferay bundles/deploy /opt/liferay/deploy
COPY --chown=liferay:liferay bundles/osgi/war /opt/liferay/osgi/war
COPY --chown=liferay:liferay bundles/osgi/modules /opt/liferay/osgi/modules
COPY --chown=liferay:liferay bundles/osgi/configs /opt/liferay/osgi/configs
COPY --chown=liferay:liferay bundles/*.properties /opt/liferay/
