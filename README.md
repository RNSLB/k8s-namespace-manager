# Kubernetes Namespace Manager

A Go CLI tool for managing Kubernetes namespaces programmatically using client-go.

## Features

- ✅ List namespaces with colors, labels, and resource counts
- ✅ Create namespaces with custom labels
- ✅ Built with Cobra framework
- ✅ Production-ready error handling

## Installation
```bash
git clone https://github.com/RNSLB/k8s-namespace-manager.git
cd k8s-namespace-manager
go build -o k8s-manager
```

## Usage

### List Namespaces
```bash
./k8s-manager list
```

Shows all namespaces with:
- Status (color-coded: green=Active, red=Terminating)
- Age
- Labels
- Pod count
- Service count

### Create Namespace
```bash
# Simple namespace
./k8s-manager create --name my-app

# With labels
./k8s-manager create --name api --labels team=engineering,env=prod
```

## Requirements

- Go 1.22+
- Kubernetes cluster (local kind or remote)
- kubectl configured (~/.kube/config)

## Learning Journey

Part of my weekend-focused learning plan to become a Staff Platform Engineer.

Built using:
- Go
- Cobra CLI framework
- Kubernetes client-go
- Claude AI as pair programmer

## Author

Rohit Narwani - Delivery Manager at SLB

## Next Steps

- [ ] Add namespace validation before creation
- [ ] Add ResourceQuota creation
- [ ] Add LimitRange support
- [ ] Add NetworkPolicy templates

## Resource Quota Management

### Create ResourceQuota

Create resource quotas to limit resource consumption in namespaces.

#### With Default Values
```bash
# Creates quota with sensible defaults:
# - CPU requests: 10 cores
# - Memory requests: 20Gi
# - CPU limits: 20 cores
# - Memory limits: 40Gi
# - Max pods: 50

k8s-manager quota create --namespace my-app --default
```

#### With Custom Values
```bash
k8s-manager quota create --namespace my-app \
  --name compute-quota \
  --requests-cpu 5 \
  --requests-memory 10Gi \
  --limits-cpu 10 \
  --limits-memory 20Gi \
  --max-pods 25
```

#### Verify Quota
```bash
kubectl describe resourcequota -n my-app
```

### What is ResourceQuota?

ResourceQuota limits resource consumption in a namespace:
- **Prevents resource hogging** - One team can't use all cluster resources
- **Cost control** - Limit spending per team/environment
- **Fair distribution** - Ensures resources are shared equitably

**Platform Engineering Best Practice:** Always set quotas on namespaces to prevent runaway resource usage.

