
# Blockless Node 

## Description

Blockless b7s Node is a peer-to-peer networking daemon for the blockless network.

A Node in the Blockless network can have one of two roles - it can be a Head Node or a Worker Node.

In short, Worker Nodes are nodes that will be doing the actual execution of work within the Blockless P2P network.
Worker Nodes do this by relying on the Blockless Runtime.
Blockless Runtime needs to be available locally on the machine where the Node is run.

Head Nodes are nodes that are performing coordination of work between a number of Worker Nodes.
When a Head Node receives an execution request to execute a piece of work (a Blockless Function), it will start a process of finding a Worker Node most suited to do this work.
Head Node does not need to have access to Blockless Runtime.

Head Nodes also serve a REST API that can be used to query or trigger certain actions.

## Usage

There are two ways of specifying configuration options - the CLI flags and a config file.
CLI flags override any configuration options in the config file.

List of supported CLI flags is listed below.

```console
Usage of b7s-node:
      --config string                  path to a config file
  -r, --role string                    role this note will have in the Blockless protocol (head or worker) (default "worker")
  -c, --concurrency uint               maximum number of requests node will process in parallel (default 10)
      --boot-nodes strings             list of addresses that this node will connect to on startup, in multiaddr format
      --workspace string               directory that the node can use for file storage
      --attributes                     node should try to load its attribute data from IPFS
      --peer-db string                 path to the database used for persisting peer data (default "peer-db")
      --function-db string             path to the database used for persisting function data (default "function-db")
      --topics strings                 topics node should subscribe to
  -l, --log-level string               log level to use (default "info")
  -a, --address string                 address that the b7s host will use (default "0.0.0.0")
  -p, --port uint                      port that the b7s host will use
      --private-key string             private key that the b7s host will use
  -w, --websocket                      should the node use websocket protocol for communication
      --websocket-port uint            port to use for websocket connections
      --dialback-address string        external address that the b7s host will advertise (default "0.0.0.0")
      --dialback-port uint             external port that the b7s host will advertise
      --websocket-dialback-port uint   external port that the b7s host will advertise for websocket connections
      --runtime-path string            Blockless Runtime location (used by the worker node)
      --runtime-cli string             runtime CLI name (used by the worker node) (default "bls-runtime")
      --cpu-percentage-limit float     amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited (default 1)
      --memory-limit int               memory limit (kB) for Blockless Functions
      --rest-api string                address where the head node REST API will listen on
```

Alternatively to the CLI flags, you can create a YAML file and specify the parameters there.
All of the CLI flags have a corresponding config option in the YAML file.
In the config file, parameters are grouped based on the functionality they impact.

Example configuration for a worker node:

```yaml
role: worker
concurrency: 10
workspace: /tmp/workspace
attributes: false

log:
  level: debug

connectivity:
  address: 127.0.0.1
  port: 9000
  private-key: /home/user/.b7s/path/to/priv/key.bin
  websocket: true


worker:
  runtime-path: /home/user/.local/blockless-runtime/bin
  cpu-percentage-limit: 0.8

```

You can find a more complete reference in [example.yaml](/cmd/node/example.yaml).

### Configuration Option Details

You can find more information about `multiaddr` format for network addresses [here](https://github.com/multiformats/multiaddr) and [here](https://multiformats.io/multiaddr/).

Private key path relates to the private key created by the [keygen](/cmd/keygen/README.md) utility.
Using the same private key in multiple `node` runs will ensure the node has the same identity on the network.
If a private key is not specified the node will start with a randomly generated identity.

## Examples

### Starting a Worker Node

```console
$ ./node --peer-db peer-database --log-level debug --port 9000 --role worker --runtime ~/.local/bin --workspace workspace --private-key ./keys/priv.bin
```

The created `node` will listen on all addresses on TCP port 9000.
Database used to persist Node data between runs will be created in the `peer-database` subdirectory.
On the other hand, Node will persist function data in the default database, in the `function-db` subdirectory.

Blockless Runtime path is given as `/home/user/.local/bin`.
At startup, node will check if the Blockless Runtime is actually found there, namely the [bls-runtime](https://blockless.network/docs/protocol/runtime).

Node Identity will be determined by the private key found in `priv.bin` file in the `keys` subdirectory.

Any transient files needed for node operation will be created in the `workspace` subdirectory.

### Starting a Head Node

```console
$ ./node --peer-db /var/tmp/b7s/peerdb --function-db /var/tmp/b7s/fdb --log-level debug --port 9002 -r head --workspace /var/tmp/b7s/workspace --private-key ~/keys/priv.bin --rest-api ':8080'
```

The created `node` will listen on all addresses on TCP port 9002.
Database used to persist Node peer data between runs will be created at `/var/tmp/b7s/peerdb`.
Database used to persist Node function data will be created at `/var/tmp/b7s/fdb`.

Any transient files needed for node operation will be created in the `/var/tmp/b7s/workspace` directory.

Head Node REST API will be available on all addresses on port 8080.
