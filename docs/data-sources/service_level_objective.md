---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "datadog_service_level_objective Data Source - terraform-provider-datadog"
subcategory: ""
description: |-
  Use this data source to retrieve information about an existing SLO for use in other resources.
---

# datadog_service_level_objective (Data Source)

Use this data source to retrieve information about an existing SLO for use in other resources.

## Example Usage

```terraform
data "datadog_service_level_objective" "test" {
  name_query = "My test SLO"
  tags_query = "foo:bar"
}

data "datadog_service_level_objective" "api_slo" {
  id = data.terraform_remote_state.api.outputs.slo
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) A SLO ID to limit the search.
- `metrics_query` (String) Filter results based on SLO numerator and denominator.
- `name_query` (String) Filter results based on SLO names.
- `tags_query` (String) Filter results based on a single SLO tag.

### Read-Only

- `description` (String) The description of the service level objective.
- `name` (String) Name of the Datadog service level objective
- `query` (List of Object) The metric query of good / total events (see [below for nested schema](#nestedatt--query))
- `target_threshold` (Number) The primary target threshold of the service level objective.
- `timeframe` (String) The primary timeframe of the service level objective.
- `type` (String) The type of the service level objective. The mapping from these types to the types found in the Datadog Web UI can be found in the Datadog API [documentation page](https://docs.datadoghq.com/api/v1/service-level-objectives/#create-a-slo-object). Available values are: `metric` and `monitor`.
- `warning_threshold` (Number) The primary warning threshold of the service level objective.

<a id="nestedatt--query"></a>
### Nested Schema for `query`

Read-Only:

- `denominator` (String)
- `numerator` (String)


