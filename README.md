![Coverage](https://img.shields.io/badge/Coverage-73.2%25-brightgreen)

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

| Flag             | Short Form | Default Value           | Description                                                                                   |
| ---------------- | ---------- | ----------------------- | --------------------------------------------------------------------------------------------- |
| log-level        | -l         | "info"                  | Specifies the level of logging to use.                                                        |
| peer-db          | N/A        | "peer-db"               | Specifies the path to database used for persisting peer data.                                 |
| function-db      | N/A        | "function-db"           | Specifies the path to database used for persisting function data.                             |
| role             | -r         | "worker"                | Specifies the role this node will have in the Blockless protocol (head or worker).            |
| address          | -a         | "0.0.0.0"               | Specifies the address that the libp2p host will use.                                          |
| port             | -p         | 0                       | Specifies the port that the libp2p host will use.                                             |
| private-key      | N/A        | N/A                     | Specifies the private key that the libp2p host will use.                                      |
| concurrency      | -c         | node.DefaultConcurrency | Specifies the maximum number of requests the node will process in parallel.                   |
| quorum           | -1         | node.DefaultQuorum      | Specifies the number of execution responses to require on each execution.                     |
| rest-api         | N/A        | N/A                     | Specifies the address where the head node REST API will listen on.                            |
| boot-nodes       | N/A        | N/A                     | Specifies a list of addresses that this node will connect to on startup, in multiaddr format. |
| workspace        | N/A        | "./workspace"           | Specifies the directory that the node can use for file storage.                               |
| runtime          | N/A        | N/A                     | Specifies the runtime address used by the worker node.                                        |
| dialback-address | N/A        | N/A                     | Specifies the advertised dialback address of the Node.                                        |
| dialback-port    | N/A        | N/A                     | Specifies the advertised dialback port of the Node.                                           |
| cpu-percentage-limit | N/A | 1.0 | Specifies the amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited. |
| memory-limit | N/A | N/A | Specifies the memory limit for Blockless Functions, in kB. |

## Dependencies

b7s depends on the following repositories:

- blocklessnetwork/runtime
- blocklessnetwork/orchestration-chain

## Contributing

See src/README for information on contributing to the b7s project.
