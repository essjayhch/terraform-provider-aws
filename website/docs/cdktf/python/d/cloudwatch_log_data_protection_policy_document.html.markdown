---
subcategory: "CloudWatch Logs"
layout: "aws"
page_title: "AWS: aws_cloudwatch_log_data_protection_policy_document"
description: |-
  Generates a CloudWatch Log Group Data Protection Policy document in JSON format
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_cloudwatch_log_data_protection_policy_document

Generates a CloudWatch Log Group Data Protection Policy document in JSON format for use with the `aws_cloudwatch_log_data_protection_policy` resource.

-> For more information about data protection policies, see the [Help protect sensitive log data with masking](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/mask-sensitive-log-data.html).

## Example Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Token, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.cloudwatch_log_data_protection_policy import CloudwatchLogDataProtectionPolicy
from imports.aws.data_aws_cloudwatch_log_data_protection_policy_document import DataAwsCloudwatchLogDataProtectionPolicyDocument
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        example = DataAwsCloudwatchLogDataProtectionPolicyDocument(self, "example",
            name="Example",
            statement=[DataAwsCloudwatchLogDataProtectionPolicyDocumentStatement(
                data_identifiers=["arn:aws:dataprotection::aws:data-identifier/EmailAddress", "arn:aws:dataprotection::aws:data-identifier/DriversLicense-US"
                ],
                operation=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperation(
                    audit=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationAudit(
                        findings_destination=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationAuditFindingsDestination(
                            cloudwatch_logs=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationAuditFindingsDestinationCloudwatchLogs(
                                log_group=audit.name
                            ),
                            firehose=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationAuditFindingsDestinationFirehose(
                                delivery_stream=Token.as_string(aws_kinesis_firehose_delivery_stream_audit.name)
                            ),
                            s3=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationAuditFindingsDestinationS3(
                                bucket=Token.as_string(aws_s3_bucket_audit.bucket)
                            )
                        )
                    )
                ),
                sid="Audit"
            ), DataAwsCloudwatchLogDataProtectionPolicyDocumentStatement(
                data_identifiers=["arn:aws:dataprotection::aws:data-identifier/EmailAddress", "arn:aws:dataprotection::aws:data-identifier/DriversLicense-US"
                ],
                operation=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperation(
                    deidentify=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationDeidentify(
                        mask_config=DataAwsCloudwatchLogDataProtectionPolicyDocumentStatementOperationDeidentifyMaskConfig()
                    )
                ),
                sid="Deidentify"
            )
            ]
        )
        aws_cloudwatch_log_data_protection_policy_example =
        CloudwatchLogDataProtectionPolicy(self, "example_1",
            log_group_name=Token.as_string(aws_cloudwatch_log_group_example.name),
            policy_document=Token.as_string(example.json)
        )
        # This allows the Terraform resource name to match the original name. You can remove the call if you don't need them to match.
        aws_cloudwatch_log_data_protection_policy_example.override_logical_id("example")
```

## Argument Reference

The following arguments are required:

* `name` - (Required) The name of the data protection policy document.
* `statement` - (Required) Configures the data protection policy.

-> There must be exactly two statements: the first with an `audit` operation, and the second with a `deidentify` operation.

The following arguments are optional:

* `description` - (Optional)
* `version` - (Optional)
* `configuration` - (Optional)

### statement Configuration Block

* `data_identifiers` - (Required) Set of at least 1 sensitive data identifiers that you want to mask. Read more in [Types of data that you can protect](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/protect-sensitive-log-data-types.html).
* `operation` - (Required) Configures the data protection operation applied by this statement.
* `sid` - (Optional) Name of this statement.

#### operation Configuration Block

* `audit` - (Optional) Configures the detection of sensitive data.
* `deidentify` - (Optional) Configures the masking of sensitive data.

-> Every policy statement must specify exactly one operation.

##### audit Configuration Block

* `findings_destination` - (Required) Configures destinations to send audit findings to.

##### findings_destination Configuration Block

* `cloudwatch_logs` - (Optional) Configures CloudWatch Logs as a findings destination.
* `firehose` - (Optional) Configures Kinesis Firehose as a findings destination.
* `s3` - (Optional) Configures S3 as a findings destination.

###### cloudwatch_logs Configuration Block

* `log_group` - (Required) Name of the CloudWatch Log Group to send findings to.

###### firehose Configuration Block

* `delivery_stream` - (Required) Name of the Kinesis Firehose Delivery Stream to send findings to.

###### s3 Configuration Block

* `bucket` - (Required) Name of the S3 Bucket to send findings to.

##### deidentify Configuration Block

* `mask_config` - (Required) An empty object that configures masking.

### configuration Configuration Block

* `custom_data_identifier` - (Optional) Configures custom regular expressions to detect sensitive data. Read more in [Custom data identifiers](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL-custom-data-identifiers.html).

#### custom_data_identifier Configuration Block

* `name` - (Required) Name of the custom data idenfitier
* `regex` - (Required) Regular expression to match sensitive data

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `json` - Standard JSON policy document rendered based on the arguments above.

<!-- cache-key: cdktf-0.20.8 input-4c58e855d6626ca8f266da57b7d0263777d8947c611506d608c213d86d1115f7 -->