package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	zeroReplicas int32

	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	listen     = flag.String("listen", ":6789", "listen address")
)

func init() {
	log.SetPrefix("kube-dash: ")
}

type Deployment struct {
	Name     string `json:"name"`
	PodCount int    `json:"pod_count"`
}

func main() {
	flag.Parse()

	config, err := getConfig(*kubeconfig)
	if err != nil {
		log.Fatalln("failed to get config:", err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/deployments", func(w http.ResponseWriter, r *http.Request) {
		result, err := listDeployments(clientset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bts, err := json.Marshal(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(bts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}).Methods(http.MethodGet)
	r.HandleFunc("/api/deployments/{ns}/{name}/up", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := scale(clientset, vars["name"], vars["ns"], 1); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}).Methods(http.MethodPut)
	r.HandleFunc("/api/deployments/{ns}/{name}/down", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := scale(clientset, vars["name"], vars["ns"], 0); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}).Methods(http.MethodPut)
	log.Println("started server at", *listen)
	if err := http.ListenAndServe(*listen, r); err != nil {
		log.Fatalln("failed to start http server:", err)
	}
}

func getConfig(cfg string) (*rest.Config, error) {
	if *kubeconfig == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", cfg)
}

func scale(clientset *kubernetes.Clientset, name, ns string, replicas int32) error {
	log.Printf("scaling %s on ns %s to %d replicas", name, ns, replicas)
	deploy, err := clientset.AppsV1beta1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deploy %s on ns %s: %s", name, ns, err)
	}
	deploy.Spec.Replicas = &replicas
	_, err = clientset.AppsV1beta1().Deployments(deploy.GetNamespace()).Update(deploy)
	return err
}

func listDeployments(clientset *kubernetes.Clientset) ([]Deployment, error) {
	var result []Deployment
	deploys, err := clientset.AppsV1beta1().Deployments(apiv1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to get deployments: %s", err)
	}
	for _, deploy := range deploys.Items {
		if skipNamespace(deploy.Namespace) {
			continue
		}
		result = append(result, Deployment{
			Name:     deploy.GetName(),
			PodCount: int(*deploy.Spec.Replicas),
		})
	}
	return result, err
}

func skipNamespace(ns string) bool {
	return ns == "kube-system"
}
