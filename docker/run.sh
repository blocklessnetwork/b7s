#!/bin/bash

if [ -n "$KEY_PATH" ]; then
  s3_uri_base="s3://${KEY_PATH}"
  aws_args="--endpoint-url ${S3_HOST}"
  if [ -n "$KEY_PASSWORD" ]; then
    file_suffix=".gpg"
  else
    file_suffix=""
  fi
fi

restore_key () {
  existing=$(aws $aws_args s3 ls "${s3_uri_base}/$1" | head -n 1)
  if [[ -z $existing ]]; then
    echo "$1 backup not found"
  else
    echo "Restoring $1"
    aws $aws_args s3 cp "${s3_uri_base}/$1" $CONFIG_PATH/$1$file_suffix
    if [ -n "$KEY_PASSWORD" ]; then
      echo "Decrypting"
      gpg --decrypt --batch --passphrase "$KEY_PASSWORD" $CONFIG_PATH/$1$file_suffix > $CONFIG_PATH/$1
      rm $CONFIG_PATH/$1$file_suffix
    fi
  fi
}

backup_key () {
  existing=$(aws $aws_args s3 ls "${s3_uri_base}/$1" | head -n 1)
  if [[ -z $existing ]]; then
    echo "Backing up $1"
    if [ -n "$KEY_PASSWORD" ]; then
      echo "Encrypting backup..."
      gpg --symmetric --batch --passphrase "$KEY_PASSWORD" $CONFIG_PATH/$1
    fi
    aws $aws_args s3 cp $CONFIG_PATH/$1$file_suffix "${s3_uri_base}/$1"
    [ -n "$KEY_PASSWORD" ] && rm $CONFIG_PATH/$1.gpg
  fi
}

# Restore keys
if [ -n "$KEY_PATH" ]; then
  for f in $(aws $aws_args s3 ls "${s3_uri_base}/" | awk '{print $4}'); do
    restore_key "$f"
  done

  if [ -n "$KEY_PASSWORD" ]; then
    if [ -f "/app/keys/wallet.key" ]; then
      echo $KEY_PASSWORD | blsd keys import node /app/keys/wallet.key --keyring-backend=test  --home=/app/.blockless-chain
    fi
  fi
fi

if [ -f "$CONFIG_PATH/identity" ]; then
  echo "Node restored from backup"
else 
  echo "Generating New Node Identity"
  blsd keys add node --keyring-backend=test --home=/app/.blockless-chain > /dev/null 2>&1
  mkdir keys && cd keys
  if [ -n "$KEY_PASSWORD" ]; then
    echo $KEY_PASSWORD | blsd keys export node --keyring-backend=test --home=/app/.blockless-chain > /app/keys/wallet.key
  fi
  ../b7s keygen
  # Backup keys
  if [ -n "$KEY_PATH" ]; then
    # backup the on chain node identity
    for f in $(ls /app/keys/*); do
      backup_key $(basename "$f")
    done
  fi
fi

# run  template against the config file
/app/gomplate -f /app/docker-config.yaml -o /app/docker-config-env.yaml

# run the node
cd /app
./b7s -c docker-config-env.yaml