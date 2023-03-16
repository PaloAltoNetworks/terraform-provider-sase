---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_applications Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_objects_applications (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `category` (String)
- `folder` (String) The folder of the entry
- `name` (String)
- `risk` (Number)
- `subcategory` (String)
- `technology` (String)

### Optional

- `able_to_transfer_file` (Boolean)
- `alg_disable_capability` (String)
- `consume_big_bandwidth` (Boolean)
- `data_ident` (Boolean)
- `default` (Attributes) (see [below for nested schema](#nestedatt--default))
- `description` (String)
- `evasive_behavior` (Boolean)
- `file_type_ident` (Boolean)
- `has_known_vulnerability` (Boolean)
- `no_appid_caching` (Boolean)
- `parent_app` (String)
- `pervasive_use` (Boolean)
- `prone_to_misuse` (Boolean)
- `signature` (Attributes List) (see [below for nested schema](#nestedatt--signature))
- `tcp_half_closed_timeout` (Number)
- `tcp_time_wait_timeout` (Number)
- `tcp_timeout` (Number)
- `timeout` (Number)
- `tunnel_applications` (Boolean)
- `tunnel_other_application` (Boolean)
- `udp_timeout` (Number)
- `used_by_malware` (Boolean)
- `virus_ident` (Boolean)

### Read-Only

- `id` (String) The object ID.
- `object_id` (String)

<a id="nestedatt--default"></a>
### Nested Schema for `default`

Optional:

- `ident_by_icmp6_type` (Attributes) (see [below for nested schema](#nestedatt--default--ident_by_icmp6_type))
- `ident_by_icmp_type` (Attributes) (see [below for nested schema](#nestedatt--default--ident_by_icmp_type))
- `ident_by_ip_protocol` (String)
- `port` (List of String)

<a id="nestedatt--default--ident_by_icmp6_type"></a>
### Nested Schema for `default.ident_by_icmp6_type`

Required:

- `type` (String)

Optional:

- `code` (String)


<a id="nestedatt--default--ident_by_icmp_type"></a>
### Nested Schema for `default.ident_by_icmp_type`

Required:

- `type` (String)

Optional:

- `code` (String)



<a id="nestedatt--signature"></a>
### Nested Schema for `signature`

Required:

- `name` (String)

Optional:

- `and_condition` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition))
- `comment` (String)
- `order_free` (Boolean)
- `scope` (String)

<a id="nestedatt--signature--and_condition"></a>
### Nested Schema for `signature.and_condition`

Required:

- `name` (String)

Optional:

- `or_condition` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition))

<a id="nestedatt--signature--and_condition--or_condition"></a>
### Nested Schema for `signature.and_condition.or_condition`

Required:

- `name` (String)
- `operator` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator))

<a id="nestedatt--signature--and_condition--or_condition--operator"></a>
### Nested Schema for `signature.and_condition.or_condition.operator`

Optional:

- `equal_to` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--equal_to))
- `greater_than` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--greater_than))
- `less_than` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--less_than))
- `pattern_match` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--pattern_match))

<a id="nestedatt--signature--and_condition--or_condition--operator--equal_to"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.equal_to`

Required:

- `context` (String)
- `value` (String)

Optional:

- `mask` (String)
- `position` (String)


<a id="nestedatt--signature--and_condition--or_condition--operator--greater_than"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.greater_than`

Required:

- `context` (String)
- `value` (Number)

Optional:

- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--greater_than--qualifier))

<a id="nestedatt--signature--and_condition--or_condition--operator--greater_than--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.greater_than.qualifier`

Required:

- `name` (String)
- `value` (String)



<a id="nestedatt--signature--and_condition--or_condition--operator--less_than"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.less_than`

Required:

- `context` (String)
- `value` (Number)

Optional:

- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--less_than--qualifier))

<a id="nestedatt--signature--and_condition--or_condition--operator--less_than--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.less_than.qualifier`

Required:

- `name` (String)
- `value` (String)



<a id="nestedatt--signature--and_condition--or_condition--operator--pattern_match"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.pattern_match`

Required:

- `context` (String)
- `pattern` (String)

Optional:

- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--pattern_match--qualifier))

<a id="nestedatt--signature--and_condition--or_condition--operator--pattern_match--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.pattern_match.qualifier`

Required:

- `name` (String)
- `value` (String)

