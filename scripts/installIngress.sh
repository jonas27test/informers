#!/bin/bash
#From https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-with-cert-manager-on-digitalocean-kubernetes

# Ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
# NodePort
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/baremetal/service-nodeport.yaml


# check nodeport
kubectl get svc -n ingress-nginx

# Create test namespace
kubectl create ns test

# Create service
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: echo1
  namespace: inf
spec:
  ports:
  - port: 80
    targetPort: 5678
    nodePort: 31111
  selector:
    app: echo1
  type: NodePort
EOF

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
  name: ingress
  namespace: katchblog
  annotations:
   nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: katchblog.com
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          serviceName: katchblog
          servicePort: 80
  - host: www.katchblog.com
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          serviceName: katchblog
          servicePort: 80
EOF
#   - host: echo.jonasburster.de
#     http:
#       paths:
#       - path: /*
#         pathType: Prefix
#         backend:
#           serviceName: echo1
#           servicePort: 80

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolume
metadata:
  name: katchblog-volume
  namespace: katchblog
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/volumes/katchblog"
EOF
---
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: katchblog-claim
  namespace: katchblog
spec:
  storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
EOF


cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: katchblog
  namespace: katchblog
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 80
    nodePort: 31111
  selector:
    app: katchblog
---
EOF

cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: katchblog
  namespace: katchblog
spec: 
  selector:
    matchLabels:
      app: katchblog
  replicas: 1
  template:
    metadata:
      labels:
        app: katchblog
    spec:
      volumes:
      - name: katchblog-storage
        persistentVolumeClaim: 
          claimName: katchblog-claim
      - name: katchblog-tls
        secret:
          secretName: katchblog-secret
      containers:
      - name: katchblog
        image: httpd:2.4
        ports:
        - containerPort: 80
        volumeMounts:
          - mountPath: /usr/local/apache2/htdocs/
            name: katchblog-storage
          - mountPath: "/tls/certs/"
            name: katchblog-tls
            readOnly: true
      restartPolicy: Always
EOF
