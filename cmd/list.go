package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kubernetes namespaces",
	Long:  `List all namespaces in the cluster with detailed information including labels and resource counts.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	fmt.Printf("%s%s🔍 Kubernetes Namespace Lister%s\n\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("Found %d namespaces:\n\n", len(namespaces.Items))

	for i, ns := range namespaces.Items {
		printNamespace(clientset, ns)
		if i < len(namespaces.Items)-1 {
			fmt.Println()
		}
	}

	return nil
}

func printNamespace(clientset *kubernetes.Clientset, ns corev1.Namespace) {
	statusColor := ColorGreen
	if ns.Status.Phase != corev1.NamespaceActive {
		statusColor = ColorRed
	}

	age := time.Since(ns.CreationTimestamp.Time).Round(time.Second)

	fmt.Printf("%s📦 %s%s%s\n", ColorBold, statusColor, ns.Name, ColorReset)
	fmt.Printf("   %sStatus:%s %s%s%s\n", ColorGray, ColorReset, statusColor, ns.Status.Phase, ColorReset)
	fmt.Printf("   %sAge:%s %s\n", ColorGray, ColorReset, formatAge(age))

	if len(ns.Labels) > 0 {
		fmt.Printf("   %sLabels:%s\n", ColorGray, ColorReset)
		for key, value := range ns.Labels {
			fmt.Printf("      %s%s%s: %s\n", ColorCyan, key, ColorReset, value)
		}
	}

	// Count pods
	pods, _ := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
	podCount := 0
	if pods != nil {
		podCount = len(pods.Items)
	}

	// Count services
	services, _ := clientset.CoreV1().Services(ns.Name).List(context.TODO(), metav1.ListOptions{})
	serviceCount := 0
	if services != nil {
		serviceCount = len(services.Items)
	}

	fmt.Printf("   %sResources:%s\n", ColorGray, ColorReset)
	fmt.Printf("      %sPods:%s %s%d%s\n", ColorGray, ColorReset, ColorYellow, podCount, ColorReset)
	fmt.Printf("      %sServices:%s %s%d%s\n", ColorGray, ColorReset, ColorYellow, serviceCount, ColorReset)
}

func formatAge(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%dd%dh", days, hours)
}
