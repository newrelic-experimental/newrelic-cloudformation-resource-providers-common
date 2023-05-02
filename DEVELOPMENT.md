## Docs
- [AWS CloudFormation CLI](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/cloudformation/index.html)
- [CFN CLI](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/resource-type-cli.html)
- [Registering an extension privately (step 1)](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/resource-type-register.html)
- [Publishing an extension publicly (step 2)](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/publish-extension.html)
## Troubleshooting
- Error log
- `Debug` log level for mustache substitution
- Validate mutation using the [Explorer](https://api.newrelic.com/graphiql)

## Building
- Install Docker
- Install Golang
- [Install the CloudFormation Command Line Interface (CFN-CLI)](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html)
- `make clean build `

## Testing
- Start Docker
- Activate `cfn-cli`, usually `source ~/.virtenv/aws-cfn-cli/bin/activate`
- Make the project `make clean build`
- Copy `TypeConfigurationLive.json` to `bin/`
- Start the container `sam local start-lambda --warm-containers eager`
- Or
```bash
rm rpdk.log ; clear ; make clean build ; cp TypeConfigurationLive.json bin/ ; sam local start-lambda --warm-containers eager 
```
- Run the [CloudFormation Contract Tests](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/contract-tests.html) `cfn test`
- NOTE: don't use a `duplicate` configured `inputs_x_create.json` file in the same `cfn test` run as a normal `create`, things break badly. Only one test type at a time.

## Publishing
```bash
# Remove TypeConfigurationLive.json from bin/ to avoid leaking the key!
rm bin/TypeConfigurationLive.json
# Double check the resulting zip file for security leaks!
cfn submit --dry-run
# Send the resource to AWS, the result is a private resource. NOTE: make clean build clears credentials from bin/
rm newrelic-cloudformation-*.zip  ; make clean build  ; cfn submit --set-default  --region us-east-1  --role-arn arn:aws:iam::830139413159:role/custom-resource-cloudformation-role
# Set the configuration on the private resource
aws cloudformation set-type-configuration --region us-east-1 --type RESOURCE --type-name NewRelic::Observability::<TYPE> --configuration-alias default --configuration "{  \"APIKey\": \"<API_KEY>\",  \"AccountID\": 
\"<ACCOUNT_ID>\",  \"Endpoint\": \"https://api.newrelic.com/graphql\"}" 
# Test the private resource with the sample template to ensure it works
aws cloudformation deploy  --force-upload --disable-rollback --region us-east-1 --template-file template-examples-live/live.yml --stack-name test-stack-workloads
# Tell AWS to run the Contract Tests, required for going public
aws cloudformation test-type --region us-east-1 --log-delivery-bucket newrelic--cloudformation--custom--resources --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-workload
# Check the result
aws cloudformation describe-type --region us-east-1  --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-workloads
# Also the logs are in CloudWatch. They end in .zip but are really gunzip so
# gunzip -S .zip <file>
# Publish the extension publicly AFTER pushing the final version to GitHub AND generating/tagging a release
# IMPORTANT!
#   --public-version-number (string)
#   The version number to assign to this version of the extension.
#   Use the following format, and adhere to semantic versioning when assigning a version number to your extension:
#     MAJOR.MINOR.PATCH
#   For more information, see Semantic Versioning 2.0.0 .
#   If you don’t specify a version number, CloudFormation increments the version number by one minor version release.
#   You cannot specify a version number the first time you publish a type. CloudFormation automatically sets the first version number to be 1.0.0 .
#
# It's a good idea to not publish until everything is ready AND version 1.0.0 is release in Git!
# KEEP IT ALL IN-SYNC!
#
# Git
#
aws cloudformation publish-type --region us-east-1  --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-workloads
aws cloudformation describe-type --region us-east-1  --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-workloads
```
### Notes
- If you see an `Exception: Could not assume specified role...` error message in the test log then you probably have a `cfn` generated role and they're not correct. 
  - Correct IAM Trust Relationship
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "cloudformation.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        },
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "resources.cloudformation.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```
  - Correct IAM Permissions Policies- `AWSCloudFormationFullAccess` 

- If a test fails and the Progress Event _looks_ right in the console check `rpdk.log`. If that doesn't help suspect that an error is being returned to the sdk in addition to the Progress Event.


## Helpful links
- [CloudFormation CLI User Guide](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html)
- [New Relic GraphQL Explorer](https://api.newrelic.com/graphiql) 
