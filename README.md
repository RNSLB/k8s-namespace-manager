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
