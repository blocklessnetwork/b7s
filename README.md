![Coverage](https://img.shields.io/badge/Coverage-48.1%25-yellow)

# b7s daemon

`blockless` the peer to peer networking daemon for the blockless network.

## installation

coming

## usage

commands
`b7s [command]`

- `help` display help menu
- `keygen` generate identity keys for the node

flags
`b7s --flag value`

- `config` path to the configuration file
- `out` style of logging used in the daemon (rich|text|json)

## depends

- [blocklessnetwork/runtime](https://github.com/blocklessnetwork/runtime)

## contributing

- golang 1.18
- clone repo
- `make`

### structure

The project structure is `mvc` style. The main entry uses `spf13/cobra`, and that setup can be found in `main.go`. The cobra `rootCMD` is set to start the `daemon`. While other sub commands, may run another service.

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
