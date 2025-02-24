package apigateway_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/apigateway"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccAPIGatewayAuthorizerDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_authorizer.test"
	dataSourceName := "data.aws_api_gateway_authorizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizerDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_credentials", dataSourceName, "authorizer_credentials"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_result_ttl_in_seconds", dataSourceName, "authorizer_result_ttl_in_seconds"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_uri", dataSourceName, "authorizer_uri"),
					resource.TestCheckResourceAttrPair(resourceName, "identity_source", dataSourceName, "identity_source"),
					resource.TestCheckResourceAttrPair(resourceName, "identity_validation_expression", dataSourceName, "identity_validation_expression"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "provider_arns.#", dataSourceName, "provider_arns.#"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
				),
			},
		},
	})
}

func testAccAuthorizerDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccAuthorizerConfig_lambda(rName), `
data "aws_api_gateway_authorizer" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  authorizer_id = aws_api_gateway_authorizer.test.id
}
`)
}
