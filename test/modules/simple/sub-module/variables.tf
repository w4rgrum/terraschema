# Copyright 2024 Hewlett Packard Enterprise Development LP

variable "should_be_ignored" {
  type = string
  description = "Variables in sub-modules are not read."
  default = "nothing to see here"
}