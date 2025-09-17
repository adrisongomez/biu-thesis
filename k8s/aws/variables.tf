variable "aws_region" {
  description = "The AWS region to deploy resources in."
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "The name of the EKS cluster."
  type        = string
  default     = "thesis-eks-cluster"
}

variable "cluster_version" {
  description = "the kubernetes version for the EKS cluster."
  type        = string
  default     = "1.29"
}
