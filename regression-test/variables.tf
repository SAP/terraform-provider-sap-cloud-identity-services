variable "tenant_url" {
  description = "The url of the SAP Cloud Identity Services tenant."
  type        = string
}

variable "password" {
  description = "Test password for the user."
  type        = string
  default     = "Abc12345@"
}

variable "certificate" {
  description = "The base64 certificate"
  type        = string
}

variable "certificate_file_path" {
  description = "The path of the file used for authentication"
  type = string
}

variable "certificate_file_password" {
  description = "The password associated with the p12 certificate file"
  type = string
}