---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_certificate_profiles Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_certificate_profiles (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ca_certificates` (Attributes List) The `ca_certificates` parameter. (see [below for nested schema](#nestedatt--ca_certificates))
- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `name` (String) The `name` parameter. String length must be at most 63.

### Optional

- `block_expired_cert` (Boolean) The `block_expired_cert` parameter.
- `block_timeout_cert` (Boolean) The `block_timeout_cert` parameter.
- `block_unauthenticated_cert` (Boolean) The `block_unauthenticated_cert` parameter.
- `block_unknown_cert` (Boolean) The `block_unknown_cert` parameter.
- `cert_status_timeout` (String) The `cert_status_timeout` parameter.
- `crl_receive_timeout` (String) The `crl_receive_timeout` parameter.
- `domain` (String) The `domain` parameter.
- `ocsp_receive_timeout` (String) The `ocsp_receive_timeout` parameter.
- `use_crl` (Boolean) The `use_crl` parameter.
- `use_ocsp` (Boolean) The `use_ocsp` parameter.
- `username_field` (Attributes) The `username_field` parameter. (see [below for nested schema](#nestedatt--username_field))

### Read-Only

- `id` (String) The object ID.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--ca_certificates"></a>
### Nested Schema for `ca_certificates`

Optional:

- `default_ocsp_url` (String) The `default_ocsp_url` parameter.
- `name` (String) The `name` parameter.
- `ocsp_verify_cert` (String) The `ocsp_verify_cert` parameter.
- `template_name` (String) The `template_name` parameter.


<a id="nestedatt--username_field"></a>
### Nested Schema for `username_field`

Optional:

- `subject` (String) The `subject` parameter. Value must be one of: `"common-name"`.
- `subject_alt` (String) The `subject_alt` parameter. Value must be one of: `"email"`.


