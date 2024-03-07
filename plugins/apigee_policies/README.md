# Apigee Policies Plugin

This plugin uses custom OpenAPI extensions within to define and insert 
Apigee policy steps into the generated API proxy bundle.


## How to create an Apigee policy

To create an Apigee policy, use the `x-Apigee-Policies` extension at the top level.
Each policy must have a unique `.name` property

e.g.

```yaml
x-Apigee-Policies:
  - FlowCallout:
      .name: FC-Callout
      DisplayName: FC-Callout
      SharedFlowBundle: FC-Callout
```

## How to insert an Apigee policy

Apigee policies can be inserted at various places within the OpenAPI spec.

1. Within an existing operation in a path
2. At the top-level as a PreFlow policy
3. At the top-level as a PostFlow policy
4. At the top-level as a PostClientFlow policy

When inserting a policy, you use the policy `.name` to reference it.


## Inserting PreFlow Policies

To insert an Apigee PreFlow policy, use the `x-Apigee-PreFlow` extension.

e.g.

```yaml
x-Apigee-PreFlow:
  Request:
    - Step:
        Condition: true
        Policy: FC-Callout
  Response:
    - Step:
        Condition: true
        Policy: FC-Callout
```


## Inserting PostFlow Policies

To insert an Apigee PostFlow policy, use the `x-Apigee-PostFlow` extension.

e.g.

```yaml
x-Apigee-PostFlow:
  Request:
    - Step:
        Condition: true
        Policy: FC-Callout
  Response:
    - Step:
        Condition: true
        Policy: FC-Callout
```


## Inserting PostClientFlow Policies

To insert an Apigee PostClientFlow policy, use the `x-Apigee-PostClientFlow` extension.

e.g.

```yaml
x-Apigee-PostClientFlow:
  Response:
    - Step:
        Condition: true
        Policy: FC-Callout

```



## Supported Policies

### Flow Callout

Example Flow Callout policy as YAML.

```yaml
FlowCallout:
  .name: FC-Callout
  .enabled: true
  .continueOnError: true
  DisplayName: FC-Callout
  SharedFlowBundle: FC-Callout
  Parameters:
    - Parameter:
        .name: param1
        .value: Literal
    - Parameter:
        .name: param2
        .ref: request.content
```

which would get generated as XML like this

```text
<FlowCallout async="false" continueOnError="false" enabled="true" name="FC-MyCallout">
  <DisplayName>FC-Callout</DisplayName>
  <Parameters>
    <Parameter name="param1">Literal</Parameter>
    <Parameter name="param2">{request.content}</Parameter>
  </Parameters>
  <SharedFlowBundle>FC-MyCallout</SharedFlowBundle>
</FlowCallout>
```

