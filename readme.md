# algo-arrays

This project **hasn't been security audited** and should not be used in a production environment.

Further, it has a hard-coded unfinished option just to demo off the idea.
~~~
BytesEq(Gtxn[0].config_asset_unit_name(), Bytes("APP#1")),
~~~

## Requirements

* Linux or macOS
* Python 3. The scripts assumes the Python executable is called `python3`.
* The [Algorand Node software](https://developer.algorand.org/docs/run-a-node/setup/install/). A private network is used, hence there is no need to sync up MainNet or TestNet. `goal` is assumed to be in the PATH.

## Setup

To install all required packages, run: 
```bash
python3 -m pip install -r requirements.txt
```

## Usage

There are a number of bash script files which execute the goal commands for you. You mush be in the same directory as the script file for it to run correctly.

They should be run in the following order:
1. **startnet.sh**: Sets up private network
2. **createapp.sh**: Compiles the PyTeal files to TEAL and deploys the stateful smart contract. There are 3 different **createapp.sh** files - run the one depending on which example you would like to try out.
3. **createasset.sh**: Calls the app to register an id and creates the asset in a grouped transaction.
3. **setmanager.sh**: Calls the app to change the manager and performs the `opt in` for the new manager in a grouped transaction as required. (breaks asset creation script as it uses the old manager)
4. **stopnet.sh**: Delete private network

Note that the scripts assume the application has ID `APP_ID=1`.
In particular, no transactions should be made before calling `createapp.sh`.

Scripts in `3` can be continuosly run (as long as the corresponding app has been created).

## Credits

Based off and ispired by: (algo-arrays)[https://github.com/gidonkatten/algo-arrays]
