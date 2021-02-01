FROM golang:1.15-alpine AS cert-builder

RUN apk add --no-cache git bash

# install mkcert
RUN git clone https://github.com/FiloSottile/mkcert /mkcert
WORKDIR /mkcert
RUN go build -ldflags "-X main.Version=$(git describe --tags)"

# make cert
WORKDIR /certs
RUN chmod +x /mkcert/mkcert
RUN /mkcert/mkcert localhost




FROM nginx:1.19-alpine AS production

COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=cert-builder /certs/localhost.pem /etc/nginx/conf.d/cert.pem
COPY --from=cert-builder /certs/localhost-key.pem /etc/nginx/conf.d/key.pem

