package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

func genCert() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1alpha3",
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":      "cert-inf-test",
				"namespace": "inf",
			},
			"spec": map[string]interface{}{
				"secretName":   "cert-inf-test",
				"duration":     "2160h",
				"renewBefore":  "360h",
				"organization": []string{"inxmail.com"},
				"isCA":         false,
				"keySize":      2048,
				"keyAlgorithm": "rsa",
				"keyEncoding":  "pkcs1",
				"usages":       []string{"server auth", "client auth"},
				"dnsNames":     []string{"inxmail.com", "internal.inxmail.com"},
				"ipAddresses":  []string{"192.168.0.1"},
				"issuerRef": map[string]interface{}{
					"name":  "ca-issuer",
					"kind":  "Issuer",
					"group": "cert-manager.io",
				},
			},
		},
	}
}

func Create(ns string, dclientset dynamic.Interface, dfactory dynamicinformer.DynamicSharedInformerFactory) {
	resource := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1alpha3",
		Resource: "Certificates",
	}
	result, err := dclientset.Resource(resource).Namespace(ns).Create(context.TODO(), genCert(), metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Created CMCertificate %q.\n", result.GetName())
}
