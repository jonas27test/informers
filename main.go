package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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
<<<<<<< HEAD
=======
	// log.
	// factory := informers.NewSharedInformerFactory(clientset, 0)
>>>>>>> 4197d16b49254c7c1087ce4f3f082549586abc90

	defer runtime.HandleCrash()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientset := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(clientset, 0)
	log.Println(factory)

	dclientset := dynamic.NewForConfigOrDie(config)
	dfactory := dynamicinformer.NewDynamicSharedInformerFactory(dclientset, 0)

	// var wg sync.WaitGroup

	// routeUpdates := watchMandants(dfactory, ctx.Done(), &wg)

	namespace := "inf"
	resource := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1alpha3",
		Resource: "certificates",
	}

controllerLoop:
	for {
		select {
		case <-ctx.Done():
			break controllerLoop
			// case: createCert := <-
		}
	}
	Create(namespace, dclientset, dfactory)
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
