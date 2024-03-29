---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_http_header_profiles Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_http_header_profiles (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `object_id` (String) The uuid of the resource.

### Read-Only

- `description` (String) The `description` parameter.
- `http_header_insertion` (Attributes List) The `http_header_insertion` parameter. (see [below for nested schema](#nestedatt--http_header_insertion))
- `id` (String) The object ID.
- `name` (String) The `name` parameter.

<a id="nestedatt--http_header_insertion"></a>
### Nested Schema for `http_header_insertion`

Read-Only:

- `name` (String) The `name` parameter.
- `type` (Attributes List) The `type` parameter. (see [below for nested schema](#nestedatt--http_header_insertion--type))

<a id="nestedatt--http_header_insertion--type"></a>
### Nested Schema for `http_header_insertion.type`

Read-Only:

- `domains` (List of String) The `domains` parameter.
- `headers` (Attributes List) The `headers` parameter. (see [below for nested schema](#nestedatt--http_header_insertion--type--headers))
- `name` (String) The `name` parameter.

<a id="nestedatt--http_header_insertion--type--headers"></a>
### Nested Schema for `http_header_insertion.type.headers`

Read-Only:

- `header` (String) The `header` parameter.
- `log` (Boolean) The `log` parameter.
- `name` (String) The `name` parameter.
- `value` (String) The `value` parameter.


