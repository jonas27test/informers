package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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

	namespace := "inf"
	resource := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1alpha3",
		Resource: "certificates",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "inf",
			Namespace: "inf",
		},
		Spec: appsv1.DeploymentSpec{
			// Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "scratch",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	r, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	log.Println(r.GetResourceVersion())

	list, err := client.Resource(resource).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
		log.Println(list)
	}
	list.GetResourceVersion()
	controller := false
	blockOwnerDeletion := false
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

	s, err := json.Marshal(c)
	var dat map[string]interface{}
	if err := json.Unmarshal(s, &dat); err != nil {
		panic(err)
	}
	result, err := dclientset.Resource(resource).Namespace("inf").Create(context.TODO(), &unstructured.Unstructured{dat}, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}
