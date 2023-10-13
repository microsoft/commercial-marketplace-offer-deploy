provider "azurerm" {
  features {}
}

variable "string_variable" {
  type = string
}

variable "string_variable_sensitive" {
  type = string
  sensitive = true
}


variable "bool_variable" {
  type = bool
}

variable "number_variable" {
  type = number
}

variable "list_variable" {
  type = list({
    name = string
  })
}

variable "set_variable" {
  type = set(object)
}

variable "object_variable" {
  type = object({
    name = string
  })
}

variable "map_variable" {
  type = map(object({
    name = string
  }))
}