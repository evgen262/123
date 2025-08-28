FROM golang:1.23-alpine AS builder
RUN apk update && \
    apk add --no-cache ca-certificates tzdata git zip && \
    update-ca-certificates

ARG RELEASE
ARG COMMIT
ARG GOPROXY_LOGIN
ARG GOPROXY_TOKEN
ENV RELEASE=${RELEASE}
ENV COMMIT=${COMMIT}
ENV GOPROXY_LOGIN=${GOPROXY_LOGIN}
ENV GOPROXY_TOKEN=${GOPROXY_TOKEN}

WORKDIR /build-app
COPY ./ .
RUN chmod +x run.sh && ./run.sh set_private_repo && ./run.sh deps && ./run.sh build linux

FROM alpine:3.21
RUN apk add --no-cache tzdata zip ca-certificates && update-ca-certificates
COPY --from=builder /build-app/bin/web-api /app/web-api
COPY --from=builder /build-app/access-list.json /app/access-list.json
WORKDIR /app/
RUN chmod +x /app/web-api
ENV TZ=Europe/Moscow
EXPOSE 80 9999 8080
ENTRYPOINT ["/app/web-api"]
