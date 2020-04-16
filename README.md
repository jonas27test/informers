# informers

Change yaml to enable owner delition
add to Deployment(cert-manager).spec.template.containers.cert-manager.args
`- --enable-certificate-owner-ref=true`

## Ingress Comparisons
https://medium.com/flant-com/comparing-ingress-controllers-for-kubernetes-9b397483b46b