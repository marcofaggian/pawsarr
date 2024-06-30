package pawsarr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetIngressControllersPods() {
	// Create a Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// List all pods in all namespaces
	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Filter pods running ingress controllers
	for _, pod := range pods.Items {
		if pod.Labels["app"] == "traefik" || pod.Labels["app"] == "nginx-ingress" {
			fmt.Printf("Pod Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
			for _, container := range pod.Spec.Containers {
				fmt.Printf("  Container Name: %s\n", container.Name)
				fmt.Println(GetPodLogs(pod.Namespace, pod.Name, container.Name, false, client))
			}
		}
	}
}

func GetPodLogs(namespace string, podName string, containerName string, follow bool, client *kubernetes.Clientset) string {
	podLogOpts := v1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	}
	req := client.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}
