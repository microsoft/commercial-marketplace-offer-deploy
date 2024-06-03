variable "resource_group_name" {
  type        = string
  description = "Name of the resource group"
}

variable "location" {
  type        = string
  description = "Location of the resource group"
}

variable "required_tags" {
  type = object({
    App         = string
    Environment = string
    GBU         = string
    ITSM        = optional(string, "MANAGEMENT")
    JobWbs      = string
    Notes       = optional(string)
    Owner       = string
  })
  description = "Required Azure tags"

  validation {
    condition     = contains(["DEV", "DR", "PROD", "QA", "TEST"], var.required_tags.Environment)
    error_message = "Environment tag must be one of: DEV, TEST, QA, PROD, DR"
  }

  validation {
    condition     = contains(["COR", "FED", "INF", "MEA"], var.required_tags.GBU)
    error_message = "GBU tag must be one of: COR, FED, INF, MEA"
  }

  validation {
    condition = contains(
      ["BACKUP", "DATABASE", "MANAGEMENT", "NETWORK", "SERVER", "STORAGE"],
      var.required_tags.ITSM
    )
    error_message = "ITSM tag must be one of: BACKUP, DATABASE, MANAGEMENT, NETWORK, SERVER, STORAGE"
  }

  validation {
    condition     = can(regex("^\\d{6}-\\d{5}$", var.required_tags.JobWbs))
    error_message = "JobWbs tag must be digits in the format xxxxxx-xxxxx"
  }

  validation {
    condition     = can(regex("(?i)[A-Z0-9+_.-]+@[A-Z0-9.-]+", var.required_tags.Owner))
    error_message = "Owner tag must be a valid email address"
  }
}
