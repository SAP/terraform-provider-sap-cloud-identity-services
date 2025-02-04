# Create a basic application in SAP Cloud Identity Services
resource "ias_application" "basic_application" {
  id          = "app_1234567890"
  name        = "My Basic Application"
  description = "A basic application in SAP Cloud Identity Services"
}
