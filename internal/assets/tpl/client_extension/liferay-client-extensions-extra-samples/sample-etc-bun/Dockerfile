FROM oven/bun:1 AS base
WORKDIR /usr/src/app

# install dependencies into temp directory
# this will cache them and speed up future builds
FROM base AS install
RUN mkdir -p /temp/dev
COPY package.json bun.lock /temp/dev/
RUN cd /temp/dev && bun install --frozen-lockfile

# install with --production (exclude devDependencies)
RUN mkdir -p /temp/prod
COPY package.json bun.lock /temp/prod/
RUN cd /temp/prod && bun install --frozen-lockfile --production

# copy node_modules from temp directory
# then copy all (non-ignored) project files into the image
FROM base AS prerelease
COPY --from=install /temp/dev/node_modules node_modules
COPY package.json .
COPY src src

COPY dxp-metadata dxp-metadata
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.domains
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.mainDomain

# copy production dependencies and source code into final image
FROM base AS release
COPY --from=install /temp/prod/node_modules node_modules
COPY --from=prerelease /usr/src/app/src src
COPY --from=prerelease /usr/src/app/package.json .
COPY --from=prerelease /usr/src/app/dxp-metadata dxp-metadata

# add path to Liferay client extension metadata
ENV LIFERAY_ROUTES_DXP=/usr/src/app/dxp-metadata

# run the app
USER bun
EXPOSE 8228/tcp
ENTRYPOINT [ "bun", "run", "src/index.ts" ]