# function execution request spec

```json
{
  // id of the function to execute
  "function_id": "com.some.foo.function",

  // method to execute in the function archive
  "method": "web-framework.wasm",

  // execution parameters
  "parameters": null,

  // execution options
  "config": {
    "stdin" : "",
    // environment variables to load before execution
    "env_vars": [
      {
        "name": "BLS_REQUEST_PATH",
        "value": "/"
      }
    ],

    // number of nodes to execute on
    "number_of_nodes": 1,

    // result aggregation method to use
    "result_aggregation": {
      "enable": false,
      "type": "none",
      "parameters": [
        {
          "name": "type",
          "value": ""
        }
      ]
    }
  }
}
```
