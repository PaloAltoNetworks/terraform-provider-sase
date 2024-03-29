---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_ipsec_crypto_profiles_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_ipsec_crypto_profiles_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.

### Optional

- `limit` (Number) The max count in result entry (count per page).
- `name` (String) The name of the entry.
- `offset` (Number) The offset of the result entry.

### Read-Only

- `data` (Attributes List) The `data` parameter. (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `total` (Number) The `total` parameter.

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `ah` (Attributes) The `ah` parameter. (see [below for nested schema](#nestedatt--data--ah))
- `dh_group` (String) The `dh_group` parameter.
- `esp` (Attributes) The `esp` parameter. (see [below for nested schema](#nestedatt--data--esp))
- `lifesize` (Attributes) The `lifesize` parameter. (see [below for nested schema](#nestedatt--data--lifesize))
- `lifetime` (Attributes) The `lifetime` parameter. (see [below for nested schema](#nestedatt--data--lifetime))
- `name` (String) The `name` parameter.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--data--ah"></a>
### Nested Schema for `data.ah`

Read-Only:

- `authentication` (List of String) The `authentication` parameter.


<a id="nestedatt--data--esp"></a>
### Nested Schema for `data.esp`

Read-Only:

- `authentication` (List of String) The `authentication` parameter.
- `encryption` (List of String) The `encryption` parameter.


<a id="nestedatt--data--lifesize"></a>
### Nested Schema for `data.lifesize`

Read-Only:

- `gb` (Number) The `gb` parameter.
- `kb` (Number) The `kb` parameter.
- `mb` (Number) The `mb` parameter.
- `tb` (Number) The `tb` parameter.


<a id="nestedatt--data--lifetime"></a>
### Nested Schema for `data.lifetime`

Read-Only:

- `days` (Number) The `days` parameter.
- `hours` (Number) The `hours` parameter.
- `minutes` (Number) The `minutes` parameter.
- `seconds` (Number) The `seconds` parameter.


