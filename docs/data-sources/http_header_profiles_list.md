---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_http_header_profiles_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_http_header_profiles_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.

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
- `http_header_insertion` (Attributes List) The `http_header_insertion` parameter. (see [below for nested schema](#nestedatt--data--http_header_insertion))
- `name` (String) The `name` parameter.
- `object_id` (String) The `object_id` parameter.

<a id="nestedatt--data--http_header_insertion"></a>
### Nested Schema for `data.http_header_insertion`

Read-Only:

- `name` (String) The `name` parameter.
- `type` (Attributes List) The `type` parameter. (see [below for nested schema](#nestedatt--data--http_header_insertion--type))

<a id="nestedatt--data--http_header_insertion--type"></a>
### Nested Schema for `data.http_header_insertion.type`

Read-Only:

- `domains` (List of String) The `domains` parameter.
- `headers` (Attributes List) The `headers` parameter. (see [below for nested schema](#nestedatt--data--http_header_insertion--type--headers))
- `name` (String) The `name` parameter.

<a id="nestedatt--data--http_header_insertion--type--headers"></a>
### Nested Schema for `data.http_header_insertion.type.name`

Read-Only:

- `header` (String) The `header` parameter.
- `log` (Boolean) The `log` parameter.
- `name` (String) The `name` parameter.
- `value` (String) The `value` parameter.


