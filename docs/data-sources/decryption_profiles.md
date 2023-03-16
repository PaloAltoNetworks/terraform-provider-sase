---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_decryption_profiles Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_decryption_profiles (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The uuid of the resource

### Read-Only

- `id` (String) The object ID.
- `name` (String)
- `ssl_forward_proxy` (Attributes) (see [below for nested schema](#nestedatt--ssl_forward_proxy))
- `ssl_inbound_proxy` (Attributes) (see [below for nested schema](#nestedatt--ssl_inbound_proxy))
- `ssl_no_proxy` (Attributes) (see [below for nested schema](#nestedatt--ssl_no_proxy))
- `ssl_protocol_settings` (Attributes) (see [below for nested schema](#nestedatt--ssl_protocol_settings))

<a id="nestedatt--ssl_forward_proxy"></a>
### Nested Schema for `ssl_forward_proxy`

Read-Only:

- `auto_include_altname` (Boolean)
- `block_client_cert` (Boolean)
- `block_expired_certificate` (Boolean)
- `block_timeout_cert` (Boolean)
- `block_tls13_downgrade_no_resource` (Boolean)
- `block_unknown_cert` (Boolean)
- `block_unsupported_cipher` (Boolean)
- `block_unsupported_version` (Boolean)
- `block_untrusted_issuer` (Boolean)
- `restrict_cert_exts` (Boolean)
- `strip_alpn` (Boolean)


<a id="nestedatt--ssl_inbound_proxy"></a>
### Nested Schema for `ssl_inbound_proxy`

Read-Only:

- `block_if_hsm_unavailable` (Boolean)
- `block_if_no_resource` (Boolean)
- `block_unsupported_cipher` (Boolean)
- `block_unsupported_version` (Boolean)


<a id="nestedatt--ssl_no_proxy"></a>
### Nested Schema for `ssl_no_proxy`

Read-Only:

- `block_expired_certificate` (Boolean)
- `block_untrusted_issuer` (Boolean)


<a id="nestedatt--ssl_protocol_settings"></a>
### Nested Schema for `ssl_protocol_settings`

Read-Only:

- `auth_algo_md5` (Boolean)
- `auth_algo_sha1` (Boolean)
- `auth_algo_sha256` (Boolean)
- `auth_algo_sha384` (Boolean)
- `enc_algo3des` (Boolean)
- `enc_algo_aes128_cbc` (Boolean)
- `enc_algo_aes128_gcm` (Boolean)
- `enc_algo_aes256_cbc` (Boolean)
- `enc_algo_aes256_gcm` (Boolean)
- `enc_algo_chacha20_poly1305` (Boolean)
- `enc_algo_rc4` (Boolean)
- `keyxchg_algo_dhe` (Boolean)
- `keyxchg_algo_ecdhe` (Boolean)
- `keyxchg_algo_rsa` (Boolean)
- `max_version` (String)
- `min_version` (String)

