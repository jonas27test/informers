#!/bin/bash

COMMON_NAME = "inf.com"

# https://docs.cert-manager.io/en/release-0.9/tasks/issuers/setup-ca.html
# Generate the singing key
# openssl genrsa -out ca.key 2048
# openssl req -x509 -new -nodes -key ca.key -subj "/CN=${COMMON_NAME}" -days 3650 -reqexts v3_req -extensions v3_ca -out ca.crt
# Save the signing key pair as a Secret in k8s
microk8s.kubectl create secret tls ca-key-pair --cert=ca.crt --key=ca.key --namespace=inf

microk8s.kubectl apply -f cert-manager/issuer.yaml