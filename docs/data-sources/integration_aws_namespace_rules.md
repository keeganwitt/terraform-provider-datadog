---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "datadog_integration_aws_namespace_rules Data Source - terraform-provider-datadog"
subcategory: ""
description: |-
  Provides a Datadog AWS Integration Namespace Rules data source. This can be used to retrieve all available namespace rules for a Datadog-AWS integration.
---

# datadog_integration_aws_namespace_rules (Data Source)

Provides a Datadog AWS Integration Namespace Rules data source. This can be used to retrieve all available namespace rules for a Datadog-AWS integration.

## Example Usage

```terraform
data "datadog_integration_aws_namespace_rules" "rules" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `namespace_rules` (List of String) The list of available namespace rules for a Datadog-AWS integration.


