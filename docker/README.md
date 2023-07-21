## Prerequisites

- A machine running Docker
- A valid AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY for an S3-compatible storage provider. This is used for backing up your node's keys and configuration.
- A valid KEY_PATH and KEY_PASSWORD for your S3-compatible storage provider. This is used for backing up your node's keys and configuration.

## Running the Image

First, you'll need to pull the latest version of the b7s Docker image from the public registry:

```bash
docker pull ghcr.io/blocklessnetwork/b7s:v0.0.25
```

Run the image

To Run the Node in the Head Node Configuration

```bash
docker run -d --name b7s \
  -e AWS_ACCESS_KEY_ID=<YOUR_AWS_ACCESS_KEY_ID> \
  -e AWS_SECRET_ACCESS_KEY=<YOUR_AWS_SECRET_ACCESS_KEY> \
  -e KEY_PATH=<YOUR_S3_KEY_PATH> \
  -e KEY_PASSWORD=<YOUR_S3_KEY_PASSWORD> \
  -e NODE_ROLE=head \
  -e P2P_PORT=9527 \
  -v /var/tmp/b7s/peerdb:/var/tmp/b7s/peerdb \
  -v /var/tmp/b7s/function-db:/var/tmp/b7s/function-db \
  -p 9527:9527 \
  ghcr.io/blocklessnetwork/b7s:v0.0.25
```
To Run the Node in the Worker Node Configuration

```bash
docker run -d --name b7s \
  -e AWS_ACCESS_KEY_ID=<YOUR_AWS_ACCESS_KEY_ID> \
  -e AWS_SECRET_ACCESS_KEY=<YOUR_AWS_SECRET_ACCESS_KEY> \
  -e KEY_PATH=<YOUR_S3_KEY_PATH> \
  -e KEY_PASSWORD=<YOUR_S3_KEY_PASSWORD> \
  -e NODE_ROLE=worker \
  -e P2P_PORT=9527 \
  -v /var/tmp/b7s/peerdb:/var/tmp/b7s/peerdb \
  -v /var/tmp/b7s/function-db:/var/tmp/b7s/function-db \
  -p 9527:9527 \
  ghcr.io/blocklessnetwork/b7s:v0.0.25
```
