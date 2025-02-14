package test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatadogHostsDatasource(t *testing.T) {
	t.Parallel()
	_, _, accProviders := testAccFrameworkMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "datadog_hosts" "test" {
						include_muted_hosts_data = true
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "total_matching"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "total_returned"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.host_name"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.id"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.name"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.up"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.meta.agent_version"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.meta.cpu_cores"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.meta.gohai"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.meta.machine"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.metrics.cpu"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.metrics.iowait"),
					resource.TestCheckResourceAttrSet("data.datadog_hosts.test", "host_list.0.metrics.load"),
				),
			},
		},
	})
}
