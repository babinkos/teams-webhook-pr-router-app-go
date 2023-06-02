#! /usr/bin/env bash
mkdir -p ./certs
cd ./certs
openssl genrsa 4096 > ca-privatekey.pem
openssl req -new -key ca-privatekey.pem -sha256 -out csr.pem
openssl x509 -req -days 3000 -in csr.pem -sha256 -signkey ca-privatekey.pem -out ca-public.crt

openssl req -newkey rsa:4096 -sha256 -nodes -subj "/CN=192.168.1.8" -keyout privkey_user.pem -out "csr_user.pem"
openssl x509 -req -in csr_user.pem -days 3000 -CA ca-public.crt -CAkey ca-privatekey.pem -CAcreateserial -out public_user.crt
