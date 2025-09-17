# outputs.tf

output "cluster_endpoint" {
  description = "The endpoint for your EKS cluster's API server."
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "The name of your EKS cluster."
  value       = module.eks.cluster_name
}

output "configure_kubectl" {
  description = "Command to configure kubectl to connect to the new cluster."
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${var.cluster_name}"
}