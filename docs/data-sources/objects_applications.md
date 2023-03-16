---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_applications Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_objects_applications (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The uuid of the resource

### Read-Only

- `able_to_transfer_file` (Boolean)
- `alg_disable_capability` (String)
- `category` (String)
- `consume_big_bandwidth` (Boolean)
- `data_ident` (Boolean)
- `default` (Attributes) (see [below for nested schema](#nestedatt--default))
- `description` (String)
- `evasive_behavior` (Boolean)
- `file_type_ident` (Boolean)
- `has_known_vulnerability` (Boolean)
- `id` (String) The object ID.
- `name` (String)
- `no_appid_caching` (Boolean)
- `parent_app` (String)
- `pervasive_use` (Boolean)
- `prone_to_misuse` (Boolean)
- `risk` (Number)
- `signature` (Attributes List) (see [below for nested schema](#nestedatt--signature))
- `subcategory` (String)
- `tcp_half_closed_timeout` (Number)
- `tcp_time_wait_timeout` (Number)
- `tcp_timeout` (Number)
- `technology` (String)
- `timeout` (Number)
- `tunnel_applications` (Boolean)
- `tunnel_other_application` (Boolean)
- `udp_timeout` (Number)
- `used_by_malware` (Boolean)
- `virus_ident` (Boolean)

<a id="nestedatt--default"></a>
### Nested Schema for `default`

Read-Only:

- `ident_by_icmp6_type` (Attributes) (see [below for nested schema](#nestedatt--default--ident_by_icmp6_type))
- `ident_by_icmp_type` (Attributes) (see [below for nested schema](#nestedatt--default--ident_by_icmp_type))
- `ident_by_ip_protocol` (String)
- `port` (List of String)

<a id="nestedatt--default--ident_by_icmp6_type"></a>
### Nested Schema for `default.ident_by_icmp6_type`

Read-Only:

- `code` (String)
- `type` (String)


<a id="nestedatt--default--ident_by_icmp_type"></a>
### Nested Schema for `default.ident_by_icmp_type`

Read-Only:

- `code` (String)
- `type` (String)



<a id="nestedatt--signature"></a>
### Nested Schema for `signature`

Read-Only:

- `and_condition` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition))
- `comment` (String)
- `name` (String)
- `order_free` (Boolean)
- `scope` (String)

<a id="nestedatt--signature--and_condition"></a>
### Nested Schema for `signature.and_condition`

Read-Only:

- `name` (String)
- `or_condition` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition))

<a id="nestedatt--signature--and_condition--or_condition"></a>
### Nested Schema for `signature.and_condition.or_condition`

Read-Only:

- `name` (String)
- `operator` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator))

<a id="nestedatt--signature--and_condition--or_condition--operator"></a>
### Nested Schema for `signature.and_condition.or_condition.operator`

Read-Only:

- `equal_to` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--equal_to))
- `greater_than` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--greater_than))
- `less_than` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--less_than))
- `pattern_match` (Attributes) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--pattern_match))

<a id="nestedatt--signature--and_condition--or_condition--operator--equal_to"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.equal_to`

Read-Only:

- `context` (String)
- `mask` (String)
- `position` (String)
- `value` (String)


<a id="nestedatt--signature--and_condition--or_condition--operator--greater_than"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.greater_than`

Read-Only:

- `context` (String)
- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--greater_than--qualifier))
- `value` (Number)

<a id="nestedatt--signature--and_condition--or_condition--operator--greater_than--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.greater_than.value`

Read-Only:

- `name` (String)
- `value` (String)



<a id="nestedatt--signature--and_condition--or_condition--operator--less_than"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.less_than`

Read-Only:

- `context` (String)
- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--less_than--qualifier))
- `value` (Number)

<a id="nestedatt--signature--and_condition--or_condition--operator--less_than--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.less_than.value`

Read-Only:

- `name` (String)
- `value` (String)



<a id="nestedatt--signature--and_condition--or_condition--operator--pattern_match"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.pattern_match`

Read-Only:

- `context` (String)
- `pattern` (String)
- `qualifier` (Attributes List) (see [below for nested schema](#nestedatt--signature--and_condition--or_condition--operator--pattern_match--qualifier))

<a id="nestedatt--signature--and_condition--or_condition--operator--pattern_match--qualifier"></a>
### Nested Schema for `signature.and_condition.or_condition.operator.pattern_match.qualifier`

Read-Only:

- `name` (String)
- `value` (String)

