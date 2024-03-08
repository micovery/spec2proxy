## spec2proxy 
This is a command-line tool that generates an Apigee X API Proxy bundle from an OpenAPI 3 Spec.

By default, the tool generates a simple API Proxy bundle that can serve as scaffolding for
building more complex proxies.

It's possible to customize the generation logic through the use of plugins.
For example, you can create a plugin that understands OpenAPI spec extensions, and
customizes the generated API Proxy bundle based on the value of the extensions.



### How to use it
```shell
spec2proxy -oas petstore.yaml -out ./petstore
```

### How to use it with plugins

You can pass one or more plugins to use with the `-plugins` parameter.

e.g.
```shell
spec2proxy -oas petstore.yaml -out ./petstore -plugins example,custom_plugin,etc
```


### How the tool works

This tool works as basic processing pipeline with three steps: *Parse*, *Transform*, and *Generate*

e.g.
```text
libopenapi_model = Parse(openapi_text) 
apigee_model = Transform(libopenapi_model)
apigee_bundle = Generate(apigee_model)

```

### How plugins work

Plugins are hooks into the processing pipeline of the tool. Each plugin has two hooks *ProcessSpecModel* and *ProcessProxyModel*

* *ProcessSpecModel* - This function runs after the input spec text has been *parsed* into the [libopenapi](https://github.com/pb33f/libopenapi) data model.
* *ProcessProxyModel* - This function runs after the spec has been *transformed* into the Apigee data model.

### How to add plugins

To get started, you make a copy the [example](/plugins/example) plugin, and register it in the [init.go](/plugins/init.go) file.

```go
plugins.RegisterPlugin(&custom_plugin.Plugin{})
```
Then, you need to re-compile the generator.

```shell
go build -o spec2proxy cmd/spec2proxy/main.go 
```


### Available plugins

 The following plugins are available to be used out of the box
 * [apigee_policies](/plugins/apigee_policies) - Supports adding and using Apigee policies
 * [custom_plugin](/plugins/custom_plugin) - Shows how to traverse and manipulate the data models
 * [example](/plugins/example) - Serves as template for creating new plugins

### What about Go-Lang Plugin package ...

Plugins are not dynamic libraries like those built with Go-Lang's plugin package. 
Instead, plugins are compiled into the generator itself. This is on purpose for the sake of portability,
and ease of development of the plugins.


### Support
This is not an officially supported Google product
