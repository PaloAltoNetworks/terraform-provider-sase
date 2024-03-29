---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_regions_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_objects_regions_list (Data Source)

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

- `address` (List of String) The `address` parameter.
- `geo_location` (Attributes) The `geo_location` parameter. (see [below for nested schema](#nestedatt--data--geo_location))
- `name` (String) The `name` parameter.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--data--geo_location"></a>
### Nested Schema for `data.geo_location`

Read-Only:

- `latitude` (Number) The `latitude` parameter.
- `longitude` (Number) The `longitude` parameter.


