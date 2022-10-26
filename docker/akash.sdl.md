# Deploying a Compute Node On Akash

This is a pretty easy process documented here!

```yaml
# this is an example of launching a node on akash
---
version: "2.0"

services:
  web:
    image: ghcr.io/blocklessnetwork/b7s:v0.0.2-rc1
    env:
      - NODE_ROLE=head
      - AWS_ACCESS_KEY_ID=s3_id
      - AWS_SECRET_ACCESS_KEY=s3_key
      - KEY_PATH=my_backup_key_path
      - KEY_PASSWORD=my_key_password
      - BOOT_NODES="/ip4/206.81.10.243/tcp/31368/p2p/12D3KooWKTKwW1y6iRoGeag1y2wpNYV9UM8QaYXaTL7DUTZdEFw6"
    expose:
      - port: 9527
        to:
          - global: true
profiles:
  compute:
    web:
      resources:
        cpu:
          units: 0.5
        memory:
          size: 1024Mi
        storage:
          size: 512Mi

  placement:
    dcloud:
      pricing:
        web:
          denom: uakt
          amount: 1000

deployment:
  web:
    dcloud:
      profile: web
      count: 1
```
