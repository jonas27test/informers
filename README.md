# informers

Change yaml to enable owner delition
add to Deployment(cert-manager).spec.template.containers.cert-manager.args
`- --enable-certificate-owner-ref=true`