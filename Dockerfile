FROM golang:1.22-alpine as builder

RUN mkdir /root/.lfr

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install ./...

FROM alpine

COPY --from=builder /go/bin/lfr /bin/

ENTRYPOINT ["/bin/lfr"]