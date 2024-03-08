# Apigee Policies Plugin

This plugin uses custom OpenAPI extensions within the spec to define and insert 
Apigee policy steps into the generated API proxy bundle.

It's also possible to move the extensions to separate files outside the main Open API Spec.
This is useful in cases where you need to share configuration across multiple API Proxies.


## How to create an Apigee policy

To create an Apigee policy, use the `x-Apigee-Policies` extension at the top level.
Each policy must have a unique `.name` property

e.g.

```yaml
x-Apigee-Policies:
  - FlowCallout:
      .name: FC-Callout
      DisplayName: FC-Callout
      SharedFlowBundle: SharedFlowName
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
        Name: FC-Callout
  Response:
    - Step:
        Condition: true
        Name: FC-Callout
```


## Inserting PostFlow Policies

To insert an Apigee PostFlow policy, use the `x-Apigee-PostFlow` extension.

e.g.

```yaml
x-Apigee-PostFlow:
  Request:
    - Step:
        Condition: true
        Name: FC-Callout
  Response:
    - Step:
        Condition: true
        Name: FC-Callout
```


## Inserting PostClientFlow Policies

To insert an Apigee PostClientFlow policy, use the `x-Apigee-PostClientFlow` extension.

e.g.

```yaml
x-Apigee-PostClientFlow:
  Response:
    - Step:
        Condition: true
        Name: FC-Callout

```

## Defining extensions in separate files

Each extension can be defined in-line within the OpenAPI spec, or within a separate file.

e.g.

```yaml
x-Apigee-Policies: 
  $ref: "./apigee-config.yaml#/Policies"

x-Apigee-PreFlow:
  $ref: "./apigee-config.yaml#/PreFlow"

x-Apigee-PostFlow:
  $ref: "./apigee-config.yaml#/PostFlow"
```



## Supported Policies

All Apigee policies are supported. 

## How to write Apigee policies as YAML

When creating an Apigee policy as YAML, you should be able to take an existing policy's
XML representation, and translate it to YAML by following a few basic rules.

  * XML elements are represented as YAML fields
  * XML elements attributes are represented as YAML fields prepended with a dot `.` 
  * XML elements content is represented as YAML fields prepended with `.@`
  * XML elements like this `<Simple>Value</Simple>` are represented like this `Simple: "Value"` 


See the examples below 

  * *Example 1*: XML element containing another XML element
    *  ```xml
       <Parent>
         <Child>foo</Child>
       </Parent>
       ```
       is equivalent to
       ```yaml
       Parent:
         Child: foo
       ```
  * *Example 2*: Simple XML element with no attributes, and scalar content
    * ```xml
      <Field>Content</Field>
      ```
      is equivalent to
      ```yaml
      Field: Content
      ```
  * *Example 3*: XML element with an attribute
    * ```xml
      <Parent foo="bar" />
      ```
      is equivalent to
      ```yaml
      Parent: 
        .foo: bar
      ```
  * *Example 4*: XML element with an attribute and  scalar content
    * ```xml
      <Parent foo="bar" >Content</Parent>
      ``` 
      is equivalent to
      ```yaml
      Parent:
        .foo: bar
        .@: Content
      ```
  * *Example 5*: XML sequence where parent has no attributes
    * ```xml
      <Parent>
        <Child>foo</Child>
        <Child>bar</Child>
      </Parent>
      ```
      is equivalent to
      ```yaml
      Parent:
        - Child: foo
        - Child: bar
      ```
  * *Example 6*: XML sequence where parent has attributes
    * ```xml
      <Parent attr1="value1" attr2="value2" >
        <Child>foo</Child>
        <Child>bar</Child>
      </Parent>
      ``` 
      is equivalent to
      ```yaml
      Parent:
        .attr1: value1
        .attr2: value2
        .@:
          - Child: foo
          - Child: bar
      ```
  * *Example 7*: XML sequence without parent 
    * ```xml
      <Root>
        <Child name="foo" />
        <Child name="bar" />
      </Root>
      ```
      is equivalent to
      ```yaml
      Root:
        .@:
          - Child:
            .name: foo
          - Child:
            .name: bar
      ```


## Apigee policies sample YAMLs

Below are several examples for common Apigee policies represented as YAML

### Flow Callout

Example Flow Callout policy as YAML.

```yaml
FlowCallout:
  .async: false
  .name: FC-Callout
  .enabled: true
  .continueOnError: true
  DisplayName: FC-Callout
  SharedFlowBundle: SharedFlowName
  Parameters:
    - Parameter:
        .name: param1
        .value: Literal
    - Parameter:
        .name: param2
        .ref: request.content
```

is equivalent to

```text
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<FlowCallout async="false" continueOnError="true" enabled="true" name="FC-Callout" >
  <DisplayName>FC-Callout</DisplayName>
  <Parameters>
    <Parameter name="param1" value="Literal" ></Parameter>
    <Parameter name="param2" ref="request.content" ></Parameter>
  </Parameters>
  <SharedFlowBundle>SharedFlowName</SharedFlowBundle>
</FlowCallout>
```


### Raise Fault

Example Raise Fault policy represented as YAML
```yaml
RaiseFault:
  .async: false
  .name: RF-Example
  .enabled: true
  .continueOnError: true
  DisplayName: RF-Example
  IgnoreUnresolvedVariables: true
  ShortFaultReason: false
  FaultResponse:
    - AssignVariable:
        Name: flow.var
        Value: 123
    - Add:
        Headers:
          - Header:
              .name: user-agent
              .@: example
    - Copy:
        .source: request
        Headers:
          - Header:
              .name: header-name
        StatusCode: 304
    - Remove:
        Headers:
          - Header:
              .name: sample-header
    - Set:
        Headers:
          - Header:
              .name: user-agent
              .@: "{request.header.user-agent}"
        Payload:
          .contentType: application/json
          .@: '{"name":"foo", "type":"bar"}'
    - Set:
        ReasonPhrase: Server Error
        StatusCode: 500

```

is equivalent to

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<RaiseFault async="false" continueOnError="true" enabled="true" name="RF-Example" >
    <DisplayName>RF-Example</DisplayName>
    <FaultResponse>
        <AssignVariable >
            <Name>flow.var</Name>
            <Value>123</Value>
        </AssignVariable>
        <Add >
            <Headers>
                <Header name="user-agent" >example</Header>
            </Headers>
        </Add>
        <Copy source="request" >
            <Headers>
                <Header name="header-name" ></Header>
            </Headers>
            <StatusCode>304</StatusCode>
        </Copy>
        <Remove >
            <Headers>
                <Header name="sample-header" ></Header>
            </Headers>
        </Remove>
        <Set >
            <Headers>
                <Header name="user-agent" >{request.header.user-agent}</Header>
            </Headers>
            <Payload contentType="application/json" >{"name":"foo", "type":"bar"}</Payload>
        </Set>
        <Set >
            <ReasonPhrase>Server Error</ReasonPhrase>
            <StatusCode>500</StatusCode>
        </Set>
    </FaultResponse>
    <IgnoreUnresolvedVariables>true</IgnoreUnresolvedVariables>
    <ShortFaultReason>false</ShortFaultReason>
</RaiseFault>

```



