FROM golang:1.21 as builder
ARG BUILD_DIRECTORY=/build

ENV CGO_ENABLED=0 GOPATH=""

RUN mkdir $BUILD_DIRECTORY
#COPY go.mod $BUILD_DIRECTORY
#COPY go.sum $BUILD_DIRECTORY

WORKDIR $BUILD_DIRECTORY
#RUN go mod download

COPY . $BUILD_DIRECTORY
RUN go build -o auth ./cmd/main

FROM scratch
ARG BUILD_DIRECTORY=/build

COPY --from=builder $BUILD_DIRECTORY/auth .
COPY --from=builder $BUILD_DIRECTORY/.env .

EXPOSE 9000

CMD ["./auth", "run"]