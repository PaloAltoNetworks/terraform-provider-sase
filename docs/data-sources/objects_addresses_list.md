---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_addresses_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_objects_addresses_list (Data Source)

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

- `description` (String) The `description` parameter.
- `fqdn` (String) The `fqdn` parameter.
- `ip_netmask` (String) The `ip_netmask` parameter.
- `ip_range` (String) The `ip_range` parameter.
- `ip_wildcard` (String) The `ip_wildcard` parameter.
- `name` (String) The `name` parameter.
- `object_id` (String) The `object_id` parameter.
- `tag` (List of String) The `tag` parameter.
- `type` (String) The `type` parameter.


