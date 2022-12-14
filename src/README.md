# contributing

- golang 1.18
- clone repo
- `make`

## structure

The main entry uses `spf13/cobra`, and that setup can be found in `main.go`. The cobra `rootCMD` is set to start the `daemon`. While other sub commands, may run another service.

```
src/
├─ main.go
├─ controller/
├─ enums/
├─ models/
├─ messaging/
Makefile
```

`controller` package will contain all the methods that are used inside the daemon to do something. `models` contain structs, but also contains helpers to extend those models, importantly all `message` structs are defined here to pass through the `p2p` channels. `enums` makes our communication consistent. `messaging` contains the wiring to send messages on the network, but also has handlers defined to react to messages sent on the network.

## messaging paths
