---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_qos_policy_rules Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_qos_policy_rules (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (Attributes) The `action` parameter. (see [below for nested schema](#nestedatt--action))
- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `name` (String) The `name` parameter.
- `position` (String) The position of a security rule. Value must be one of: `"pre"`, `"post"`.

### Optional

- `description` (String) The `description` parameter.
- `dscp_tos` (Attributes) The `dscp_tos` parameter. (see [below for nested schema](#nestedatt--dscp_tos))
- `schedule` (String) The `schedule` parameter.

### Read-Only

- `id` (String) The object ID.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--action"></a>
### Nested Schema for `action`

Optional:

- `class` (String) The `class` parameter.


<a id="nestedatt--dscp_tos"></a>
### Nested Schema for `dscp_tos`

Optional:

- `codepoints` (Attributes List) The `codepoints` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints))

<a id="nestedatt--dscp_tos--codepoints"></a>
### Nested Schema for `dscp_tos.codepoints`

Optional:

- `name` (String) The `name` parameter.
- `type` (Attributes) The `type` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type))

<a id="nestedatt--dscp_tos--codepoints--type"></a>
### Nested Schema for `dscp_tos.codepoints.type`

Optional:

- `af` (Attributes) The `af` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type--af))
- `cs` (Attributes) The `cs` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type--cs))
- `custom` (Attributes) The `custom` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type--custom))
- `ef` (Boolean) The `ef` parameter.
- `tos` (Attributes) The `tos` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type--tos))

<a id="nestedatt--dscp_tos--codepoints--type--af"></a>
### Nested Schema for `dscp_tos.codepoints.type.tos`

Optional:

- `codepoint` (String) The `codepoint` parameter.


<a id="nestedatt--dscp_tos--codepoints--type--cs"></a>
### Nested Schema for `dscp_tos.codepoints.type.tos`

Optional:

- `codepoint` (String) The `codepoint` parameter.


<a id="nestedatt--dscp_tos--codepoints--type--custom"></a>
### Nested Schema for `dscp_tos.codepoints.type.tos`

Optional:

- `codepoint` (Attributes) The `codepoint` parameter. (see [below for nested schema](#nestedatt--dscp_tos--codepoints--type--tos--codepoint))

<a id="nestedatt--dscp_tos--codepoints--type--tos--codepoint"></a>
### Nested Schema for `dscp_tos.codepoints.type.tos.codepoint`

Optional:

- `binary_value` (String) The `binary_value` parameter.
- `codepoint_name` (String) The `codepoint_name` parameter.



<a id="nestedatt--dscp_tos--codepoints--type--tos"></a>
### Nested Schema for `dscp_tos.codepoints.type.tos`

Optional:

- `codepoint` (String) The `codepoint` parameter.


