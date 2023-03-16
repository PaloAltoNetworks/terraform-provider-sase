---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_decryption_rules Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_decryption_rules (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The uuid of the resource

### Read-Only

- `action` (String)
- `category` (List of String)
- `description` (String)
- `destination` (List of String)
- `destination_hip` (List of String)
- `disabled` (Boolean)
- `from` (List of String)
- `id` (String) The object ID.
- `log_fail` (Boolean)
- `log_setting` (String)
- `log_success` (Boolean)
- `name` (String)
- `negate_destination` (Boolean)
- `negate_source` (Boolean)
- `profile` (String)
- `service` (List of String)
- `source` (List of String)
- `source_hip` (List of String)
- `source_user` (List of String)
- `tag` (List of String)
- `to` (List of String)
- `type` (Attributes) (see [below for nested schema](#nestedatt--type))

<a id="nestedatt--type"></a>
### Nested Schema for `type`

Read-Only:

- `ssl_forward_proxy` (Boolean)
- `ssl_inbound_inspection` (String)

