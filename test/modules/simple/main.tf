# Copyright 2024 Hewlett Packard Enterprise Development LP

terraform {
    required_version = ">= 0.13.0"
}

output "hello" {
    value = "hello"
}

output "name" {
    value = var.name
}

output "age" {
    value = var.age
}