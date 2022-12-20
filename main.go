// This program uses the Kubernetes client-go library to create a new token using the TokenRequest API, and then creates a
// kubeconfig file using the token.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Cluster holds the cluster data
type Cluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

//Clusters hold an array of the clusters that would exist in the config file
type Clusters []struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}

//Context holds the cluster context
type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

//Contexts holds an array of the contexts
type Contexts []struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

//Users holds an array of the users that would exist in the config file
type Users []struct {
	User User   `yaml:"user"`
	Name string `yaml:"name"`
}

//User holds the user authentication data
type User struct {
	Token string `yaml:"token"`
}

//KubeConfig holds the necessary data for creating a new KubeConfig file
type KubeConfig struct {
	APIVersion     string   `yaml:"apiVersion"`
	Clusters       Clusters `yaml:"clusters"`
	Contexts       Contexts `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
	Kind           string   `yaml:"kind"`
	Preferences    struct{} `yaml:"preferences"`
	Users          Users    `yaml:"users"`
}

func initKubeClient() (*kubernetes.Clientset, clientcmd.ClientConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Printf("initKubeClient: failed creating ClientConfig with", err)
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("initKubeClient: failed creating Clientset with", err)
		return nil, nil, err
	}
	return clientset, kubeConfig, nil
}

func main() {
	serviceAccountName := flag.String("service-account", "default", "The service account to use for the token")
	namespace := flag.String("namespace", "default", "The namespace to use for the token")
	outputFile := flag.String("output-file", "", "The name of the kubeconfig file to create")
	expirationSeconds := flag.Int64("expiration-seconds", 3600, "The expiration time of the token in seconds")
	flag.Parse()

	if *outputFile == "" {
		*outputFile = *serviceAccountName + ".kubeconfig"
	}

	clientset, kubeConfig, err := initKubeClient()
	if err != nil {
		log.Fatal(err)
	}

	tokenRequest := &authv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *serviceAccountName,
			Namespace: *namespace,
		},
		Spec: authv1.TokenRequestSpec{
			ExpirationSeconds: expirationSeconds,
			Audiences:         []string{"api"},
		},
	}

	result, err := clientset.CoreV1().ServiceAccounts(*namespace).CreateToken(context.TODO(), *serviceAccountName, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status.Token)
	raw, err := kubeConfig.RawConfig()
	fmt.Println(raw.Contexts[raw.CurrentContext].Cluster)
}
