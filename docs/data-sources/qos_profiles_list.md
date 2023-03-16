---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_qos_profiles_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_qos_profiles_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry

### Optional

- `limit` (Number) The max count in result entry (count per page)
- `name` (String) The name of the entry
- `offset` (Number) The offset of the result entry

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `total` (Number)

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `aggregate_bandwidth` (Attributes) (see [below for nested schema](#nestedatt--data--aggregate_bandwidth))
- `class_bandwidth_type` (Attributes) (see [below for nested schema](#nestedatt--data--class_bandwidth_type))
- `name` (String)
- `object_id` (String)

<a id="nestedatt--data--aggregate_bandwidth"></a>
### Nested Schema for `data.aggregate_bandwidth`

Read-Only:

- `egress_guaranteed` (Number)
- `egress_max` (Number)


<a id="nestedatt--data--class_bandwidth_type"></a>
### Nested Schema for `data.class_bandwidth_type`

Read-Only:

- `mbps` (Attributes) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--mbps))
- `percentage` (Attributes) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--percentage))

<a id="nestedatt--data--class_bandwidth_type--mbps"></a>
### Nested Schema for `data.class_bandwidth_type.mbps`

Read-Only:

- `class` (Attributes List) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--mbps--class))

<a id="nestedatt--data--class_bandwidth_type--mbps--class"></a>
### Nested Schema for `data.class_bandwidth_type.mbps.class`

Read-Only:

- `class_bandwidth` (Attributes) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--mbps--class--class_bandwidth))
- `name` (String)
- `priority` (String)

<a id="nestedatt--data--class_bandwidth_type--mbps--class--class_bandwidth"></a>
### Nested Schema for `data.class_bandwidth_type.mbps.class.class_bandwidth`

Read-Only:

- `egress_guaranteed` (Number)
- `egress_max` (Number)




<a id="nestedatt--data--class_bandwidth_type--percentage"></a>
### Nested Schema for `data.class_bandwidth_type.percentage`

Read-Only:

- `class` (Attributes List) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--percentage--class))

<a id="nestedatt--data--class_bandwidth_type--percentage--class"></a>
### Nested Schema for `data.class_bandwidth_type.percentage.class`

Read-Only:

- `class_bandwidth` (Attributes) (see [below for nested schema](#nestedatt--data--class_bandwidth_type--percentage--class--class_bandwidth))
- `name` (String)
- `priority` (String)

<a id="nestedatt--data--class_bandwidth_type--percentage--class--class_bandwidth"></a>
### Nested Schema for `data.class_bandwidth_type.percentage.class.class_bandwidth`

Read-Only:

- `egress_guaranteed` (Number)
- `egress_max` (Number)

