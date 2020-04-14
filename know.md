docker desktop (dd):
 PVC handling - dd legt PVC 
 PVs darf ich nicht anlegen da dd sie selber anlegt

Base image baut ein debian und installiert kubectl, operator, und andere dependencies 
    operator sdk redhat, wird verwendet um openshift zu machen.
    K8s hat custom resources cr --> Red hat hat cr in openshift was leben einfacher macht

Linux dev tools
    lokale scripts von Alex

Wir er

Developer images:
    wird verwndet um operator auf jenkins zu bauen, wird von keinem verwendet.

overlays
    admin bruacht admin rechte yum installieren
    cloud braucht 

Operator ist programmierter admin
    Er reagiert auf events und reagiert darauf.
    Ignore event aber reconcile alles

    an manchen stellen PV (persistent volumens) unabhaengig 


for step loop:
    Implementierung soll geandert werden



Wie sollen wir die certificates anlegen


XPro mandant andert sich

Operator bekommt es mit und kombiniert neue eintraege in neue domain.

Pro register domain gibt es ein zertificate (welche alle dns pro mandant) with tracking domain
operator muss wissen wie die reg domains aussieht
psl (public suffix list) tool um register , hat func print-reg-domain
go --> https://godoc.org/golang.org/x/net/

Delete:
Idee eine CR custom resource pro tracking domain mit dem gleichen name wie mandant und dann finalizer.
Alle  resourcen bekommen ownership reference


Falls update von xpromandat
    check ob neue dns die noch in keinem certificate ist.
    Falls ja, trage neue subdomain ein und generiere neues ssl.
    Falls nein, tu nichts.
Falls mandant geloescht wird   
    check certificate, trage dns aus und fordere neues certifikat an.
Maybe batch/background job: j

Operator legt einen Certificate an was sagt ich bruach ein certificate --> certmanager - produziert -> neues secret



Wie ist es mit managed falsed/true
    sollten mandaten auserhalb von k8s auch certificate bekommen.

docker run -it correkter image tag ist xpro-operator-677f84b499-swb5n

Fragen:
 - Wie sollen wir long-term die DNS namen sichern. publicsuffix gibt nur ob eine domain gemanaged ist
 - Was wir machen wenn ein mandat gelöscht wird
 - Was sollen wir damit umgehen, dass manadaten gemanaged sein können.




Xpro-Operator zum laufen bekommen


docker build -t local/xpro:v0.jonas .
docker build -t inx.dockreg.net/xpro-operator-dev:VERSION developer-image/docker
docker build -f .  -t inx.dockreg.net/xpro-operator-dev:VERSION developer-image/docker
docker build -t inx.dockreg.net/xpro-operator-dev:VERSION -f developer-image/docker .
docker build -t inx.dockreg.net/xpro-operator-dev:VERSION -f developer-image/docker/Dockerfile .
kubectl create namespace xpro-operator
kubectl -n xpro-operator create secret docker-registry docker-artifactory --docker-server=https://inx.dockreg.net/ --docker-username=jobu --docker-password=AP3QFSCYmJBNeALRb57KXsa5DZN --docker-email=Jonas.Burster@inxmail.de

cd ..\xpro-k8s-operator\xpro-operator\
docker run --rm -it -v C:\Users\jobu\repos\xpro-k8s-operator\xpro-operator:/code  inx.dockreg.net/xpro-operator-dev:latest /bin/bash