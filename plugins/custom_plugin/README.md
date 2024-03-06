## Custom Plugin Example

This plugin detects if any operation in the OpenAPI spec has an extension named "x-visibility" with
the following structure

```json
{
  "x-visibility": {
    "extent": "INTERNAL"
  }
}
```

If operation is tagged as internal, it's removed from the spec model before it's passed down to the rest of the generator pipeline.