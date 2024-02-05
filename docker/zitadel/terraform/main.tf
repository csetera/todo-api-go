#
# https://registry.terraform.io/providers/zitadel/zitadel/latest
#

terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.7"
    }
  }
}

provider "zitadel" {
  domain           = "localhost"
  insecure         = "true"
  port             = "8088"
  jwt_profile_file = "admin-key.json"
}

#
# Organization and project
#
resource "zitadel_org" "default" {
  name = "todo-api-org"
}

resource "zitadel_project" "default" {
  name                     = "todo-api-go-project"
  org_id                   = zitadel_org.default.id
  project_role_assertion   = true
  project_role_check       = true
  has_project_check        = true
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}

#
# Some roles for testing
#
variable "roles" {
  description = "Roles for testing"
  type        = list(string)
  default     = ["create", "retrieve", "update", "delete"]
}

resource "zitadel_project_role" "name" {
  count = length(var.roles)

  org_id       = zitadel_org.default.id
  project_id   = zitadel_project.default.id
  role_key = var.roles[count.index]
  display_name = var.roles[count.index]
}

#
# Machine users with private access tokens
#
resource "zitadel_machine_user" "readonly" {
  org_id      = zitadel_org.default.id
  user_name   = "readonly"
  name        = "Read Only User"
  description = "Read-only user"
}

resource "zitadel_personal_access_token" "readonly" {
  org_id          = zitadel_org.default.id
  user_id         = zitadel_machine_user.readonly.id
  expiration_date = "2519-04-01T08:45:00Z"
}

resource "zitadel_machine_user" "readwrite" {
  org_id      = zitadel_org.default.id
  user_name   = "readwrite"
  name        = "Read Write User"
  description = "Read/Write user"
}

resource "zitadel_personal_access_token" "readwrite" {
  org_id          = zitadel_org.default.id
  user_id         = zitadel_machine_user.readwrite.id
  expiration_date = "2519-04-01T08:45:00Z"
}

#
# Human Users
#
resource "zitadel_human_user" "readonly" {
  org_id             = zitadel_org.default.id
  user_name          = "readonly@localhost.com"
  first_name         = "Read"
  last_name          = "Only"
  email              = "readonly@localhost.com"
  is_email_verified  = true
  initial_password   = "P455w0rd!"
}

resource "zitadel_user_grant" "readonly" {
  project_id = zitadel_project.default.id
  org_id     = zitadel_org.default.id
  user_id    = zitadel_human_user.readonly.id
  role_keys  = ["retrieve"]
}

resource "zitadel_human_user" "readwrite" {
  org_id             = zitadel_org.default.id
  user_name          = "readwrite@localhost.com"
  first_name         = "Read"
  last_name          = "Write"
  email              = "readwrite@localhost.com"
  is_email_verified  = true
  initial_password   = "P455w0rd!"
}

resource "zitadel_user_grant" "readwrite" {
  project_id = zitadel_project.default.id
  org_id     = zitadel_org.default.id
  user_id    = zitadel_human_user.readwrite.id
  role_keys  = ["create", "retrieve", "update", "delete"]
}

#
# API Application Setup
#
resource "zitadel_application_api" "default" {
  org_id           = zitadel_org.default.id
  project_id       = zitadel_project.default.id
  name             = "todo-api-go"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}

resource "zitadel_application_key" "default" {
  org_id          = zitadel_org.default.id
  project_id      = zitadel_project.default.id
  app_id          = zitadel_application_api.default.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}

#
# OIDC Application Setup
#
resource "zitadel_application_oidc" "default" {
  org_id     = zitadel_org.default.id
  project_id = zitadel_project.default.id

  name                        = "todo-api-go-oidc"
  app_type                    = "OIDC_APP_TYPE_USER_AGENT"
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  dev_mode                    = true
  redirect_uris               = ["http://localhost", "http://localhost:28080" ]
  post_logout_redirect_uris   = ["http://localhost", "http://localhost:28080" ]
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_NONE"
  access_token_role_assertion = true
  id_token_role_assertion     = true
  id_token_userinfo_assertion = true
}

#
# Output variables
#
output "org_id" {
  value = zitadel_org.default.id
  description = "The identifier of the organization"
}

output "project_id" {
  value = zitadel_project.default.id
  description = "The identifier of the project"
}

output "application_api_id" {
  value = zitadel_application_api.default.id
  description = "The identifier of the application for the API"
}

output "application_api_key" {
  value = nonsensitive(zitadel_application_key.default.key_details)
  description = "The JSON key for the application of the API"
  sensitive = true
}

output "oidc_app_client_id" {
  value = zitadel_application_oidc.default.client_id
  description = "The client identifier of the OIDC application"
  sensitive = true
}

output "readonly_pat" {
  value         = zitadel_personal_access_token.readonly.token
  description   = "The Bearer token value for the read-only user"
  sensitive = true
}

output "readwrite_pat" {
  value         = zitadel_personal_access_token.readwrite.token
  description   = "The Bearer token value for the read-write user"
  sensitive = true
}

#
# tofu init
# tofu apply
#

#
# Retrieve the API Key as JSON using
# tofu output -raw oidc_app_client_id
# tofu output -raw application_api_key > todo-api-go-key.json
# tofu output -raw readonly_pat > readonly_pat.txt
# tofu output -raw readwrite_pat > readwrite_pat.txt
#