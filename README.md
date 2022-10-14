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

```bash
b7s --config=../configs/head-config.yaml --out=json
```

## depends

- [blocklessnetwork/runtime](https://github.com/blocklessnetwork/runtime)

## contributing

see [src/readme](src/README.md)
