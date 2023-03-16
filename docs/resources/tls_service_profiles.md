---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_tls_service_profiles Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_tls_service_profiles (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `certificate` (String)
- `folder` (String) The folder of the entry
- `name` (String)
- `protocol_settings` (Attributes) (see [below for nested schema](#nestedatt--protocol_settings))

### Read-Only

- `id` (String) The object ID.
- `object_id` (String)

<a id="nestedatt--protocol_settings"></a>
### Nested Schema for `protocol_settings`

Optional:

- `auth_algo_sha1` (Boolean)
- `auth_algo_sha256` (Boolean)
- `auth_algo_sha384` (Boolean)
- `enc_algo3des` (Boolean)
- `enc_algo_aes128_cbc` (Boolean)
- `enc_algo_aes128_gcm` (Boolean)
- `enc_algo_aes256_cbc` (Boolean)
- `enc_algo_aes256_gcm` (Boolean)
- `enc_algo_rc4` (Boolean)
- `keyxchg_algo_dhe` (Boolean)
- `keyxchg_algo_ecdhe` (Boolean)
- `keyxchg_algo_rsa` (Boolean)
- `max_version` (String)
- `min_version` (String)

