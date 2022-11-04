![Coverage](https://img.shields.io/badge/Coverage-48.1%25-yellow)

# b7s daemon

`blockless` the peer to peer networking daemon for the blockless network.

Supported Platforms

| OS      | arm64 | x64 |
| ------- | ----- | --- |
| Windows |       | x   |
| Linux   | x     | x   |
| MacOS   | x     | x   |

Using **curl**:

```bash
sudo sh -c "curl https://raw.githubusercontent.com/blocklessnetwork/b7s/main/download.sh | bash"
```

Using **wget**:

```bash
sudo sh -c "wget https://raw.githubusercontent.com/blocklessnetwork/b7s/main/download.sh -v -O download.sh; chmod +x download.sh; ./download.sh; rm -rf download.sh"
```

Use `docker` see docker [docs](docker/README.md)

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
