---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase Provider"
subcategory: ""
description: |-
  Interact with HashiCups.
---

# sase Provider

Interact with HashiCups.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) The CLient ID for the connection. Environment variable: `SASE_CLIENT_ID`. JSON config file variable: `client_id`.
- `client_secret` (String, Sensitive) The client secret for the connection. Environment variable: `SASE_CLIENT_SECRET`. JSON config file variable: `client_secret`.
- `scope` (String) The client scope. Environment variable: `SASE_SCOPE`. JSON config file variable: `scope`.

### Optional

- `auth_file` (String) The file path to the JSON file with auth creds for SASE.
- `host` (String) The hostname. Default: `api.sase.paloaltonetworks.com`. Environment variable: `SASE_HOST`. JSON config file variable: `host`.
- `logging` (String) The logging level of the provider and underlying communication. Default: `basic`. Environment variable: `SASE_LOGGING`. JSON config file variable: `logging`.