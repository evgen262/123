FROM golang:1.23-alpine AS builder
RUN apk update && \
    apk add --no-cache ca-certificates tzdata git zip && \
    update-ca-certificates

ARG RELEASE
ARG COMMIT
ARG GOPROXY_LOGIN
ARG GOPROXY_TOKEN
ENV RELEASE=${RELEASE} \
    COMMIT=${COMMIT} \
    GOPROXY_LOGIN=${GOPROXY_LOGIN} \
    GOPROXY_TOKEN=${GOPROXY_TOKEN}

WORKDIR /build-app
COPY ./ .
RUN chmod +x run.sh && ./run.sh set_private_repo && ./run.sh deps && ./run.sh build linux

FROM alpine:3.21
RUN apk add --no-cache tzdata zip ca-certificates && update-ca-certificates
COPY --from=builder /build-app/bin/auth /app/auth
WORKDIR /app/
RUN chmod +x /app/auth
ENV TZ=Europe/Moscow
EXPOSE 9999 8080
ENTRYPOINT ["/app/auth"]