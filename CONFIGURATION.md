## Description
This Cloud Formation Custom Resource provides a CRUDL interface to the New Relic [NerdGraph (GraphQL) Workloads API](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-workloads-api-tutorials/) for Cloud Formation stacks.

## Prerequisites
This document assumes familiarity with using CloudFormation Public extensions in CloudFormation templates. If you are not familiar with this [start here](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/registry-public.html)

## CloudFormation Configuration
NewRelic::Observability::Workloads requires the following AWS Type Configuration per activated region

| Field           | Type   | Default                          | Required | Notes                                                                                                                       |
|-----------------|--------|----------------------------------|:--------:|-----------------------------------------------------------------------------------------------------------------------------|
| AccountID       | string | none                             |    Y     | [New Relic Account ID](https://docs.newrelic.com/docs/accounts/accounts-billing/account-structure/account-id/)              |
| APIKey          | string | none                             |    Y     | [New Relic User Key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#overview-keys)                      |
| Endpoint        | string | https://api.newrelic.com/graphql |    N     | [API endpoints](https://docs.newrelic.com/docs/apis/nerdgraph/get-started/introduction-new-relic-nerdgraph/#authentication) |

### AWS Console
[See bullet #6 here.](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/registry-public.html#registry-public-activate-extension-console)

### AWS CLI
```bash
aws cloudformation set-type-configuration --region <AWS_REGION> --type RESOURCE --type-name NewRelic::Observability::Workloads --configuration-alias default --configuration "{  \"APIKey\": \"<YOUR_NEWRELIC_API_KEY\",  
\"AccountID\": \"YOUR_NEWRELIC_ACCOUNT_ID\"}"
```
