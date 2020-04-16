package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type CMCertificate struct {
	Certificate *v1alpha3.Certificate
	Fields      CertFields
}

var (
	resource = schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1alpha3", Resource: "certificates"}
)

func (c *CMCertificate) Create(dclientset dynamic.Interface) {
	// r, err := dclientset.Resource(resource).Namespace(c.Fields.Namespace).List(context.TODO(), metav1.ListOptions{})
	c.genCert()
	unstructured := c.transformToUnstructured()
	result, err := dclientset.Resource(resource).Namespace(c.Fields.Namespace).Create(context.TODO(), unstructured, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
		log.Println(result)
	}
	// log.Printf("Created CMCertificate %q.\n", result.GetName())
}

func Get(dclientset dynamic.Interface) *unstructured.Unstructured {
	result, err := dclientset.Resource(resource).Namespace("inf").Get(context.TODO(), "cert-inf", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	return result
}

func (c *CMCertificate) Update(dclientset dynamic.Interface) {
	c.genCert()
	unstructured := c.transformToUnstructured()
	_, err := dclientset.Resource(resource).Namespace("inf").Update(context.TODO(), unstructured, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
	}
}

func (c *CMCertificate) transformToUnstructured() *unstructured.Unstructured {
	s, err := json.Marshal(c.Certificate)
	if err != nil {
		log.Println(err)
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(s, &dat); err != nil {
		panic(err)
	}
	// log.Println(dat)
	return &unstructured.Unstructured{dat}
}

func (c *CMCertificate) genCert() {
	c.Certificate = &v1alpha3.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1alpha3",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Fields.Name,
			Namespace:       c.Fields.Namespace,
			OwnerReferences: c.Fields.Owner,
		},
		Spec: v1alpha3.CertificateSpec{
			SecretName:  "cert-inf-test",
			Duration:    &metav1.Duration{365 * 24 * time.Hour},
			RenewBefore: &metav1.Duration{300 * 24 * time.Hour},
			// CommonName:   c.Fields.CommonName,
			IsCA:         false,
			KeySize:      2048,
			KeyAlgorithm: "rsa",
			KeyEncoding:  "pkcs1",
			Usages:       []v1alpha3.KeyUsage{v1alpha3.UsageAny},
			DNSNames:     c.Fields.DNSNames,
			IPAddresses:  c.Fields.IPAddresses,
			IssuerRef: cmmeta.ObjectReference{
				Name:  "cl-issuer",
				Kind:  "ClusterIssuer",
				Group: "cert-manager.io",
			},
		},
	}
}

// Only define adjustablefields here.
// Issuer Global?
// Is Duration global? Then it should not be included here!
type CertFields struct {
	Name        string
	Namespace   string
	CommonName  string
	DNSNames    []string
	IPAddresses []string
	// Duration    *metav1.Duration
	// RenewBefore *metav1.Duration
	Owner  []metav1.OwnerReference
	Issuer cmmeta.ObjectReference
}
