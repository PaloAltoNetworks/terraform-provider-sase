---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_hip_profiles Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_objects_hip_profiles (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `match` (String) The `match` parameter. String length must be at most 2048.
- `name` (String) The `name` parameter. String length must be at most 31.

### Optional

- `description` (String) The `description` parameter. String length must be between 0 and 255.

### Read-Only

- `id` (String) The object ID.
- `object_id` (String) The `object_id` parameter.


