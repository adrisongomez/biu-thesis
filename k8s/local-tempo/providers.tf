# providers.tf

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.11.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.5.1"
    }
  }
}

# This provider connects to Kubernetes itself
provider "kubernetes" {
    config_path    = "~/.kube/config"
  config_context = "minikube"
}

provider "helm" {
  kubernetes = {
    config_path    = "~/.kube/config"
    config_context = "minikube"
  }
}