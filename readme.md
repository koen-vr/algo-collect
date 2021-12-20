# algo-collect

Demo's a collection manager for NFT's represented by an ASA.

This project **hasn't been security audited** and should only be used with caution and understanding.

## Brief

The collection is a list of unique IDs assigned to ASAs within a grouped atomic transaction.

## Requirements

- Linux or macOS
- Golang version 1.17.0 or higher
- Python 3. The tools assume the Python executable is called `python3`.
- The [Algorand Node software](https://developer.algorand.org/docs/run-a-node/setup/install/)- A private network is used, hence there is no need to sync up MainNet or TestNet. `goal` is assumed to be in the PATH.

## Setup

To install all required packages, run:

```bash
python3 -m pip install -r requirements.txt
```

## The collection

By using this contract a collection and the maximum amount of assets within a collection is set in an Algorand application.

By looking at the transaction history of the collections application one can find all other assets within the collection.

By looking at the assets creation and the related group transaction proof is provided that the asset was minted within the collection.

## Demo Contract and Code

A demo can be found in the `cmd/publish` folder, this tool will take assets in a folder build metadata for them, and publish them onto the blockchain. A full guide can be found in the readme.md file of the demo.

## Credits

Based off and inspired by: [algo-arrays](https://github.com/gidonkatten/algo-arrays)
