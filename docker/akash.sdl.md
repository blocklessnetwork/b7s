# Deploying a Compute Node On Akash

This is a pretty easy process documented here!

Here is the `Akash` SDL used to launch a `worker` node on `Blockless`!

```yaml
# this is an example of launching a node on akash
---
version: "2.0"

services:
  web:
    image: ghcr.io/blocklessnetwork/b7s:v0.0.5-rc1
    env:
      - NODE_ROLE=worker
      - AWS_ACCESS_KEY_ID=s3_id
      - AWS_SECRET_ACCESS_KEY=s3_key
      - KEY_PATH=my_backup_key_path
      - KEY_PASSWORD=my_key_password
      - BOOT_NODES=/ip4/206.81.10.243/tcp/31368/p2p/12D3KooWKTKwW1y6iRoGeag1y2wpNYV9UM8QaYXaTL7DUTZdEFw6,/ip4/147.75.84.103/tcp/31282/p2p/12D3KooWFJxHN8NMAQkU6HdqvMM4wpSKnsyzSb2rJwMhg2w4UsWh,/ip4/147.75.199.31/tcp/31969/p2p/12D3KooWH7TiTXb4UzKKU5zMkCDSMx5KAwP4waBbknxnoXfWL7Db
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
