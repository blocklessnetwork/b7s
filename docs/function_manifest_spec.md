```json
{
  "function": {
    "id": "org.blockless.functions.myfunction",
    "name": "The Best Little Function In TX (Labs)",
    "version": "1.0.0",
    "build-command": "npm build",
    "build-output": "./output.wasm",
    "runtime": "<1.0.0",
    "extensions": ["foundation/ipfs/>=1.0.0"]
  },
  "deployment": {
    "uri": "http://localhost:8080/someid/web-framework.tar.gz",
    "checksum": "0cbaf5c9d0aa075d546a9084096ce380",
    "permission": "private",
    "methods": [
      {
        "name": "web-framework",
        "entry": "web-framework.wasm",
        "arguments": [{name:, value:],
        "envvars":[{name:, value:}]
      }
    ],
    "aggregation": "foundation/average/>=1.0.0",
    "nodes": 5
  }
}
```
