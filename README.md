![Coverage](https://img.shields.io/badge/Coverage-47.9%25-yellow)

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

You can also use Docker to install b7s. See the [Docker documentation](/docker/README.md) for more information.

## Usage

| Flag        | Short Form | Default Value           | Description                                                                                   |
| ----------- | ---------- | ----------------------- | --------------------------------------------------------------------------------------------- |
| log-level   | -l         | "info"                  | Specifies the level of logging to use.                                                        |
| db          | -d         | "db"                    | Specifies the path to the database used for persisting node data.                             |
| role        | -r         | "worker"                | Specifies the role this node will have in the Blockless protocol (head or worker).            |
| address     | -a         | "0.0.0.0"               | Specifies the address that the libp2p host will use.                                          |
| port        | -p         | 0                       | Specifies the port that the libp2p host will use.                                             |
| private-key | N/A        | N/A                     | Specifies the private key that the libp2p host will use.                                      |
| concurrency | -c         | node.DefaultConcurrency | Specifies the maximum number of requests the node will process in parallel.                   |
| rest-api    | N/A        | N/A                     | Specifies the address where the head node REST API will listen on.                            |
| boot-nodes  | N/A        | N/A                     | Specifies a list of addresses that this node will connect to on startup, in multiaddr format. |
| workspace   | N/A        | "./workspace"           | Specifies the directory that the node can use for file storage.                               |
| runtime     | N/A        | N/A                     | Specifies the runtime address used by the worker node.                                        |

## Dependencies

b7s depends on the following repositories:

- blocklessnetwork/runtime
- blocklessnetwork/orchestration-chain

## Contributing

See src/README for information on contributing to the b7s project.
