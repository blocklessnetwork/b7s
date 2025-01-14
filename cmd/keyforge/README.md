## KeyForge

The `keyforge` utility is used for managing cryptographic keys and signing/verifying messages or data. It provides the following options:

### Generate Keys

Generate a new keypair and save it to a file:

$ ./keyforge -o

### Sign and Verify

Sign and verify messages or data using your generated keys:

#### Sign a Message

Sign a message and save the signature:

$ ./keyforge -s "Your message" -o

#### Sign a File

Sign a file and save the signature:

$ ./keyforge -f -o

#### Verify a Signature

Verify a message or file's signature using the \`keyforge\` utility:

$ ./keyforge -pubkey -message "Original message" -signature

#### Verify a Signature with PeerID

Verify a message or file's signature using a PeerID with the \`keyforge\` utility:

$ ./keyforge -peerid -message "Original message" -signature

#### Verify a Signature with OpenSSL

Verify a message or file's signature using OpenSSL:

##### Create a Signature

Use OpenSSL to create a signature:

$ openssl dgst -sha256 -sign -out message.sig

##### Verify a Signature

Use OpenSSL to verify a signature:

$ openssl dgst -sha256 -verify -signature message.sig -in

These commands enable you to manage cryptographic keys and perform signing and verification operations, including using OpenSSL for verification, conveniently within the Bless b7s Node network.
