package logs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tflogs "github.com/hashicorp/terraform-provider-aws/internal/service/logs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccLogsMetricFilter_basic(t *testing.T) {
	var mf cloudwatchlogs.MetricFilter
	resourceName := "aws_cloudwatch_log_metric_filter.test"
	logGroupResourceName := "aws_cloudwatch_log_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricFilterConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMetricFilterExists(resourceName, &mf),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.default_value", ""),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.name", "metric1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.namespace", "ns1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.unit", "None"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "pattern", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMetricFilterImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogsMetricFilter_disappears(t *testing.T) {
	var mf cloudwatchlogs.MetricFilter
	resourceName := "aws_cloudwatch_log_metric_filter.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricFilterConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricFilterExists(resourceName, &mf),
					acctest.CheckResourceDisappears(acctest.Provider, tflogs.ResourceMetricFilter(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLogsMetricFilter_Disappears_logGroup(t *testing.T) {
	var mf cloudwatchlogs.MetricFilter
	resourceName := "aws_cloudwatch_log_metric_filter.test"
	logGroupResourceName := "aws_cloudwatch_log_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricFilterConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricFilterExists(resourceName, &mf),
					acctest.CheckResourceDisappears(acctest.Provider, tflogs.ResourceGroup(), logGroupResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLogsMetricFilter_many(t *testing.T) {
	resourceName := "aws_cloudwatch_log_metric_filter.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricFilterConfig_many(rName, 15),
				Check:  testAccCheckMetricFilterManyExists(resourceName, 15),
			},
		},
	})
}

func TestAccLogsMetricFilter_update(t *testing.T) {
	var mf cloudwatchlogs.MetricFilter
	resourceName := "aws_cloudwatch_log_metric_filter.test"
	logGroupResourceName := "aws_cloudwatch_log_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricFilterConfig_allAttributes1(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMetricFilterExists(resourceName, &mf),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.default_value", "2.5"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.name", "metric1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.namespace", "ns1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.unit", "Terabytes"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.value", "3"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "pattern", "[TEST]"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMetricFilterImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccMetricFilterConfig_allAttributes2(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMetricFilterExists(resourceName, &mf),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_name", logGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.default_value", ""),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.d1", "$.d1"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.d2", "$.d2"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.dimensions.d3", "$.d3"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.name", "metric2"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.namespace", "ns2"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.unit", "Gigabits"),
					resource.TestCheckResourceAttr(resourceName, "metric_transformation.0.value", "10"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "pattern", `{ $.d1 = "OK" }`),
				),
			},
		},
	})
}

func testAccMetricFilterImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["log_group_name"] + ":" + rs.Primary.Attributes["name"], nil
	}
}

func testAccCheckMetricFilterExists(n string, v *cloudwatchlogs.MetricFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudWatch Logs Metric Filter ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LogsConn

		output, err := tflogs.FindMetricFilterByTwoPartKey(context.Background(), conn, rs.Primary.Attributes["log_group_name"], rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckMetricFilterDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).LogsConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloudwatch_log_metric_filter" {
			continue
		}

		_, err := tflogs.FindMetricFilterByTwoPartKey(context.Background(), conn, rs.Primary.Attributes["log_group_name"], rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("CloudWatch Logs Metric Filter still exists: %s", rs.Primary.ID)
	}

	return nil
}

func testAccCheckMetricFilterManyExists(basename string, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for i := 0; i < n; i++ {
			n := fmt.Sprintf("%s.%d", basename, i)
			var v cloudwatchlogs.MetricFilter

			err := testAccCheckMetricFilterExists(n, &v)(s)

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccMetricFilterConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_cloudwatch_log_metric_filter" "test" {
  name           = %[1]q
  pattern        = ""
  log_group_name = aws_cloudwatch_log_group.test.name

  metric_transformation {
    name      = "metric1"
    namespace = "ns1"
    value     = "1"
  }
}
`, rName)
}

func testAccMetricFilterConfig_many(rName string, n int) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_cloudwatch_log_metric_filter" "test" {
  count = %[2]d

  name           = "%[1]s-${count.index}"
  pattern        = "TEST"
  log_group_name = aws_cloudwatch_log_group.test.name

  metric_transformation {
    name      = "metric${count.index}"
    namespace = "ns1"
    value     = count.index
  }
}
`, rName, n)
}

func testAccMetricFilterConfig_allAttributes1(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_cloudwatch_log_metric_filter" "test" {
  name           = %[1]q
  pattern        = "[TEST] "
  log_group_name = aws_cloudwatch_log_group.test.name

  metric_transformation {
    name          = "metric1"
    namespace     = "ns1"
    unit          = "Terabytes"
    value         = "3"
    default_value = "2.5"
  }
}
`, rName)
}

func testAccMetricFilterConfig_allAttributes2(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_cloudwatch_log_metric_filter" "test" {
  name    = %[1]q
  pattern = <<EOS
    { $.d1 = "OK" }
EOS

  log_group_name = aws_cloudwatch_log_group.test.name

  metric_transformation {
    name      = "metric2"
    namespace = "ns2"
    unit      = "Gigabits"
    value     = "10"

    dimensions = {
      d1 = "$.d1"
      d2 = "$.d2"
      d3 = "$.d3"
    }
  }
}
`, rName)
}