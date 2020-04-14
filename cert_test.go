package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/fatih/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func before(name string) {
	log.SetFlags(log.Lshortfile)
	log.Println(name)
}

func TestCMCertificate(t *testing.T) {
	before("CertTest started")
	cert := certGen()
	log.Println(structs.Map(cert))
	certString, err := json.Marshal(cert)
	if err != nil {
		log.Println(err)
	}
	result := `{"TypeMeta":{"kind":"Certificate","apiVersion":"cert-manager.io/v1alpha2"},"metadata":{"name":"cert-inf-test","namespace":"inf","creationTimestamp":null},"spec":{"secretName":"cert-inf-test","duration":"2160h","renewBefore":"360h","organization":["inxmail.com"],"keySize":2048,"keyAlgorithm":"rsa","KeyEncoding":"pkcs1","usages":["server auth","client auth"],"dnsNames":["inxmail.com","internal.inxmail.com"],"ipAddresses":["192.168.0.1"],"issuerRef":{"name":"cl-ca-issuer","kind":"ClusterIssuer"}}}`
	if string(certString) != result {
		panic("Strings don't match!")
	}
}

func TestStruct(t *testing.T) {
	log.Println("")
}

func certGen() CMCertificate {
	return CMCertificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1alpha2",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-inf-test",
			Namespace: "inf",
		},
		Spec: CMCertSpec{
			SecretName:   "cert-inf-test",
			Duration:     "2160h",
			RenewBefore:  "360h",
			Organization: []string{"inxmail.com"},
			IsCA:         false,
			KeySize:      2048,
			KeyAlgorithm: "rsa",
			KeyEncoding:  "pkcs1",
			Usages:       []string{"server auth", "client auth"},
			DNSNames:     []string{"inxmail.com", "internal.inxmail.com"},
			IPAddresses:  []string{"192.168.0.1"},
			IssuerRef: CMIssuerRef{
				Name: "cl-ca-issuer",
				Kind: "ClusterIssuer",
			},
		},
	}
}
