package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

  "github.com/RNSLB/k8s-namespace-manager/pkg/validator" 
)

var (
	namespaceName string
	labelsFlag    string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Kubernetes namespace",
	Long: `Create a Kubernetes namespace with optional labels.

Labels should be provided as comma-separated key=value pairs.

Examples:
  k8s-manager create --name demo-app
  k8s-manager create --name api --labels team=engineering,env=dev`,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&namespaceName, "name", "", "Namespace name (required)")
	createCmd.Flags().StringVar(&labelsFlag, "labels", "", "Labels (key=value,key2=value2)")
	createCmd.MarkFlagRequired("name")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// STEP 1: Validate the namespace name BEFORE creating
	valid, reason := validator.Validate(namespaceName)
	if !valid {
		return fmt.Errorf("invalid namespace name '%s': %s\n\nRequirements:\n• Must be lowercase\n• Only letters (a-z), numbers (0-9), and hyphens (-)\n• Length: 1-63 characters\n• Cannot start or end with hyphen", 
			namespaceName, reason)
	}

	// STEP 2: Check if it's a reserved name
	if validator.IsReserved(namespaceName) {
		return fmt.Errorf("'%s' is a reserved Kubernetes namespace and cannot be created", namespaceName)
	}

	// STEP 3: Get Kubernetes clientset
	clientset, err := getKubernetesClient()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// STEP 4: Parse labels from flag
	labels := parseLabels(labelsFlag)

	// STEP 5: Create namespace object
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespaceName,
			Labels: labels,
		},
	}

	// STEP 6: Create namespace in Kubernetes
	created, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	// STEP 7: Print success message
	fmt.Printf("✅ Created namespace: %s\n", created.Name)
	if len(labels) > 0 {
		fmt.Println("   Labels:")
		for k, v := range labels {
			fmt.Printf("   - %s: %s\n", k, v)
		}
	}

	return nil
}

func getKubernetesClient() (*kubernetes.Clientset, error) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Build kubeconfig path
	kubeconfigPath := filepath.Join(home, ".kube", "config")

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Create and return clientset
	return kubernetes.NewForConfig(config)
}

func parseLabels(labelStr string) map[string]string {
	labels := make(map[string]string)
	if labelStr == "" {
		return labels
	}

	// Split by comma
	pairs := strings.Split(labelStr, ",")
	for _, pair := range pairs {
		// Split by equals sign
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			labels[key] = value
		}
	}

	return labels
}
