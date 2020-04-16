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

func TestStruct(t *testing.T) {
	log.Println("")
}

func genCert(name string, ownerRef string, ownderID string, ownerKind string, controller bool) *unstructured.Unstructured {
	c := v1alpha3.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1alpha3",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-inf-test",
			Namespace: "inf",
			// map[string]interface{}{},
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					Kind:               "Deployment",
					APIVersion:         r.GetResourceVersion(),
					Controller:         &controller,
					UID:                r.UID,
					Name:               r.Name,
					BlockOwnerDeletion: &blockOwnerDeletion,
				},
			},
		},
		Spec: v1alpha3.CertificateSpec{
			SecretName:  "cert-inf-test",
			Duration:    &metav1.Duration{365 * 24 * time.Hour},
			RenewBefore: &metav1.Duration{300 * 24 * time.Hour},
			// Organization: []string{"inxmail.com"},
			IsCA:         false,
			KeySize:      2048,
			KeyAlgorithm: "rsa",
			KeyEncoding:  "pkcs1",
			Usages:       []v1alpha3.KeyUsage{v1alpha3.UsageAny},
			DNSNames:     []string{"inxmail.com", "internal.inxmail.com", "inx.com"},
			IPAddresses:  []string{"192.168.0.1"},
			IssuerRef: cmmeta.ObjectReference{
				Name:  "cl-issuer",
				Kind:  "ClusterIssuer",
				Group: "cert-manager.io",
			},
		},
	}

}
