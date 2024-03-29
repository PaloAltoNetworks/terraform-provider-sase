---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_trusted_certificate_authorities_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_trusted_certificate_authorities_list (Data Source)

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

- `common_name` (String) The `common_name` parameter.
- `expiry_epoch` (String) The `expiry_epoch` parameter.
- `filename` (String) The `filename` parameter.
- `issuer` (String) The `issuer` parameter.
- `name` (String) The `name` parameter.
- `not_valid_after` (String) The `not_valid_after` parameter.
- `not_valid_before` (String) The `not_valid_before` parameter.
- `object_id` (String) The `object_id` parameter.
- `serial_number` (String) The `serial_number` parameter.
- `subject` (String) The `subject` parameter.


