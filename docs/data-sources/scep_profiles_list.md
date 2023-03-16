---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_scep_profiles_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_scep_profiles_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry

### Optional

- `limit` (Number) The max count in result entry (count per page)
- `name` (String) The name of the entry
- `offset` (Number) The offset of the result entry

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `total` (Number)

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `algorithm` (Attributes) (see [below for nested schema](#nestedatt--data--algorithm))
- `ca_identity_name` (String)
- `certificate_attributes` (Attributes) (see [below for nested schema](#nestedatt--data--certificate_attributes))
- `digest` (String)
- `fingerprint` (String)
- `name` (String)
- `object_id` (String)
- `scep_ca_cert` (String)
- `scep_challenge` (Attributes) (see [below for nested schema](#nestedatt--data--scep_challenge))
- `scep_client_cert` (String)
- `scep_url` (String)
- `subject` (String)
- `use_as_digital_signature` (Boolean)
- `use_for_key_encipherment` (Boolean)

<a id="nestedatt--data--algorithm"></a>
### Nested Schema for `data.algorithm`

Read-Only:

- `rsa` (Attributes) (see [below for nested schema](#nestedatt--data--algorithm--rsa))

<a id="nestedatt--data--algorithm--rsa"></a>
### Nested Schema for `data.algorithm.rsa`

Read-Only:

- `rsa_nbits` (String)



<a id="nestedatt--data--certificate_attributes"></a>
### Nested Schema for `data.certificate_attributes`

Read-Only:

- `dnsname` (String)
- `rfc822name` (String)
- `uniform_resource_identifier` (String)


<a id="nestedatt--data--scep_challenge"></a>
### Nested Schema for `data.scep_challenge`

Read-Only:

- `dynamic_value` (Attributes) (see [below for nested schema](#nestedatt--data--scep_challenge--dynamic_value))
- `fixed` (String)
- `none` (String)

<a id="nestedatt--data--scep_challenge--dynamic_value"></a>
### Nested Schema for `data.scep_challenge.dynamic_value`

Read-Only:

- `otp_server_url` (String)
- `password` (String)
- `username` (String)

