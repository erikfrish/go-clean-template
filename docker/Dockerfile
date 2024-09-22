ARG GO_VER
ARG ALPINE_VER
FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder

ARG VERSION
WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -a -installsuffix cgo -o app -ldflags "-X 'main.version=${VERSION}'" ./cmd/app


FROM alpine:${ALPINE_VER}

RUN apk add --no-cache tzdata
ENV TZ=UTC
WORKDIR /root/
COPY --from=builder /src/app .
COPY ["config/config.yml", "config/"]
COPY ["deploy/migrations/*.sql", "deploy/migrations/"]
CMD ["./app"]
