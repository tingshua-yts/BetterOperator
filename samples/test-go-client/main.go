package main

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func test_list() {
	log.Print("Starting the shared informed app")
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/tingshuai.yts/.kube/config.new-ai-studio")
	if err != nil {
		log.Panic(err.Error())
	}
	// 1. create go-client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}
	// 2. create informer factory
	factory := informers.NewSharedInformerFactory(clientset, 0)

	stopCh := make(chan struct{})

	// 3. create  lister and informer, ps: must before start factory
	podLister := factory.Core().V1().Pods().Lister()
	informer := factory.Core().V1().Pods().Informer()

	// 4. start factory
	factory.Start(stopCh)

	// 5. wait cache sync inform
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	// 6. list or get object
	namespace := "default"
	pods, err := podLister.Pods(namespace).List(labels.Everything())

	//pods, err := podLister.List(nil)
	if err != nil {
		log.Fatalf("list ERROR: %s", err.Error())
	} else {
		log.Print("success list pod")
		log.Print(len(pods))
		for _, pod := range pods {
			log.Print("print pod")
			log.Printf("pod name is %s", pod.ObjectMeta.Name)
		}
	}

}

func test_inform() {
	log.Print("Starting the shared informed app")
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/.../.kube/config")
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		DeleteFunc: onDelete,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// when a new pod is deployed the onAdd function would be invoked
// for now just print the event.
func onAdd(obj interface{}) {
	// Cast the obj as Pod
	pod := obj.(*corev1.Pod)
	podName := pod.GetName()
	fmt.Println("Pod started -> ", podName)
}

// when a pod is deleted the onDelete function would be invoked
// for now just print the event
func onDelete(obj interface{}) {
	pod := obj.(*corev1.Pod)
	podName := pod.GetName()
	fmt.Println("Pod deleted -> ", podName)
}
func main() {
	test_list()
}
