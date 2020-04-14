package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	scheme "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/fatih/structs"
	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
)

// var (
// 	localSchemeBuilder = &SchemeBuilder
// )

func main() {

	log.SetFlags(log.Lshortfile)
	log.Println("Shared Informer app started")
	config, err := clientcmd.BuildConfigFromFlags("", "./conf/config")
	if err != nil {
		log.Panicln(err.Error())
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Panicln(err.Error())
	}
	// log.
	// factory := informers.NewSharedInformerFactory(clientset, 0)

	defer runtime.HandleCrash()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientset := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(clientset, 0)
	log.Println(factory)

	dclientset := dynamic.NewForConfigOrDie(config)
	_ = dynamicinformer.NewDynamicSharedInformerFactory(dclientset, 0)

	// var wg sync.WaitGroup

	// routeUpdates := watchMandants(dfactory, ctx.Done(), &wg)

	namespace := "inf"
	resource := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1alpha3",
		Resource: "certificates",
	}

	// controllerLoop:
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			break controllerLoop
	// 			// case: createCert := <-
	// 		}
	// 	}
	// Create(namespace, dclientset, dfactory)
	// result, err := client.Resource(resource).Namespace(namespace).Create(context.TODO(), Create(namespace, dclientset, dfactory), metav1.CreateOptions{})
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("Created CMCertificate %q.\n", result.GetName())
	list, err := client.Resource(resource).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}
	log.Println(list)
	s := CMCertificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-inf-test",
			Namespace: "inf",
			// OwnerReferences: map[string]interface{}{},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1alpha2",
			Kind:       "Certificate",
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
	log.Println(s.TypeMeta)
	c := v1alpha3.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-inf-test",
			Namespace: "inf",
			// OwnerReferences: map[string]interface{}{},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1alpha3",
			Kind:       "Certificate",
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
			DNSNames:     []string{"inxmail.com", "internal.inxmail.com"},
			IPAddresses:  []string{"192.168.0.1"},
			IssuerRef:    cmmeta.ObjectReference{},
		},
	}
	log.Println(c.TypeMeta.Kind)
	localSchemeBuilder.Register(addTypes)
	localSchemeBuilder

	result, err := dclientset.Resource(resource).Namespace("inf").Create(context.TODO(), &unstructured.Unstructured{Object: structs.Map(c)}, metav1.CreateOptions{})
	log.Println(result)
	log.Println(err)
	log.Println("asd")

	// informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc: onAdd,
	// })
	// go informer.Run(stopper)
	// if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
	// 	runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	// 	return
	// }
	// <-stopper
}

func addTypes(scheme *scheme.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1alpha3.Certificate{},
	)
	return nil
}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onAdd(obj interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	s, ok := pod.GetLabels()[pod.Spec.String()]
	if strings.Contains(pod.Name, "informer") {
		log.Println(pod.Name)
		log.Println(s)
	}
	if ok {
		fmt.Printf("It has the label!")
	}
}
