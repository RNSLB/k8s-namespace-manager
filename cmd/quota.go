package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/RNSLB/k8s-namespace-manager/pkg/quota"
)

var (
	quotaNamespace    string
	quotaName         string
	requestsCPU       string
	requestsMemory    string
	limitsCPU         string
	limitsMemory      string
	maxPods           string
	useDefaultQuota   bool
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage resource quotas for namespaces",
	Long:  `Create and manage Kubernetes ResourceQuota objects for namespace resource limits.`,
}

var quotaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a ResourceQuota in a namespace",
	Long: `Create a ResourceQuota to limit resource consumption in a namespace.

Examples:
  # Create with default values (10 CPU, 20Gi memory)
  k8s-manager quota create --namespace my-app --default

  # Create with custom values
  k8s-manager quota create --namespace my-app \
    --requests-cpu 5 \
    --requests-memory 10Gi \
    --limits-cpu 10 \
    --limits-memory 20Gi \
    --max-pods 25

  # Create named quota
  k8s-manager quota create --namespace my-app --name compute-quota --default`,
	RunE: runQuotaCreate,
}

func init() {
	rootCmd.AddCommand(quotaCmd)
	quotaCmd.AddCommand(quotaCreateCmd)

	// Required flags
	quotaCreateCmd.Flags().StringVar(&quotaNamespace, "namespace", "", "Namespace to create quota in (required)")
	quotaCreateCmd.MarkFlagRequired("namespace")

	// Optional flags
	quotaCreateCmd.Flags().StringVar(&quotaName, "name", "default-quota", "Name of the ResourceQuota")
	quotaCreateCmd.Flags().BoolVar(&useDefaultQuota, "default", false, "Use default quota values (10 CPU, 20Gi memory)")
	quotaCreateCmd.Flags().StringVar(&requestsCPU, "requests-cpu", "", "CPU requests limit (e.g., '10')")
	quotaCreateCmd.Flags().StringVar(&requestsMemory, "requests-memory", "", "Memory requests limit (e.g., '20Gi')")
	quotaCreateCmd.Flags().StringVar(&limitsCPU, "limits-cpu", "", "CPU limits (e.g., '20')")
	quotaCreateCmd.Flags().StringVar(&limitsMemory, "limits-memory", "", "Memory limits (e.g., '40Gi')")
	quotaCreateCmd.Flags().StringVar(&maxPods, "max-pods", "", "Maximum number of pods (e.g., '50')")
}

func runQuotaCreate(cmd *cobra.Command, args []string) error {
	// Build quota config
	var config quota.QuotaConfig

	if useDefaultQuota {
		// Use defaults
		config = quota.DefaultQuotaConfig(quotaNamespace)
		config.Name = quotaName
	} else {
		// Use provided values
		config = quota.QuotaConfig{
			Name:           quotaName,
			Namespace:      quotaNamespace,
			RequestsCPU:    requestsCPU,
			RequestsMemory: requestsMemory,
			LimitsCPU:      limitsCPU,
			LimitsMemory:   limitsMemory,
			MaxPods:        maxPods,
		}
	}

	// Validate config
	if err := quota.ValidateQuotaConfig(config); err != nil {
		return fmt.Errorf("invalid quota configuration: %w", err)
	}

	// Build ResourceQuota object
	resourceQuota, err := quota.BuildResourceQuota(config)
	if err != nil {
		return fmt.Errorf("failed to build ResourceQuota: %w", err)
	}

	// Get Kubernetes client
	clientset, err := getKubernetesClient()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Create ResourceQuota in Kubernetes
	created, err := clientset.CoreV1().ResourceQuotas(quotaNamespace).Create(
		context.TODO(),
		resourceQuota,
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create ResourceQuota: %w", err)
	}

	// Display success message
	fmt.Printf("✅ Created ResourceQuota: %s in namespace: %s\n", created.Name, created.Namespace)
	fmt.Println("\nHard Limits:")
	for resourceName, quantity := range created.Spec.Hard {
		fmt.Printf("  • %s: %s\n", resourceName, quantity.String())
	}

	return nil
}
