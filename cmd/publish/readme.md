# algo-collect publisher

A simple demo project to act as a base for publishing collections using the algo-collect contract.

This project is **a quick demo** and needs some clean up after the initial round of feedback.

This project **hasn't been security audited** and should only be used with caution and understanding.

## Requirements

- Linux or macOS
- Golang version 1.17.0 or higher
- Python 3. The scripts assumes the Python executable is called `python3`.
- The [Algorand Node software][algorand-install] installed using the updater script; with the enviorment variable `ALGORAND_DATA` set and the toolset added to the systems `PATH`.

## Install and Setup

### Configuration

Configuration is done through environment variables or a .env file. Copy the .env.defaults as a starting point.

```
TYPE=[devnet, testnet or mainnet]
PASS=<a passphrase to encrypt and decrypt account info>

NODE=<path to the Algorand node **>
DATA=<path to the collections data ***>

APP_PROG=<name of the collection contract>
APP_CLEAR=<name of the clear state contract>
APP_MANAGER=<default name for the managing account>

META_COLLECT=<name of the collection>
META_COLLECT_TAG=<two letter tag as in the app****>
META_COLLECT_MAXCOUNT=<Max amount as in the app****>

PINATA_KEY=<your-api-key>
PINATA_SEC=<your-api-secret>
```

** The path to the node needs: - To be writable - Hold the required tools - Hold network genesis files. \*** See below for the required structure and files.
\*\*\*\* Needs to match the information in the application contract.

### Algorand Node

`NODE` defined as `./node`

To interact with the Algorand network the toolchain needs access to a locally installed. The toolchain has been tested with the node software installed using the updater-script. ([algorand-updater-script])

Executables need to be at `./node` (NODE)
Genesis Files: need to be at `./node/genesisfiles`

Depending on the network in the configuration, one folder will be made within `./node` and be named `devnet-data`, `mainnet-data`, or `testnet-data`.

### Collections Data

`DATA` defined as `./assets`

The toolchain expects the following structure:
(with default contract files shown down below)

```
./assets
./assets/collection.png

./assets/accounts

./assets/contracts
./assets/contracts/clear.py
./assets/contracts/collection.py
./assets/contracts/utility.py

./assets/images
./assets/images/<name-x>.png
...
./assets/images/<name-xx>.png
```

The account is generated when going through the steps. As a backup, the mnemonic recovery words will be shown (only on creation), and an encrypted key file based on web3 standard is created in the accounts folder.

The contracts can be copied over from this repository but need minor tweaks to fit your collection.

###### Maximum Assets

By default, the contract can not handle more than 64512, if you desire to lower this value to assure potential buyers no extra assets will be added to the collection the following edit is required:

Inside collection.py on line 72

```py
        # Is smaller then max
        index.load() < Int(64512),
```

change to

```py
        # Is smaller then max
        index.load() < Int(<META_COLLECT_MAXCOUNT>),
```

###### Unit Name validation

The contract will validate the format of the unit name in the ASA. In short: The unit name of an ASA on the Algorand network has a max length of 8 bytes, and the contract can hand out a max of 64512 ids. The result is that`Nr#64512` is the max value possible; 2-byte tag, 1-byte separator, 5 bytes for the id.

If you are not using the defaults, the following edits need to happen:

Inside utility.py on line 42

```py
def get_asset_unit_name(arg):
...
        out.store(Bytes("Nr#")),
        tmp.store(uint_to_bytes(arg)),
        For(num.store(Int(5) - Len(tmp.load())) ...
...
```

Values here need to match up with `META_COLLECT_TAG` and `META_COLLECT_MAXCOUNT` as they relate to the unit name verification. (where the length of `MAXCOUNT` string is 5 or less)

```py
        out.store(Bytes("<META_COLLECT_TAG>#")),
        tmp.store(uint_to_bytes(arg)),
        For(num.store(Int(<len(str(META_COLLECT_MAXCOUNT))>)  ...
```

## Usage

### Usage Routine

Step by step setup and collection publication routine. While the goal is to automate things. At this point, steps are split up for manual verification and corrections when needed.

**ToDo: Varify IPFS Uploads: download and check hashes**
**Verification needs to be done manually for now**

> 0. Setup the configuration

```
setup the environment variables
avoid storing pass-phrase in files
```

> 1. Create and start a network node

```
go run ./cmd/publish network create
go run ./cmd/publish network start
```

> 2. Create and set up the manager account.
>    - use info to check on account balance
>    - on a `devnet` node it funds the account

```
go run ./cmd/publish account create
go run ./cmd/publish account info
```

> 3. Build and push the collection contract

```
go run ./cmd/publish contract build
go run ./cmd/publish contract push
```

> 4. Setup the ASA Data for the contract

```
go run ./cmd/publish contract image
go run ./cmd/publish contract meta
```

> 5. Setup the ASA Data for the assets

```
go run ./cmd/publish assets image
go run ./cmd/publish assets meta
```

> 6. Build transactions and mint assets

```
go run ./cmd/publish assets mint
```

### Extra commands

Network management

```
go run ./cmd/publish network create
go run ./cmd/publish network destroy
go run ./cmd/publish network status
go run ./cmd/publish network start
go run ./cmd/publish network stop
```

Test Pinata API-Keys

```
go run ./cmd/publish asset
:: Testing Pinata api keys on https://api.pinata.cloud/data/testAuthentication
>> Response: Congratulations! You are communicating with the Pinata API!
```

[algorand-install]: https://developer.algorand.org/docs/run-a-node/setup/install/
[algorand-updater-script]: https://developer.algorand.org/docs/run-a-node/setup/install/#installation-with-the-updater-script
