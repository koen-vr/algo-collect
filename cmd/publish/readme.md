# algo-collect publisher

A simple demo project to act as a base for publishing collections using the algo-collect contract.

## Requirements

- Linux or macOS
- Golang version 1.17.0 or higher
- Python 3. The scripts assumes the Python executable is called `python3`.
- The [Algorand Node software][algorand-install] installed using the updater script; with the enviorment variable `ALGORAND_DATA` set and the toolset added to the systems `PATH`.

[algorand-install]: https://developer.algorand.org/docs/run-a-node/setup/install/

## Node install

The tools will configure private, testnet and mainnet nodes depending on what you do; to do this it does require the node software to be available.

Defailt steps to get it installed:

> Create a temporary folder to hold the install package and files.

```
mkdir ~/node
cd ~/node
```

Download the install/update script

```
wget https://raw.githubusercontent.com/algorand/go-algorand-doc/master/downloads/installers/update.sh
```

```
chmod 544 update.sh
```

```
./update.sh -i -c stable -p ~/node -d ~/node/data -n
```

## Usage

Network managment

```
go run ./cmd/publish network create
go run ./cmd/publish network destroy
go run ./cmd/publish network status
go run ./cmd/publish network start
go run ./cmd/publish network stop
```

Wallet managment

```
go run ./cmd/publish account info (--name my-account)
go run ./cmd/publish account create (--name my-account)
```

Test pinata api keys

```
go run ./cmd/publish asset
```

## Usage Routines

Most steps are split in to two parts, one where all data is prepaird and verified and then a second step to publish the data on to the target network. Seems reduent but it is safer and easier to catch and correct issues before things are put on to the blockchain.

**ToDo: Varify IPFS Uploads: download and check hashes**

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

> 2. Create and setup the manager account.
>    - using info to check on acount balance
>    - on a `devnet` node it funds the acocunt

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

**TODO: Distrebute NFTs to address list**
