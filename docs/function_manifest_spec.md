# function manifest spec

this is the `json` specification for the manifests that describe how a function should work.

```json
{
  "function": {
    // organization and function identifier
    "id": "org.blockless.functions.myfunction",

    // descriptive name of the function
    "name": "The Best Little Function In TX (Labs)",

    // version of the function
    "version": "1.0.0",

    // whats the runtime specification
    "runtime": "<1.0.0",

    // what extensions are needed to run this function
    "extensions": ["foundation/ipfs/>=1.0.0"]
  },
  "deployment": {
    // a manifest can provide a CID meaning its
    // been stored on filecoin
    // we will try to use multiple known filecoin portals
    // to retrieve the package
    "cid": "basy1234566",

    // because we want blockless to remain decoupled
    // a uri can also be specified to retrieve the package
    "uri": "basy1234566",

    // this is the checksum of the gzipped package
    "checksum": "0cbaf5c9d0aa075d546a9084096ce380",

    // describes the wasm builds located inside
    "methods": [
      {
        "name": "web-framework",
        "entry": "web-framework.wasm",
        "arguments":[{"name": "name", "value": "value"}]
        "envvars":[{"name": "name", "value": "value"}]
      }
    ],

    // what aggregation
    "aggregation": "foundation/average/>=1.0.0",

    // how many executions
    "nodes": 5
  }
}
```
