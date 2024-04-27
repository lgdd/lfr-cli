FROM {{.Image}}

ENV JAVA_VERSION=zulu{{.JavaVersion}}

ARG TARGET_ENV=prod

COPY --chown=liferay:liferay build/docker/deploy /mnt/liferay/deploy
COPY --chown=liferay:liferay build/docker/configs/$TARGET_ENV /mnt/liferay/files
