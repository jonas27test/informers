package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

type CMCertificate struct {
	Certificate &v1alpha3.Certificate
}



func (c *CMCertificate) Create(ns string, dclientset dynamic.Interface, dfactory dynamicinformer.DynamicSharedInformerFactory) {
	resource := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1alpha3",
		Resource: "Certificates",
	}
	s, err := json.Marshal(c)
	var dat map[string]interface{}
	if err := json.Unmarshal(s, &dat); err != nil {
		panic(err)
	}
	result, err := dclientset.Resource(resource).Namespace("inf").Create(context.TODO(), &unstructured.Unstructured{dat}, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("Created CMCertificate %q.\n", result.GetName())
}

func (c *CMCertificate) transformToUnstructured() {
	s, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(s, &dat); err != nil {
		panic(err)
	}
	return &unstructured.Unstructured{dat}
}
