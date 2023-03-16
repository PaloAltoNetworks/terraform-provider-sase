---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_wildfire_anti_virus_profiles Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_wildfire_anti_virus_profiles (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry
- `name` (String)

### Optional

- `description` (String)
- `mlav_exception` (Attributes List) (see [below for nested schema](#nestedatt--mlav_exception))
- `packet_capture` (Boolean)
- `rules` (Attributes List) (see [below for nested schema](#nestedatt--rules))
- `threat_exception` (Attributes List) (see [below for nested schema](#nestedatt--threat_exception))

### Read-Only

- `id` (String) The object ID.
- `object_id` (String)

<a id="nestedatt--mlav_exception"></a>
### Nested Schema for `mlav_exception`

Optional:

- `description` (String)
- `filename` (String)
- `name` (String)


<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Optional:

- `analysis` (String)
- `application` (List of String)
- `direction` (String)
- `file_type` (List of String)
- `name` (String)


<a id="nestedatt--threat_exception"></a>
### Nested Schema for `threat_exception`

Optional:

- `name` (String)
- `notes` (String)

