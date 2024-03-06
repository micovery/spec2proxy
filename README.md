## spec2proxy 
This is a command-line tool that generates an Apigee X API Proxy from an OpenAPI Spec.

By default, the tool generates a very simple API Proxy that can serve as scaffolding for
building more complex proxies.

However, it's possible to customize the generation logic through the use of Go-Lang plugins.
This allows you to define your own OpenAPI extensions, and use them to affect how the
Apigee API Proxies are generated.



### How to use it
```shell
spec2proxy --spec petstore.yaml --out ./petstore
```

### How to use it with plugins

```shell
spec2proxy --spec petstore.yaml --out ./petstore --plugins myplugin
```

