FROM nginx:stable-alpine

ARG HOST_URL

RUN apk update && apk add openssl

RUN openssl dhparam -out /etc/ssl/certs/dhparam.pem 2048 &&\
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/certs/dev.key -out /etc/ssl/certs/dev.crt \
    -subj "/C=RU/ST=Moscow/L=Kremlin/O=DIT Code/OU=Org/CN=${HOST_URL}"

RUN rm /etc/nginx/conf.d/default.conf
COPY dev.conf /etc/nginx/conf.d/dev.conf
COPY proxy.conf.template /etc/nginx/templates/10-variables.conf.template
COPY ssl-params.conf /etc/nginx/ssl-params.conf
