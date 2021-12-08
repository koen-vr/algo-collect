package asset

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/vecno-io/algo-collection/shared/account"
	"github.com/vecno-io/algo-collection/shared/contract"

	exec "github.com/vecno-io/algo-collection/shared/execute"
	net "github.com/vecno-io/algo-collection/shared/network"
)

func Setup() error {
	fmt.Println("### Creating an asset for the collection")

	fmt.Println("--- Asset creation transaction")
	cln, err := net.MakeClient()
	if err != nil {
		return fmt.Errorf("asset: asset: failed to make algod client: %s", err)
	}

	txnParams, err := cln.SuggestedParams().Do(context.Background())
	if err != nil {
		return fmt.Errorf("asset: failed to get transaction params: %s", err)
	}

	// ToDo replace these with the account loads

	ac1, err := account.GetMainAccount()
	if err != nil {
		return fmt.Errorf("asset: failed to get the main account: %s", err)
	}
	ac2, err := account.GetUserAccount()
	if err != nil {
		return fmt.Errorf("asset: failed to get the user account: %s", err)
	}

	// Setup the Asset Creation Transaction

	unitIdx := uint32(64511)
	unitName := fmt.Sprintf("#%d", unitIdx)
	// unitIdx := uint32(1)
	// unitName := fmt.Sprintf("#0000%d", unitIdx)
	// unitIdx := uint32(255)
	// unitName := fmt.Sprintf("#00%d", unitIdx)

	assetName := "latinum"
	assetHash := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	assetPath := "https://path/to/my/asset/details"

	note := []byte(nil)
	total := uint64(1)
	frozen := false
	decimals := uint32(0)

	address := ac2.Address.String()
	manager := types.ZeroAddress.String()

	freeze := types.ZeroAddress.String()
	reserve := types.ZeroAddress.String()
	clawback := types.ZeroAddress.String()

	tx1, err := future.MakeAssetCreateTxn(
		address, note, txnParams, total, decimals,
		frozen, manager, reserve, freeze, clawback,
		unitName, assetName, assetPath, assetHash,
	)
	if err != nil {
		return fmt.Errorf("asset: failed to create asset creation transaction: %s", err)
	}

	// Setup the Index Reservation Transaction

	appIdx, err := contract.GetMainApp()
	if err != nil {
		return fmt.Errorf("asset: fail to get main app id: %s", err)
	}

	appArgs := [][]byte{
		[]byte("reserve"),
		i32tob(unitIdx),
	}
	accounts := []string{}
	foreignApps := []uint64{}
	foreignAssets := []uint64{}
	group := types.Digest{}
	lease := [32]byte{}
	rekeyTo := types.ZeroAddress

	// Note: Inconsitency in the api, needs a decoded address?
	sender, err := types.DecodeAddress(ac1.Address.String())
	if err != nil {
		return fmt.Errorf("asset: failed to decode address: %s", err)
	}

	tx2, err := future.MakeApplicationNoOpTx(
		appIdx, appArgs, accounts, foreignApps, foreignAssets,
		txnParams, sender, note, group, lease, rekeyTo,
	)
	if err != nil {
		return fmt.Errorf("asset: failed to create application call transaction: %s", err)
	}

	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})
	if err != nil {
		return fmt.Errorf("asset: failed to create group id: %s", err)
	}
	tx1.Group = gid
	tx2.Group = gid

	// Note Account 1 is the contract owner, it does the contract call. The minter
	// can be annyone so the collection logic is not bount by the 1k minting limit.
	_, stx1, err := crypto.SignTransaction(ac2.PrivateKey, tx1)
	if err != nil {
		return fmt.Errorf("asset: failed to sign transaction one: %s", err)
	}
	_, stx2, err := crypto.SignTransaction(ac1.PrivateKey, tx2)
	if err != nil {
		return fmt.Errorf("asset: failed to sign transaction two: %s", err)
	}

	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)

	pendingTxID, err := cln.SendRawTransaction(signedGroup).Do(context.Background())
	if err != nil {
		return fmt.Errorf("asset: failed to send transaction: %s", err)
	}

	txInfo, err := net.WaitForConfirmation(cln, pendingTxID, 24, context.Background())
	if err != nil {
		return fmt.Errorf("asset: failed to confirm transaction: %s", err)
	}
	if len(txInfo.PoolError) > 0 {
		return fmt.Errorf("asset: error while confirm transaction: %s", txInfo.PoolError)
	}

	fmt.Printf("Asset Created deployed with id: %d", txInfo.AssetIndex)

	out, err := exec.List([]string{"-c", fmt.Sprintf(
		"goal app read -d %s --app-id %d --guess-format --global",
		net.NodePath(), appIdx,
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return fmt.Errorf("asset: failed to read app state: %s", err)
	}
	return nil
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * (3 - i))) & 0xff)
	}
	return r
}
