package main

import (
	"context"
	"log"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func before(name string) {
	log.SetFlags(log.Lshortfile)
	log.Println(name)
}

func TestDeleteService(t *testing.T) {
	before("TestDeleteService")
	clientset, _ := connect()
	ns := "inf"
	svcName1 := "inf-svc1"
	deleteService(svcName1, ns, clientset)
}

// go test -run TestCreate
func TestCreate(t *testing.T) {
	before("TestCreate")
	clientset, dclientset := connect()

	ns := "inf"
	svcName1 := "inf-svc1"
	deleteService(svcName1, ns, clientset)
	svc := createService(svcName1, ns, clientset)
	// log.Println(svc.GetResourceVersion())
	con := false
	c := CMCertificate{
		Fields: CertFields{
			Name:        "cert-inf",
			Namespace:   ns,
			CommonName:  "jonasburster.de",
			DNSNames:    []string{"jonasburster.de", "www.jonasburster.de"},
			IPAddresses: []string{"123.123.123.123"},
			Owner: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: svc.GetResourceVersion(),
				// Kind:       svc.Kind,
				Kind:               "Service",
				Name:               svc.Name,
				UID:                svc.UID,
				Controller:         &con,
				BlockOwnerDeletion: &con,
			}}}}
	c.Create(dclientset)

	// deleteService(svcName1, ns, clientset)
}

// how should we handle get. Maybe we must register certificates before, that would be way nicer
func TestGet(t *testing.T) {
	before("TestGet")
	_, dclientset := connect()
	cert := Get(dclientset)
	log.Println(cert)

	// cm := &v1alpha3.Certificate{}
	c := []byte{}
	err := cert.UnmarshalJSON(c)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(c))

	// json.Unmarshal(, cm)

}

func TestUpdate(t *testing.T) {

}

func connect() (kubernetes.Clientset, dynamic.Interface) {
	config, err := clientcmd.BuildConfigFromFlags("", "./conf/config")
	if err != nil {
		log.Panicln(err.Error())
	}
	defer runtime.HandleCrash()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	clientset := kubernetes.NewForConfigOrDie(config)
	// factory := informers.NewSharedInformerFactory(clientset, 0)
	dclientset := dynamic.NewForConfigOrDie(config)
	_ = dynamicinformer.NewDynamicSharedInformerFactory(dclientset, 0)
	return *clientset, dclientset
}

func deleteService(name string, ns string, clientset kubernetes.Clientset) {
	err := clientset.CoreV1().Services(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err)
	}
}

func createService(name string, ns string, clientset kubernetes.Clientset) *v1.Service {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name: "port",
					Port: 8080,
				},
			},
		},
	}
	svc, err := clientset.CoreV1().Services(ns).Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	return svc
}
