[![Coverage](https://img.shields.io/badge/Coverage-64.5%25-yellow)](https://img.shields.io/badge/Coverage-64.5%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/blocklessnetwork/b7s)](https://goreportcard.com/report/github.com/blocklessnetwork/b7s) 
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/blocklessnetwork/b7s/blob/main/LICENSE.md) 
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/blocklessnetwork/b7s)](https://img.shields.io/github/v/release/blocklessnetwork/b7s)


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

For a more detailed overview of the configuration options, see the [b7s-node Readme](/cmd/node/README.md#usage).

| Flag                      | Short Form | Default Value           | Description                                                                                           |
| ------------------------- | ---------- | ----------------------- | ----------------------------------------------------------------------------------------------------- |
| config                    | N/A        | N/A                     | Specifies the config file to load.                                                                    |
| log-level                 | -l         | "info"                  | Specifies the level of logging to use.                                                                |
| db                        | N/A        | "db"                    | Specifies the path to database used for persisting peer and function data.                            |
| role                      | -r         | "worker"                | Specifies the role this node will have in the Blockless protocol (head or worker).                    |
| address                   | -a         | "0.0.0.0"               | Specifies the address that the libp2p host will use.                                                  |
| port                      | -p         | 0                       | Specifies the port that the libp2p host will use.                                                     |
| websocket-port            | N/A        | 0                       | Specifies the port that the libp2p host will use for websocket connections.                           |
| private-key               | N/A        | N/A                     | Specifies the private key that the libp2p host will use.                                              |
| concurrency               | -c         | node.DefaultConcurrency | Specifies the maximum number of requests the node will process in parallel.                           |
| rest-api                  | N/A        | N/A                     | Specifies the address where the head node REST API will listen on.                                    |
| boot-nodes                | N/A        | N/A                     | Specifies a list of addresses that this node will connect to on startup, in multiaddr format.         |
| workspace                 | N/A        | "./workspace"           | Specifies the directory that the node can use for file storage.                                       |
| dialback-address          | N/A        | N/A                     | Specifies the advertised dialback address of the Node.                                                |
| dialback-port             | N/A        | N/A                     | Specifies the advertised dialback port of the Node.                                                   |
| websocket-dialback-port   | N/A        | 0                       | Specifies the advertised dialback port for Websocket connections.                                     |
| cpu-percentage-limit      | N/A        | 1.0                     | Specifies the amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited. |
| memory-limit              | N/A        | N/A                     | Specifies the memory limit for Blockless Functions, in kB.                                            |
| no-dialback-peers         | N/A        | false                   | Specifies if the node should avoid dialing back peers known from past runs                            |
| load-attributes           | N/A        | false                   | Specifies if the node should load attributes from the environment.                                    |
| topics                    | N/A        | N/A                     | Specifies topics that the node sould subscribe to.                                                    |
| websocket                 | -w         | false                   | Specifies if the node should use websocket protocol for communication, besides TCP.                   |
| must-reach-boot-nodes     | N/A        | false                   | Specifies if the node should fail if it cannot reach boot nodes on startup.                           |
| disable-connection-limits | N/A        | false                   | Specifies if the node should try to maintain as many connections as possible.                         |
| connection-count          | N/A        | N/A                     | Specifies the number of connections that the node will aim to have.                                   |
| runtime-path              | N/A        | N/A                     | Specifies the local path to the Blockless Runtime.                                                    |
| runtime-cli               | N/A        | N/A                     | Specifies the name of the Blockless Runtime executable, as found in the runtime-path.                 |
| enable-tracing            | N/A        | false                   | Specifies whether the node should emit tracing data.                                                  |
| tracing-grpc-endpoint     | N/A        | N/A                     | Specifies the GRPC endpoint where node should send tracing data.                                      |
| tracing-http-endpoint     | N/A        | N/A                     | Specifies the HTTP endpoint where node should send tracing data.                                      |
| enable-metrics            | N/A        | false                   | Specifies whether the node should serve metrics.                                                      |
| prometheus-address        | N/A        | N/A                     | Specifies the address where node should serve metrics (for head node this is the REST API address)    |





## Dependencies

b7s depends on the following repositories:

- blocklessnetwork/runtime
- blocklessnetwork/orchestration-chain

## Contributing

See src/README for information on contributing to the b7s project.
