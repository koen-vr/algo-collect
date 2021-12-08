# algo-collection

Demo's a collection manager for NFT's represented by an ASA.

This demo stores private keys in plain text and should not be used in production.

This project **hasn't been security audited** and should not be used in a production environment.

## Brief

The collection is a list of unique IDs assigned to the NFTs. The manager tracks free IDs and specific NFTs store their ID in the 'unitname' parameter. The creation of an asset and reservation of an ID is wrapped in an atomic transaction to make a secure link.

## Requirements

* Linux or macOS
* Golang version 1.17.0 or higher
* Python 3. The scripts assumes the Python executable is called `python3`.
* The [Algorand Node software](https://developer.algorand.org/docs/run-a-node/setup/install/). A private network is used, hence there is no need to sync up MainNet or TestNet. `goal` is assumed to be in the PATH.

## Setup

To install all required packages, run: 
```bash
python3 -m pip install -r requirements.txt
```

## Usage

TODO: Refresh to golang tools
~~~
go run ./cmd/collection network start
go run ./cmd/collection create wallet
go run ./cmd/collection create app
go run ./cmd/collection create asset
go run ./cmd/collection network stop
~~~

## Credits

Based off and inspired by: [algo-arrays](https://github.com/gidonkatten/algo-arrays)
