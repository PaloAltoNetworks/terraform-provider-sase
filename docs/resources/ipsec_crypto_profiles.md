---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_ipsec_crypto_profiles Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_ipsec_crypto_profiles (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `lifetime` (Attributes) The `lifetime` parameter. (see [below for nested schema](#nestedatt--lifetime))
- `name` (String) The `name` parameter. String length must be at most 31.

### Optional

- `ah` (Attributes) The `ah` parameter. (see [below for nested schema](#nestedatt--ah))
- `dh_group` (String) The `dh_group` parameter. Default: `"group2"`. Value must be one of: `"no-pfs"`, `"group1"`, `"group2"`, `"group5"`, `"group14"`, `"group19"`, `"group20"`.
- `esp` (Attributes) The `esp` parameter. (see [below for nested schema](#nestedatt--esp))
- `lifesize` (Attributes) The `lifesize` parameter. (see [below for nested schema](#nestedatt--lifesize))

### Read-Only

- `id` (String) The object ID.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--lifetime"></a>
### Nested Schema for `lifetime`

Optional:

- `days` (Number) The `days` parameter. Value must be between 1 and 365.
- `hours` (Number) The `hours` parameter. Value must be between 1 and 65535.
- `minutes` (Number) The `minutes` parameter. Value must be between 3 and 65535.
- `seconds` (Number) The `seconds` parameter. Value must be between 180 and 65535.


<a id="nestedatt--ah"></a>
### Nested Schema for `ah`

Required:

- `authentication` (List of String) The `authentication` parameter.


<a id="nestedatt--esp"></a>
### Nested Schema for `esp`

Required:

- `authentication` (List of String) The `authentication` parameter.
- `encryption` (List of String) The `encryption` parameter.


<a id="nestedatt--lifesize"></a>
### Nested Schema for `lifesize`

Optional:

- `gb` (Number) The `gb` parameter. Value must be between 1 and 65535.
- `kb` (Number) The `kb` parameter. Value must be between 1 and 65535.
- `mb` (Number) The `mb` parameter. Value must be between 1 and 65535.
- `tb` (Number) The `tb` parameter. Value must be between 1 and 65535.


