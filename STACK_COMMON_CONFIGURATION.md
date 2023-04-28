# Stack common configuration

This document documents the configuration common to all Stacks using NewRelic::Observability public extensions.

## Model
| Field           | Type   | Default | Create | Update | Delete | Read | List | Notes |
|-----------------|--------|---------|:------:|:------:|:------:|:----:|:----:|-------|
| Guid            | string | none    |        |   R    |   R    |  R   |      |       |
| ListQueryFilter | string | none    |        |        |        |      |  R   |       |
| Variables       | Object | none    |   O    |   O    |        |  O   |  O   |       |
| Tags            | Object | none    |   O    |   O    |        |      |      |       |                                                                                                                             |


### Guid
`Guid` New Relic entity identifier. Typically the `guid` or `id` value from Nerd Graph.

### ListQueryFilter
Entity search query string. The query string can search for an exact or fuzzy match on name, as well as searching several other attributes.

Note: you must supply either a query OR a queryBuilder argument, not both.

Operators available: =, AND, IN, LIKE

Special characters (.,;:*-_) are treated as whitespace. For example, name LIKE ':aws:' will match -aws. or foo aws.

Tags can be referenced in multiple ways with or without backticks.

Examples:
```
"name = 'MyApp (Staging)'
"name LIKE 'MyApp' AND type IN ('APPLICATION')"
"reporting = 'false' AND type IN ('HOST')"
"domain IN ('INFRA', 'APM')"
tags.Environment = 'staging' AND type IN ('APPLICATION')
```

[See also](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-entities-api-tutorial/#search-entity)

### Variables
`Variables` is a list of key/value pairs (string/string) (yaml object) that are substituted in any other Model variable using [Moustache](#Moustache) allowing for parameterized input at the CloudFormation level. For instance this Model 
fragment:
```yaml
AWSTemplateFormatVersion: 2010-09-09
Description: Sample New Relic Dashboards Template
Resources:
  Resource1:
    Type: 'NewRelic::Observability::Dashboards'
    Properties:
      Dashboard: >-
        dashboard: {description: "CloudFormation test dashboard", name: "CloudFormation Test Dashboard", pages: 
            {description: "TD PAGE", name: "TD Page", widgets: 
              {title: "Widget Title", configuration: 
                {markdown: {text: "Some markdown"}}}}, permissions: PRIVATE}
      Tags:
        SomeTag: "{{{Environment}}}"
      Variables:
        Environment: "Production"

Outputs:
  CustomResourceGuid:
    Value: !GetAtt Resource1.Guid
```

Variable names reserved by the system:
- ACCOUNTID
- GUID
- FRAGMENT

If you use a reserved Variable your value will be overwritten by the system.

### Tags
`Tags` is a list of key/value pairs (string/string) (yaml object) that are attached to the Created/Updated entity with the [NerdGraph Tagging API](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-tagging-api-tutorial/)

A `Create` operation uses the `taggingAddTagsToEntity` mutation, an `Update` operatino uses the `taggingReplaceTagsOnEntity` mutation.

***IMPORTANT NOTE***: An `Update` operation _completely_ overwrites all previous tags on the entity.

## Moustache
All text substitution is done using a [Go implementation](https://github.com/cbroglie/mustache) of the [Moustache specification](https://github.com/mustache/spec)`. [The manual is here](http://mustache.github.io/mustache.5.html), in
general all you need to know is use triple curly braces.
