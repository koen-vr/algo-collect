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

Most steps are split in to two parts, one where all data is prepaird and verified and then a second step to publish the data on to the target network. Seems reduent but it is safer and easier to catch and correct issues before things are put on to the blockchain.

> 1. Start and stop the network node

```
go run ./cmd/publish network start testnet
go run ./cmd/publish network stop
```

> 2. Verify and deploy the collectors contract

```
go run ./cmd/publish deploy build
go run ./cmd/publish deploy publish
```

> 3. Verify and publish the image files

```
go run ./cmd/publish image build
go run ./cmd/publish image publish
```

> 4. Verify and publish the metadata files

```
go run ./cmd/publish meta build
go run ./cmd/publish meta publish
```

> 5. Verify transaction and mint the nfts

```
go run ./cmd/publish nft setup
go run ./cmd/publish nft publish
```

**TODO: Distrebute NFTs to address list**