#!/bin/bash

mkdir -p deploy/tls/certs

openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout deploy/tls/certs/key.pem \
  -out deploy/tls/certs/cert.pem \
  -days 365 \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:nginx,DNS:tasks"

echo "Сертификаты сгенерированы в deploy/tls/certs/"