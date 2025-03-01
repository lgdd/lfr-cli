FROM golang:1.24-alpine3.21

WORKDIR /usr/src/app

COPY go.mod .
RUN go mod download

COPY . .

RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.domains
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.mainDomain

RUN go build -v -o /usr/local/bin/app ./...

ENTRYPOINT ["app"]