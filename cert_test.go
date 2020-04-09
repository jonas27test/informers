package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/fatih/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCMCertificate(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	log.Println("CertTest started")
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
			DnsNames:     []string{"inxmail.com", "internal.inxmail.com"},
			IpAddresses:  []string{"192.168.0.1"},
			IssuerRef: CMIssuerRef{
				Name: "cl-ca-issuer",
				Kind: "ClusterIssuer",
			},
		},
	}
}

// func dumpMap(space string, m map[string]interface{}) {
// 	for k, v := range m {
// 		if mv, ok := v.(map[string]interface{}); ok {
// 			fmt.Printf("{ \"%v\": \n", k)
// 			dumpMap(space+"\t", mv)
// 			fmt.Printf("}\n")
// 		} else {
// 			fmt.Printf("%v %v : %v\n", space, k, v)
// 		}
// 	}
// }

type CMCertificate struct {
	TypeMeta   metav1.TypeMeta   `json:",inline"`
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CMCertSpec `json:"spec"`
}

type CMCertSpec struct {
	// Secret names are always required.
	SecretName   string   `json:"secretName,required"`
	Duration     string   `json:"duration,omitempty"`
	RenewBefore  string   `json:"renewBefore,omitempty"`
	Organization []string `json:"organization,omitempty"`
	// The use of the common name field has been deprecated since 2000 and is
	// discouraged from being used.
	CommonName   string   `json:"commonName,omitempty"`
	IsCA         bool     `json:"isCA,omitempty"`
	KeySize      int64    `json:"keySize,omitempty"`
	KeyAlgorithm string   `json:"keyAlgorithm,omitempty"`
	KeyEncoding  string   `json:"KeyEncoding,omitempty"`
	Usages       []string `json:"usages,omitempty"`
	// At least one of a DNS Name, USI SAN, or IP address is required.
	DnsNames    []string `json:"dnsNames,required"`
	UriSANs     []string `json:"uriSANs,omitempty"`
	IpAddresses []string `json:"ipAddresses,omitempty"`
	// Issuer references are always required. (Either Cluster or NS)
	IssuerRef CMIssuerRef `json:"issuerRef,required"`
}

type CMIssuerRef struct {
	Name string `json:"name,required"`
	// We can reference ClusterIssuers by changing the kind here.
	// The default value is Issuer (i.e. a locally namespaced Issuer)
	Kind string `json:"kind,omitempty"`
	// This is optional since cert-manager will default to this value however
	// if you are using an external issuer, change this to that issuer group.
	Group string `json:"group,omitempty"`
}
