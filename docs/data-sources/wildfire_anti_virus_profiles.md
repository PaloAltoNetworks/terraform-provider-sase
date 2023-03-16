---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_wildfire_anti_virus_profiles Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_wildfire_anti_virus_profiles (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The uuid of the resource

### Read-Only

- `description` (String)
- `id` (String) The object ID.
- `mlav_exception` (Attributes List) (see [below for nested schema](#nestedatt--mlav_exception))
- `name` (String)
- `packet_capture` (Boolean)
- `rules` (Attributes List) (see [below for nested schema](#nestedatt--rules))
- `threat_exception` (Attributes List) (see [below for nested schema](#nestedatt--threat_exception))

<a id="nestedatt--mlav_exception"></a>
### Nested Schema for `mlav_exception`

Read-Only:

- `description` (String)
- `filename` (String)
- `name` (String)


<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Read-Only:

- `analysis` (String)
- `application` (List of String)
- `direction` (String)
- `file_type` (List of String)
- `name` (String)


<a id="nestedatt--threat_exception"></a>
### Nested Schema for `threat_exception`

Read-Only:

- `name` (String)
- `notes` (String)

