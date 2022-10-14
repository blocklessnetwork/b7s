## daemon

This part of the application is responsible for 'glueing' a running service together. It runs all other parts of the applications as `go functions` with the main `func` yielding to a `select {}`
