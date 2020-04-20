#!/bin/bash
#From https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-with-cert-manager-on-digitalocean-kubernetes

# Ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
# NodePort
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/baremetal/service-nodeport.yaml
# check nodeport
kubectl get svc --namespace=ingress-nginx

# Create test namespace
kubectl create ns test

# Create service
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: echo1
  namespace: test
spec:
  ports:
  - port: 80
    targetPort: 5678
    nodePort: 31111
  selector:
    app: echo1
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo1
  namespace: test
spec:
  selector:
    matchLabels:
      app: echo1
  replicas: 1
  template:
    metadata:
      labels:
        app: echo1
    spec:
      containers:
      - name: echo1
        image: hashicorp/http-echo
        args:
        - "-text=echo1"
        ports:
        - containerPort: 5678
EOF

# Creae Ingress resource
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: echo-ingress
#   namespace: test
spec:
  rules:
  - host: echo1.jonasburster.de
    http:
      paths:
      - backend:
          serviceName: echo1
          servicePort: 80
EOF

# cat <<EOF | kubectl apply -f -
# apiVersion: apps/v1
# kind: Deployment
# metadata: 
#   name: hello-world
#   namespace: test
# spec: 
#   selector:
#     matchLabels:
#       app: hello-world
#   replicas: 1
#   template: 
#     metadata: 
#       labels: 
#         app: hello-world
#     spec: 
#       containers: 
#         - image: "gokul93/hello-world:latest"
#           imagePullPolicy: Always
#           name: hello-world-container
#           ports: 
#             - containerPort: 8080
# ---
# apiVersion: v1
# kind: Service
# metadata: 
#   name: hello-world
#   namespace: test
# spec: 
#   ports: 
#      -  port: 8080
#         protocol: TCP
#         targetPort: 8080
#         nodePort: 31112
#   selector: 
#     app: hello-world
#   type: NodePort
# EOF