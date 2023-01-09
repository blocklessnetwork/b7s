![Coverage](https://img.shields.io/badge/Coverage-48.1%25-yellow)

# b7s daemon

b7s is a peer-to-peer networking daemon for the blockless network. It is supported on Windows, Linux, and MacOS platforms for both x64 and arm64 architectures.

## Installation

You can install b7s using either curl or wget:

```bash
# using curl
sudo sh -c "curl https://raw.githubusercontent.com/blocklessnetwork/b7s/main/download.sh | bash"

# using wget
sudo sh -c "wget https://raw.githubusercontent.com/blocklessnetwork/b7s/main/download.sh -v -O download.sh; chmod +x download.sh; ./download.sh; rm -rf download.sh"
```

You can also use Docker to install b7s. See the [Docker documentation](https://chat.openai.com/chat/docker/README.md) for more information.

Usage
b7s can be run with a number of commands and flags:

Commands:

- `help`: display the help menu
- `keygen`: generate identity keys for the node
  Flags:

- `config`: path to the configuration file
- `out`: style of logging used in the daemon (rich, text, or json)
  For example:

## Dependencies

b7s depends on the following repositories:

- blocklessnetwork/runtime
- blocklessnetwork/orchestration-chain

## Contributing

See src/README for information on contributing to the b7s project.
