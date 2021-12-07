package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/vecno-io/algo-collection/shared/account"
	"github.com/vecno-io/algo-collection/shared/client"
	exec "github.com/vecno-io/algo-collection/shared/execute"
)

var CreateappCmd = &cobra.Command{
	Use:   "createapp",
	Short: "Create the collection app",
	Long:  `Compiles and publishes the collection program to the private network.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		approvalProg, err := ioutil.ReadFile("./contracts/collection.teal")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		clearProg, err := ioutil.ReadFile("./contracts/clear.teal")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("--- Loading frst account")
		acc, err := account.LoadFromFile("./ac1.frag")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("--- App creation transaction")
		cln, err := client.MakeAlgodClient()
		if err != nil {
			fmt.Printf("failed to make algod client: %s\n", err)
			os.Exit(1)
		}

		txnParams, err := cln.SuggestedParams().Do(context.Background())
		if err != nil {
			fmt.Printf("failed to get transaction params: %s\n", err)
			os.Exit(1)
		}

		apx, err := cln.TealCompile(approvalProg).Do(context.Background())
		if err != nil {
			fmt.Printf("failed to compile approval program: %s\n", err)
			os.Exit(1)
		}
		ap, err := base64.StdEncoding.DecodeString(apx.Result)
		if err != nil {
			fmt.Printf("failed to decode approval program: %s\n", err)
			os.Exit(1)
		}
		err = logic.CheckProgram(ap, make([][]byte, 0))
		if nil != err {
			fmt.Fprintln(os.Stderr, err)
			fmt.Printf("check on application program failed: %s\n", err)
			os.Exit(1)
		}

		cpx, err := cln.TealCompile(clearProg).Do(context.Background())
		if err != nil {
			fmt.Printf("failed to compile clear program: %s\n", err)
			os.Exit(1)
		}
		cp, err := base64.StdEncoding.DecodeString(cpx.Result)
		if err != nil {
			fmt.Printf("failed to decode clear program: %s\n", err)
			os.Exit(1)
		}
		err = logic.CheckProgram(cp, make([][]byte, 0))
		if nil != err {
			fmt.Fprintln(os.Stderr, err)
			fmt.Printf("check on clear program failed: %s\n", err)
			os.Exit(1)
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
			fmt.Printf("failed to create app creation transaction: %s\n", err)
			os.Exit(1)
		}

		_, stx1, err := crypto.SignTransaction(acc.PrivateKey, tx1)
		if err != nil {
			fmt.Printf("failed to sign transaction: %s\n", err)
			return
		}

		pendingTxID, err := cln.SendRawTransaction(stx1).Do(context.Background())
		if err != nil {
			fmt.Printf("failed to send transaction: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("transaction send: %s\n", pendingTxID)
		txInfo, err := client.WaitForConfirmation(cln, pendingTxID, 24, context.Background())
		if err != nil {
			fmt.Printf("failed to confirm transaction: %s\n", err)
			os.Exit(1)
		}
		if len(txInfo.PoolError) > 0 {
			fmt.Printf("error while confirm transaction: %s\n", txInfo.PoolError)
			os.Exit(1)
		}
		fmt.Printf("Application deployed with id: %d\n", txInfo.ApplicationIndex)

		// Save the application id to file

		jsonString, _ := json.Marshal(txInfo.ApplicationIndex)
		ioutil.WriteFile("./app.frag", jsonString, os.ModePerm)

		out, err = exec.List([]string{
			"-c", fmt.Sprintf("goal app read -d ./net1/primary --app-id %d --guess-format --local --from %s", txInfo.ApplicationIndex, sender.String()),
		})
		if len(out) > 0 {
			fmt.Println()
			fmt.Println(out)
		}
		if nil != err {
			fmt.Printf("failed to read app state: %s\n", err)
			os.Exit(1)
		}
	},
}
