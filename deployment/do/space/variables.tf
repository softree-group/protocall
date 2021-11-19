variable "name" {
    description = "The name of the bucket"
    type = string
}

variable "region" {
  description = "Region to place the resources in."
  type        = string
  default     = "nyc3"
}

variable "acl" {
    description = "Canned ACL applied on bucket creation (private or public-read)"
    type = string
}
