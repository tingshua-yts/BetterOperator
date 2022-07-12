package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func getClinet() clientset.Interface {
	kubeconfig := "/Users/tingshuai.yts/.kube/config.online"
	var (
		config *rest.Config
		err    error
	)

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}

	client := clientset.NewForConfigOrDie(config)
	return client
}
func GetUserToken(loginName string) (string, error) {
	kubeAINamespace := "kube-ai"

	kubeclient := getClinet()
	//name := "test-ram-1673747224140186.onaliyun.com"
	sa, err := kubeclient.CoreV1().ServiceAccounts(kubeAINamespace).Get(context.TODO(), loginName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		klog.Errorf("fail to get sa: %v", err)
	}
	if len(sa.Secrets) < 1 {
		return "", errors.New(fmt.Sprintf("service account %s secret count is zero", loginName))
	}

	secretName := sa.Secrets[0].Name

	secret, err := kubeclient.CoreV1().Secrets(kubeAINamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		klog.Errorf("get secret failed, ns:%s name:%s, err:%v", kubeAINamespace, secretName, err)
		return "", err
	}
	token := string(secret.Data["token"])
	return token, nil
}
func main() {

	token, err := GetUserToken("test-ram-1673747224140186.onaliyun.com")
	if err != nil {
		klog.Info("fail to get user token", err)
	}
	klog.Info("token:" + token)

}
