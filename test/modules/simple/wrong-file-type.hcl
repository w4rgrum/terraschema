# Copyright 2024 Hewlett Packard Enterprise Development LP

variable "a_variable_in_the_wrong_file" {
    type = string
    description = "A string. This should not show up in the schema, and is ignored by terraform."
}