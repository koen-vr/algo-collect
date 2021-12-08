package contract

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/vecno-io/algo-collection/shared/account"

	exec "github.com/vecno-io/algo-collection/shared/execute"
	net "github.com/vecno-io/algo-collection/shared/network"
)

func Deploy() error {
	fmt.Println("### Creating the collection application")

	fmt.Println("--- Compile approval prog")
	out, err := exec.List([]string{
		"-c", "python3 ./contracts/collection.py > ./contracts/collection.teal",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return fmt.Errorf("contract: compiling approval failed: %s", err)
	}
	approvalProg, err := ioutil.ReadFile("./contracts/collection.teal")
	if err != nil {
		return fmt.Errorf("contract: reading approval failed: %s", err)
	}

	fmt.Println("--- Compile clear prog")
	out, err = exec.List([]string{
		"-c", "python3 ./contracts/clear.py > ./contracts/clear.teal",
	})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return fmt.Errorf("contract: compiling clear failed: %s", err)
	}
	clearProg, err := ioutil.ReadFile("./contracts/clear.teal")
	if err != nil {
		return fmt.Errorf("contract: reading clear failed: %s", err)
	}

	fmt.Println("--- Loading main account")
	acc, err := account.GetMainAccount()
	if err != nil {
		return fmt.Errorf("contract: reading main account failed: %s", err)
	}

	fmt.Println("--- App creation transaction")
	cln, err := net.MakeClient()
	if err != nil {
		return fmt.Errorf("contract: failed to make algod client: %s", err)
	}

	txnParams, err := cln.SuggestedParams().Do(context.Background())
	if err != nil {
		return fmt.Errorf("contract: failed to get transaction params: %s", err)
	}

	apx, err := cln.TealCompile(approvalProg).Do(context.Background())
	if err != nil {
		return fmt.Errorf("contract: failed to compile approval program: %s", err)
	}
	ap, err := base64.StdEncoding.DecodeString(apx.Result)
	if err != nil {
		return fmt.Errorf("contract: failed to decode approval program: %s", err)
	}
	err = logic.CheckProgram(ap, make([][]byte, 0))
	if nil != err {
		return fmt.Errorf("contract: check on application program failed: %s", err)
	}

	cpx, err := cln.TealCompile(clearProg).Do(context.Background())
	if err != nil {
		return fmt.Errorf("contract: failed to compile clear program: %s", err)
	}
	cp, err := base64.StdEncoding.DecodeString(cpx.Result)
	if err != nil {
		return fmt.Errorf("contract: failed to decode clear program: %s", err)
	}
	err = logic.CheckProgram(cp, make([][]byte, 0))
	if nil != err {
		return fmt.Errorf("contract: check on clear program failed: %s", err)
	}

	optIn := true
	globalSchema := types.StateSchema{
		NumUint:      0,
		NumByteSlice: 64,
	}
	localSchema := types.StateSchema{
		NumUint:      2,
		NumByteSlice: 0,
	}
	appArgs := [][]byte{}
	accounts := []string{}
	foreignApps := []uint64{}
	foreignAssets := []uint64{}
	sender := acc.Address
	note := []byte{}
	group := types.Digest{}
	lease := [32]byte{}
	rekeyTo := types.ZeroAddress

	// Even with opt-in to true it needs to be called seperataly?
	tx1, err := future.MakeApplicationCreateTxWithExtraPages(
		optIn, ap, cp, globalSchema, localSchema,
		appArgs, accounts, foreignApps, foreignAssets, txnParams,
		sender, note, group, lease, rekeyTo, 1,
	)
	// Double check to make sure
	tx1.OnCompletion = types.OptInOC

	if err != nil {
		return fmt.Errorf("contract: failed to create app creation transaction: %s", err)
	}

	_, stx1, err := crypto.SignTransaction(acc.PrivateKey, tx1)
	if err != nil {
		return fmt.Errorf("contract: failed to sign transaction: %s", err)
	}

	pendingTxID, err := cln.SendRawTransaction(stx1).Do(context.Background())
	if err != nil {
		return fmt.Errorf("contract: failed to send transaction: %s", err)
	}

	fmt.Printf("transaction send: %s\n", pendingTxID)
	txInfo, err := net.WaitForConfirmation(cln, pendingTxID, 24, context.Background())
	if err != nil {
		return fmt.Errorf("contract: failed to confirm transaction: %s", err)
	}
	if len(txInfo.PoolError) > 0 {
		return fmt.Errorf("contract: error while confirm transaction: %s", txInfo.PoolError)
	}
	fmt.Printf("Application deployed with id: %d\n", txInfo.ApplicationIndex)

	// Save the application id to file

	if err := saveToFile(
		fmt.Sprintf("%s/app1.frag", net.Path()),
		txInfo.ApplicationIndex,
	); err != nil {
		return fmt.Errorf("contract: failed to save app: %s", err)
	}

	out, err = exec.List([]string{"-c", fmt.Sprintf(
		"goal app read -d %s --app-id %d --guess-format --local --from %s",
		net.NodePath(), txInfo.ApplicationIndex, sender.String(),
	)})
	if len(out) > 0 {
		fmt.Println()
		fmt.Println(out)
	}
	if nil != err {
		return fmt.Errorf("contract: failed to read app state: %s", err)
	}
	return nil
}

func GetMainApp() (uint64, error) {
	return loadFromFile(fmt.Sprintf("%s/app1.frag", net.Path()))
}
