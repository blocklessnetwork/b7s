
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

```console
Usage of node:
  -a, --address string       address that the libp2p host will use (default "0.0.0.0")
      --boot-nodes strings   list of addresses that this node will connect to on startup, in multiaddr format
  -c, --concurrency uint     maximum number of requests node will process in parallel (default 10)
  -d, --db string            path to the database used for persisting node data (default "db")
  -l, --log-level string     log level to use (default "info")
  -p, --port uint            port that the libp2p host will use
      --private-key string   private key that the libp2p host will use
      --rest-api string      address where the head node REST API will listen on
  -r, --role string          role this note will have in the Blockless protocol (head or worker) (default "worker")
      --runtime string       runtime address (used by the worker node)
      --workspace string     directory that the node can use for file storage (default "./workspace")
```

You can find more information about `multiaddr` format for network addresses [here](https://github.com/multiformats/multiaddr) and [here](https://multiformats.io/multiaddr/).

Private key path relates to the private key created by the [keygen](/cmd/keygen/README.md) utility.
Using the same private key in multiple `node` runs will ensure the node has the same identity on the network.
If a private key is not specified the node will start with a randomly generated identity.

## Examples

### Starting a Worker Node

```console
$ ./node --db ./database --log-level debug --port 9000 --role worker --runtime ~/.local/bin --workspace workspace --private-key ./keys/priv.bin
```

The created `node` will listen on all addresses on TCP port 9000.
Database used to persist Node data between runs will be created in the `database` subdirectory.

Blockless Runtime path is given as `~/.local/bin`.
At startup, node will check if the Blockless Runtime is actually found there, namely the [blockless-cli](https://blockless.network/docs/cli).

Node Identity will be determined by the private key found in `priv.bin` file in the `keys` subdirectory.

Any transient files needed for node operation will be created in the `workspace` subdirectory.

### Starting a Head Node

```console
$ ./node --db /var/tmp/b7s/db --log-level debug --port 9002 -r head --workspace /var/tmp/b7s/workspace --private-key ~/keys/priv.bin --rest-api ':8080'
```

The created `node` will listen on all addresses on TCP port 9002.
Database used to persist Node data between runs will be created at `/var/tmp/b7s/db`.

Any transient files needed for node operation will be created in the `/var/tmp/b7s/workspace` directory.

Head Node REST API will be available on all addresses on port 8080.
