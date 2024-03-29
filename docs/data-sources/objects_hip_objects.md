---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_hip_objects Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_objects_hip_objects (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.
- `object_id` (String) The uuid of the resource.

### Read-Only

- `anti_malware` (Attributes) The `anti_malware` parameter. (see [below for nested schema](#nestedatt--anti_malware))
- `certificate` (Attributes) The `certificate` parameter. (see [below for nested schema](#nestedatt--certificate))
- `custom_checks` (Attributes) The `custom_checks` parameter. (see [below for nested schema](#nestedatt--custom_checks))
- `data_loss_prevention` (Attributes) The `data_loss_prevention` parameter. (see [below for nested schema](#nestedatt--data_loss_prevention))
- `description` (String) The `description` parameter.
- `disk_backup` (Attributes) The `disk_backup` parameter. (see [below for nested schema](#nestedatt--disk_backup))
- `disk_encryption` (Attributes) The `disk_encryption` parameter. (see [below for nested schema](#nestedatt--disk_encryption))
- `firewall` (Attributes) The `firewall` parameter. (see [below for nested schema](#nestedatt--firewall))
- `host_info` (Attributes) The `host_info` parameter. (see [below for nested schema](#nestedatt--host_info))
- `id` (String) The object ID.
- `mobile_device` (Attributes) The `mobile_device` parameter. (see [below for nested schema](#nestedatt--mobile_device))
- `name` (String) The `name` parameter.
- `network_info` (Attributes) The `network_info` parameter. (see [below for nested schema](#nestedatt--network_info))
- `patch_management` (Attributes) The `patch_management` parameter. (see [below for nested schema](#nestedatt--patch_management))

<a id="nestedatt--anti_malware"></a>
### Nested Schema for `anti_malware`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--anti_malware--vendor))

<a id="nestedatt--anti_malware--criteria"></a>
### Nested Schema for `anti_malware.criteria`

Read-Only:

- `is_installed` (Boolean) The `is_installed` parameter.
- `last_scan_time` (Attributes) The `last_scan_time` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--last_scan_time))
- `product_version` (Attributes) The `product_version` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--product_version))
- `real_time_protection` (String) The `real_time_protection` parameter.
- `virdef_version` (Attributes) The `virdef_version` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--virdef_version))

<a id="nestedatt--anti_malware--criteria--last_scan_time"></a>
### Nested Schema for `anti_malware.criteria.last_scan_time`

Read-Only:

- `not_available` (Boolean) The `not_available` parameter.
- `not_within` (Attributes) The `not_within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--last_scan_time--not_within))
- `within` (Attributes) The `within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--last_scan_time--within))

<a id="nestedatt--anti_malware--criteria--last_scan_time--not_within"></a>
### Nested Schema for `anti_malware.criteria.last_scan_time.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `hours` (Number) The `hours` parameter.


<a id="nestedatt--anti_malware--criteria--last_scan_time--within"></a>
### Nested Schema for `anti_malware.criteria.last_scan_time.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `hours` (Number) The `hours` parameter.



<a id="nestedatt--anti_malware--criteria--product_version"></a>
### Nested Schema for `anti_malware.criteria.product_version`

Read-Only:

- `contains` (String) The `contains` parameter.
- `greater_equal` (String) The `greater_equal` parameter.
- `greater_than` (String) The `greater_than` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.
- `less_equal` (String) The `less_equal` parameter.
- `less_than` (String) The `less_than` parameter.
- `not_within` (Attributes) The `not_within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--product_version--not_within))
- `within` (Attributes) The `within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--product_version--within))

<a id="nestedatt--anti_malware--criteria--product_version--not_within"></a>
### Nested Schema for `anti_malware.criteria.product_version.within`

Read-Only:

- `versions` (Number) The `versions` parameter.


<a id="nestedatt--anti_malware--criteria--product_version--within"></a>
### Nested Schema for `anti_malware.criteria.product_version.within`

Read-Only:

- `versions` (Number) The `versions` parameter.



<a id="nestedatt--anti_malware--criteria--virdef_version"></a>
### Nested Schema for `anti_malware.criteria.virdef_version`

Read-Only:

- `not_within` (Attributes) The `not_within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--virdef_version--not_within))
- `within` (Attributes) The `within` parameter. (see [below for nested schema](#nestedatt--anti_malware--criteria--virdef_version--within))

<a id="nestedatt--anti_malware--criteria--virdef_version--not_within"></a>
### Nested Schema for `anti_malware.criteria.virdef_version.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `versions` (Number) The `versions` parameter.


<a id="nestedatt--anti_malware--criteria--virdef_version--within"></a>
### Nested Schema for `anti_malware.criteria.virdef_version.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `versions` (Number) The `versions` parameter.




<a id="nestedatt--anti_malware--vendor"></a>
### Nested Schema for `anti_malware.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.



<a id="nestedatt--certificate"></a>
### Nested Schema for `certificate`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--certificate--criteria))

<a id="nestedatt--certificate--criteria"></a>
### Nested Schema for `certificate.criteria`

Read-Only:

- `certificate_attributes` (Attributes List) The `certificate_attributes` parameter. (see [below for nested schema](#nestedatt--certificate--criteria--certificate_attributes))
- `certificate_profile` (String) The `certificate_profile` parameter.

<a id="nestedatt--certificate--criteria--certificate_attributes"></a>
### Nested Schema for `certificate.criteria.certificate_attributes`

Read-Only:

- `name` (String) The `name` parameter.
- `value` (String) The `value` parameter.




<a id="nestedatt--custom_checks"></a>
### Nested Schema for `custom_checks`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria))

<a id="nestedatt--custom_checks--criteria"></a>
### Nested Schema for `custom_checks.criteria`

Read-Only:

- `plist` (Attributes List) The `plist` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria--plist))
- `process_list` (Attributes List) The `process_list` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria--process_list))
- `registry_key` (Attributes List) The `registry_key` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria--registry_key))

<a id="nestedatt--custom_checks--criteria--plist"></a>
### Nested Schema for `custom_checks.criteria.plist`

Read-Only:

- `key` (Attributes List) The `key` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria--plist--key))
- `name` (String) The `name` parameter.
- `negate` (Boolean) The `negate` parameter.

<a id="nestedatt--custom_checks--criteria--plist--key"></a>
### Nested Schema for `custom_checks.criteria.plist.negate`

Read-Only:

- `name` (String) The `name` parameter.
- `negate` (Boolean) The `negate` parameter.
- `value` (String) The `value` parameter.



<a id="nestedatt--custom_checks--criteria--process_list"></a>
### Nested Schema for `custom_checks.criteria.process_list`

Read-Only:

- `name` (String) The `name` parameter.
- `running` (Boolean) The `running` parameter.


<a id="nestedatt--custom_checks--criteria--registry_key"></a>
### Nested Schema for `custom_checks.criteria.registry_key`

Read-Only:

- `default_value_data` (String) The `default_value_data` parameter.
- `name` (String) The `name` parameter.
- `negate` (Boolean) The `negate` parameter.
- `registry_value` (Attributes List) The `registry_value` parameter. (see [below for nested schema](#nestedatt--custom_checks--criteria--registry_key--registry_value))

<a id="nestedatt--custom_checks--criteria--registry_key--registry_value"></a>
### Nested Schema for `custom_checks.criteria.registry_key.registry_value`

Read-Only:

- `name` (String) The `name` parameter.
- `negate` (Boolean) The `negate` parameter.
- `value_data` (String) The `value_data` parameter.





<a id="nestedatt--data_loss_prevention"></a>
### Nested Schema for `data_loss_prevention`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--data_loss_prevention--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--data_loss_prevention--vendor))

<a id="nestedatt--data_loss_prevention--criteria"></a>
### Nested Schema for `data_loss_prevention.criteria`

Read-Only:

- `is_enabled` (String) The `is_enabled` parameter.
- `is_installed` (Boolean) The `is_installed` parameter.


<a id="nestedatt--data_loss_prevention--vendor"></a>
### Nested Schema for `data_loss_prevention.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.



<a id="nestedatt--disk_backup"></a>
### Nested Schema for `disk_backup`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--disk_backup--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--disk_backup--vendor))

<a id="nestedatt--disk_backup--criteria"></a>
### Nested Schema for `disk_backup.criteria`

Read-Only:

- `is_installed` (Boolean) The `is_installed` parameter.
- `last_backup_time` (Attributes) The `last_backup_time` parameter. (see [below for nested schema](#nestedatt--disk_backup--criteria--last_backup_time))

<a id="nestedatt--disk_backup--criteria--last_backup_time"></a>
### Nested Schema for `disk_backup.criteria.last_backup_time`

Read-Only:

- `not_available` (Boolean) The `not_available` parameter.
- `not_within` (Attributes) The `not_within` parameter. (see [below for nested schema](#nestedatt--disk_backup--criteria--last_backup_time--not_within))
- `within` (Attributes) The `within` parameter. (see [below for nested schema](#nestedatt--disk_backup--criteria--last_backup_time--within))

<a id="nestedatt--disk_backup--criteria--last_backup_time--not_within"></a>
### Nested Schema for `disk_backup.criteria.last_backup_time.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `hours` (Number) The `hours` parameter.


<a id="nestedatt--disk_backup--criteria--last_backup_time--within"></a>
### Nested Schema for `disk_backup.criteria.last_backup_time.within`

Read-Only:

- `days` (Number) The `days` parameter.
- `hours` (Number) The `hours` parameter.




<a id="nestedatt--disk_backup--vendor"></a>
### Nested Schema for `disk_backup.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.



<a id="nestedatt--disk_encryption"></a>
### Nested Schema for `disk_encryption`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--disk_encryption--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--disk_encryption--vendor))

<a id="nestedatt--disk_encryption--criteria"></a>
### Nested Schema for `disk_encryption.criteria`

Read-Only:

- `encrypted_locations` (Attributes List) The `encrypted_locations` parameter. (see [below for nested schema](#nestedatt--disk_encryption--criteria--encrypted_locations))
- `is_installed` (Boolean) The `is_installed` parameter.

<a id="nestedatt--disk_encryption--criteria--encrypted_locations"></a>
### Nested Schema for `disk_encryption.criteria.encrypted_locations`

Read-Only:

- `encryption_state` (Attributes) The `encryption_state` parameter. (see [below for nested schema](#nestedatt--disk_encryption--criteria--encrypted_locations--encryption_state))
- `name` (String) The `name` parameter.

<a id="nestedatt--disk_encryption--criteria--encrypted_locations--encryption_state"></a>
### Nested Schema for `disk_encryption.criteria.encrypted_locations.name`

Read-Only:

- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.




<a id="nestedatt--disk_encryption--vendor"></a>
### Nested Schema for `disk_encryption.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.



<a id="nestedatt--firewall"></a>
### Nested Schema for `firewall`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--firewall--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--firewall--vendor))

<a id="nestedatt--firewall--criteria"></a>
### Nested Schema for `firewall.criteria`

Read-Only:

- `is_enabled` (String) The `is_enabled` parameter.
- `is_installed` (Boolean) The `is_installed` parameter.


<a id="nestedatt--firewall--vendor"></a>
### Nested Schema for `firewall.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.



<a id="nestedatt--host_info"></a>
### Nested Schema for `host_info`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--host_info--criteria))

<a id="nestedatt--host_info--criteria"></a>
### Nested Schema for `host_info.criteria`

Read-Only:

- `client_version` (Attributes) The `client_version` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--client_version))
- `domain` (Attributes) The `domain` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--domain))
- `host_id` (Attributes) The `host_id` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--host_id))
- `host_name` (Attributes) The `host_name` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--host_name))
- `managed` (Boolean) The `managed` parameter.
- `os` (Attributes) The `os` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--os))
- `serial_number` (Attributes) The `serial_number` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--serial_number))

<a id="nestedatt--host_info--criteria--client_version"></a>
### Nested Schema for `host_info.criteria.client_version`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--host_info--criteria--domain"></a>
### Nested Schema for `host_info.criteria.domain`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--host_info--criteria--host_id"></a>
### Nested Schema for `host_info.criteria.host_id`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--host_info--criteria--host_name"></a>
### Nested Schema for `host_info.criteria.host_name`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--host_info--criteria--os"></a>
### Nested Schema for `host_info.criteria.os`

Read-Only:

- `contains` (Attributes) The `contains` parameter. (see [below for nested schema](#nestedatt--host_info--criteria--os--contains))

<a id="nestedatt--host_info--criteria--os--contains"></a>
### Nested Schema for `host_info.criteria.os.contains`

Read-Only:

- `apple` (String) The `apple` parameter.
- `google` (String) The `google` parameter.
- `linux` (String) The `linux` parameter.
- `microsoft` (String) The `microsoft` parameter.
- `other` (String) The `other` parameter.



<a id="nestedatt--host_info--criteria--serial_number"></a>
### Nested Schema for `host_info.criteria.serial_number`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.




<a id="nestedatt--mobile_device"></a>
### Nested Schema for `mobile_device`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria))

<a id="nestedatt--mobile_device--criteria"></a>
### Nested Schema for `mobile_device.criteria`

Read-Only:

- `applications` (Attributes) The `applications` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--applications))
- `disk_encrypted` (Boolean) The `disk_encrypted` parameter.
- `imei` (Attributes) The `imei` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--imei))
- `jailbroken` (Boolean) The `jailbroken` parameter.
- `last_checkin_time` (Attributes) The `last_checkin_time` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--last_checkin_time))
- `model` (Attributes) The `model` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--model))
- `passcode_set` (Boolean) The `passcode_set` parameter.
- `phone_number` (Attributes) The `phone_number` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--phone_number))
- `tag` (Attributes) The `tag` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--tag))

<a id="nestedatt--mobile_device--criteria--applications"></a>
### Nested Schema for `mobile_device.criteria.applications`

Read-Only:

- `has_malware` (Attributes) The `has_malware` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--applications--has_malware))
- `has_unmanaged_app` (Boolean) The `has_unmanaged_app` parameter.
- `includes` (Attributes List) The `includes` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--applications--includes))

<a id="nestedatt--mobile_device--criteria--applications--has_malware"></a>
### Nested Schema for `mobile_device.criteria.applications.includes`

Read-Only:

- `no` (Boolean) The `no` parameter.
- `yes` (Attributes) The `yes` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--applications--includes--yes))

<a id="nestedatt--mobile_device--criteria--applications--includes--yes"></a>
### Nested Schema for `mobile_device.criteria.applications.includes.yes`

Read-Only:

- `excludes` (Attributes List) The `excludes` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--applications--includes--yes--excludes))

<a id="nestedatt--mobile_device--criteria--applications--includes--yes--excludes"></a>
### Nested Schema for `mobile_device.criteria.applications.includes.yes.excludes`

Read-Only:

- `hash` (String) The `hash` parameter.
- `name` (String) The `name` parameter.
- `package` (String) The `package` parameter.




<a id="nestedatt--mobile_device--criteria--applications--includes"></a>
### Nested Schema for `mobile_device.criteria.applications.includes`

Read-Only:

- `hash` (String) The `hash` parameter.
- `name` (String) The `name` parameter.
- `package` (String) The `package` parameter.



<a id="nestedatt--mobile_device--criteria--imei"></a>
### Nested Schema for `mobile_device.criteria.imei`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--mobile_device--criteria--last_checkin_time"></a>
### Nested Schema for `mobile_device.criteria.last_checkin_time`

Read-Only:

- `not_within` (Attributes) The `not_within` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--last_checkin_time--not_within))
- `within` (Attributes) The `within` parameter. (see [below for nested schema](#nestedatt--mobile_device--criteria--last_checkin_time--within))

<a id="nestedatt--mobile_device--criteria--last_checkin_time--not_within"></a>
### Nested Schema for `mobile_device.criteria.last_checkin_time.within`

Read-Only:

- `days` (Number) The `days` parameter.


<a id="nestedatt--mobile_device--criteria--last_checkin_time--within"></a>
### Nested Schema for `mobile_device.criteria.last_checkin_time.within`

Read-Only:

- `days` (Number) The `days` parameter.



<a id="nestedatt--mobile_device--criteria--model"></a>
### Nested Schema for `mobile_device.criteria.model`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--mobile_device--criteria--phone_number"></a>
### Nested Schema for `mobile_device.criteria.phone_number`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.


<a id="nestedatt--mobile_device--criteria--tag"></a>
### Nested Schema for `mobile_device.criteria.tag`

Read-Only:

- `contains` (String) The `contains` parameter.
- `is` (String) The `is` parameter.
- `is_not` (String) The `is_not` parameter.




<a id="nestedatt--network_info"></a>
### Nested Schema for `network_info`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--network_info--criteria))

<a id="nestedatt--network_info--criteria"></a>
### Nested Schema for `network_info.criteria`

Read-Only:

- `network` (Attributes) The `network` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network))

<a id="nestedatt--network_info--criteria--network"></a>
### Nested Schema for `network_info.criteria.network`

Read-Only:

- `is` (Attributes) The `is` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is))
- `is_not` (Attributes) The `is_not` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is_not))

<a id="nestedatt--network_info--criteria--network--is"></a>
### Nested Schema for `network_info.criteria.network.is_not`

Read-Only:

- `mobile` (Attributes) The `mobile` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is_not--mobile))
- `unknown` (Boolean) The `unknown` parameter.
- `wifi` (Attributes) The `wifi` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is_not--wifi))

<a id="nestedatt--network_info--criteria--network--is_not--mobile"></a>
### Nested Schema for `network_info.criteria.network.is_not.mobile`

Read-Only:

- `carrier` (String) The `carrier` parameter.


<a id="nestedatt--network_info--criteria--network--is_not--wifi"></a>
### Nested Schema for `network_info.criteria.network.is_not.wifi`

Read-Only:

- `ssid` (String) The `ssid` parameter.



<a id="nestedatt--network_info--criteria--network--is_not"></a>
### Nested Schema for `network_info.criteria.network.is_not`

Read-Only:

- `ethernet` (Boolean) The `ethernet` parameter.
- `mobile` (Attributes) The `mobile` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is_not--mobile))
- `unknown` (Boolean) The `unknown` parameter.
- `wifi` (Attributes) The `wifi` parameter. (see [below for nested schema](#nestedatt--network_info--criteria--network--is_not--wifi))

<a id="nestedatt--network_info--criteria--network--is_not--mobile"></a>
### Nested Schema for `network_info.criteria.network.is_not.mobile`

Read-Only:

- `carrier` (String) The `carrier` parameter.


<a id="nestedatt--network_info--criteria--network--is_not--wifi"></a>
### Nested Schema for `network_info.criteria.network.is_not.wifi`

Read-Only:

- `ssid` (String) The `ssid` parameter.






<a id="nestedatt--patch_management"></a>
### Nested Schema for `patch_management`

Read-Only:

- `criteria` (Attributes) The `criteria` parameter. (see [below for nested schema](#nestedatt--patch_management--criteria))
- `exclude_vendor` (Boolean) The `exclude_vendor` parameter.
- `vendor` (Attributes List) The `vendor` parameter. (see [below for nested schema](#nestedatt--patch_management--vendor))

<a id="nestedatt--patch_management--criteria"></a>
### Nested Schema for `patch_management.criteria`

Read-Only:

- `is_enabled` (String) The `is_enabled` parameter.
- `is_installed` (Boolean) The `is_installed` parameter.
- `missing_patches` (Attributes) The `missing_patches` parameter. (see [below for nested schema](#nestedatt--patch_management--criteria--missing_patches))

<a id="nestedatt--patch_management--criteria--missing_patches"></a>
### Nested Schema for `patch_management.criteria.missing_patches`

Read-Only:

- `check` (String) The `check` parameter.
- `patches` (List of String) The `patches` parameter.
- `severity` (Attributes) The `severity` parameter. (see [below for nested schema](#nestedatt--patch_management--criteria--missing_patches--severity))

<a id="nestedatt--patch_management--criteria--missing_patches--severity"></a>
### Nested Schema for `patch_management.criteria.missing_patches.severity`

Read-Only:

- `greater_equal` (Number) The `greater_equal` parameter.
- `greater_than` (Number) The `greater_than` parameter.
- `is` (Number) The `is` parameter.
- `is_not` (Number) The `is_not` parameter.
- `less_equal` (Number) The `less_equal` parameter.
- `less_than` (Number) The `less_than` parameter.




<a id="nestedatt--patch_management--vendor"></a>
### Nested Schema for `patch_management.vendor`

Read-Only:

- `name` (String) The `name` parameter.
- `product` (List of String) The `product` parameter.


