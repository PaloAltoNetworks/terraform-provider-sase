---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_app_override_rules Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_app_override_rules (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `object_id` (String) The uuid of the resource.

### Read-Only

- `application` (String) The `application` parameter.
- `description` (String) The `description` parameter.
- `destination` (List of String) The `destination` parameter.
- `disabled` (Boolean) The `disabled` parameter.
- `from` (List of String) The `from` parameter.
- `group_tag` (String) The `group_tag` parameter.
- `id` (String) The object ID.
- `name` (String) The `name` parameter.
- `negate_destination` (Boolean) The `negate_destination` parameter.
- `negate_source` (Boolean) The `negate_source` parameter.
- `port` (Number) The `port` parameter.
- `protocol` (String) The `protocol` parameter.
- `source` (List of String) The `source` parameter.
- `tag` (List of String) The `tag` parameter.
- `to` (List of String) The `to` parameter.


