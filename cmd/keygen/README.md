# KeyGen

## Description

The `keygen` utility can be used to create keys that will determine Blockless b7s Node identity.

## Usage

```console
Usage of keygen:
  -o, --output string   directory where keys should be stored (default ".")
```

## Examples

Create keys in the `keys` directory of the users `home` directory:

```console
$ ./keygen --output ~/keys
generated private key: /home/user/keys/priv.bin
generated public key: /home/user/keys/pub.bin
generated identity file: /home/user/keys/identity
```
