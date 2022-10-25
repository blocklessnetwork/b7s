## Rest Api Discussion

There should be a small discussion about the shape of the responses when using the rest API. The Responses will have multiple status.

```json
{
  "status": "200", // standard HTTP API return, if the API is functioning normally, should always return 200
  "data": {
    "code": "200" // inner API response, from querying workers, or trying to do user actions
  }
}
```
