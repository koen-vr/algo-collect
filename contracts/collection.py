from pyteal import *
from utility import *

# Implements a basic byte array example
def approval_program():
    
    # Index in to the global array
    index = ScratchVar(TealType.uint64)
    # Key for the global key value pair
    string = ScratchVar(TealType.bytes)

    # [x, y, z] in to the global array
    key = ScratchVar(TealType.uint64)
    idx = ScratchVar(TealType.uint64)
    bit = ScratchVar(TealType.uint64)

    # Access global state for index key
    key_value = App.globalGet(string.load())
    key_has_value = App.globalGetEx(Int(0), string.load())

    # Check to see if the caller is the collection manager
    is_manager = App.localGet(Txn.sender(), Bytes("manager")) > Int(0)

    # Grab the first application argument from the call
    store_index = index.store(Btoi(Txn.application_args[1]))

    # Convert the index to [x, y, z] keys
    convert_to_keys = Seq([
        key.store(index.load() / Int(1008)),
        idx.store((index.load() % Int(1008)) / Int(8)),
        bit.store((index.load() % Int(1008)) % Int(8))
    ])

    # Convert the main index to a string for the global store
    convert_to_string = Seq([
        string.store(uint_to_bytes(key.load()))
    ])

    # If the key has no value initialize all bits to 0
    initialize_key = If(
        Not(key_has_value.hasValue()),
        App.globalPut(
            string.load(),
            Bytes("base16", "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
        )
    )

    # TODO assert that the asset can not be destroyed?
    # TODO assert is a valid collection asset? (nft)
    do_reserve_index = If(And(
            index.load() < Int(64512),
            Gtxn[0].config_asset() == Int(0),
            Gtxn[0].type_enum() == TxnType.AssetConfig,
            GetBit(GetByte(key_value, idx.load()), bit.load()) == Int(0),
            BytesEq(Gtxn[0].config_asset_unit_name(), get_asset_unit_name(index.load()))
        ),
        Seq([
            App.globalPut(
                string.load(),
                SetByte(
                    key_value,
                    idx.load(),
                    SetBit(
                        GetByte(key_value, idx.load()),
                        bit.load(),
                        Int(1)
                    )
                )
            ),
            Int(1)
        ]),
        Int(0)
    )

    handle_manager = If(And(
            is_manager,
            # Needs a transaction group
            Global.group_size() == Int(2),
            # Needs one new manager account
            Txn.accounts.length() == Int(1),
            # Needs that account to call opt-in
            Gtxn[0].sender() == Txn.accounts[1],
            # Needs the other call to be of type opt-in
            Gtxn[0].on_completion() == OnComplete.OptIn,
            # Needs the call to be to the same application
            Gtxn[0].application_id() == Txn.application_id(),
        ), 
        Seq([
            App.localPut(
                Txn.accounts[0],
                Bytes("manager"), 
                Int(0)
            ),
            App.localPut(
                Txn.accounts[1], 
                Bytes("manager"), 
                Int(1)
            ),
            Int(1)
        ]),
        Int(0)
    )

    handle_reserve = If(And(
            is_manager,
            # Require asset creation tx
            Global.group_size() == Int(2),
            # Require function call and index
            Txn.application_args.length() == Int(2)
        ),
        Seq([
            # Store the index
            # Convert index to key / idx / bit
            # Convert key to mapped string value
            store_index,
            convert_to_keys,
            convert_to_string,
            # Init values on key
            # Final try to reserve
            key_has_value,
            initialize_key,
            do_reserve_index
        ]),
        Int(0)
    )

    handle_create = Seq([
        # Set the app creator as the initial manager
        App.localPut(Int(0), Bytes("manager"), Int(1)),
        Int(1)
    ])

    handle_noop = Cond(
        [Txn.application_args[0] == Bytes("manager"), handle_manager],
        [Txn.application_args[0] == Bytes("reserve"), handle_reserve],
    )

    handle_optin = And(
        # Require a second call
        Global.group_size() == Int(2),
        # Needs to be an application call
        Gtxn[1].type_enum() == TxnType.ApplicationCall,
        # Needs the call to setup a new manager
        Gtxn[1].application_args[0] == Bytes("manager"),
        # Needs the call to be to the same application
        Gtxn[1].application_id() == Txn.application_id(),
        # Needs that call to set the sender as the manager
        Gtxn[1].accounts[1] == Txn.accounts[0]
    )

    handle_closeout = Seq([ 
        # Only close out if the sender is not the manager
        App.localGet(Txn.sender(), Bytes("manager")) == Int(0)
    ])

    # Disable updates, final no changes
    handle_updateapp = Seq([ Int(0) ])

    # Disable removal, this data lives for ever
    handle_deleteapp = Seq([ Int(0) ])

    program = Cond(
        [Txn.application_id() == Int(0), handle_create],
        [Txn.on_completion() == OnComplete.NoOp, handle_noop],
        [Txn.on_completion() == OnComplete.OptIn, handle_optin],
        [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
        [Txn.on_completion() == OnComplete.DeleteApplication, handle_updateapp],
        [Txn.on_completion() == OnComplete.UpdateApplication, handle_deleteapp],
    )
    return program

if __name__ == "__main__":
    print(compileTeal(approval_program(), Mode.Application, version=5))