package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func MakeAlgodClient() (*algod.Client, error) {
	addr, err := getFirstLineFromFile("./net1/Primary/algod.net")
	if err != nil {
		return nil, fmt.Errorf("read network file: %v", err)
	}

	token, err := getFirstLineFromFile("./net1/Primary/algod.token")
	if err != nil {
		return nil, fmt.Errorf("read token file: %v", err)
	}
	return algod.MakeClient("http://"+addr, token)
}

func WaitForConfirmation(c *algod.Client, txid string, waitRounds uint64, ctx context.Context, headers ...*common.Header) (txInfo models.PendingTransactionInfoResponse, err error) {
	response, err := c.Status().Do(ctx, headers...)
	if err != nil {
		return
	}

	lastRound := response.LastRound
	currentRound := lastRound + 1

	for {
		// Check that the `waitRounds` has not passed
		if waitRounds > 0 && currentRound > lastRound+waitRounds {
			err = fmt.Errorf("wait for transaction id %s timed out", txid)
			return
		}

		txInfo, _, err = c.PendingTransactionInformation(txid).Do(ctx, headers...)
		if err != nil {
			return
		}

		// The transaction has been confirmed
		if txInfo.ConfirmedRound > 0 {
			return
		}

		// Wait until the block for the `currentRound` is confirmed
		response, err = c.StatusAfterBlock(currentRound).Do(ctx, headers...)
		if err != nil {
			return
		}

		// Increment the `currentRound`
		currentRound += 1
	}
}

func getFirstLineFromFile(file string) (string, error) {
	addrStr, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(addrStr), "\n")
	return lines[0], err
}
