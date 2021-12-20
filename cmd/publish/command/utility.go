package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	cfg "github.com/vecno-io/go-pyteal/config"
	net "github.com/vecno-io/go-pyteal/network"
)

var setup = cfg.Setup{}

func onInitialize(validate bool) error {
	if err := viper.Unmarshal(&setup); nil != err {
		return err
	}
	if err := cfg.OnInitialize(setup, validate); nil != err {
		return err
	}
	return nil
}

func getListOfFiles(ext, path string) ([]string, error) {
	list := make([]string, 0)
	info, err := os.Stat(path)
	if nil != err {
		return list[:], fmt.Errorf("get file list: %s", err)
	}
	if !info.IsDir() {
		return list[:], fmt.Errorf("get file list:: not a direcotry: %s", path)
	}

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ext {
			list = append(list, path)
		}
		return nil
	}); err != nil {
		return list[:], fmt.Errorf("get file list: walk tree: %s", err)
	}

	return list[:], nil
}

func getUnitName(nr uint32) string {
	str := fmt.Sprintf("%%s#%%0%dd", len(viper.GetString("META_COLLECT_MAXCOUNT")))
	return fmt.Sprintf(str, viper.GetString("META_COLLECT_TAG"), nr)
}

func getPinData(path string) (PinFileResponse, error) {
	res := PinFileResponse{}
	data, err := os.ReadFile(path)
	if err != nil {
		return res, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return PinFileResponse{}, err
	}
	return res, nil
}

func getAssetMeta(file string) (Arc3Asset, error) {
	asset := Arc3Asset{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return asset, err
	}
	if err := json.Unmarshal(data, &asset); nil != err {
		return asset, err
	}
	return asset, nil
}

func getApplicationId(name string) (string, error) {
	id, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.id", viper.GetString("DATA"), name))
	if err != nil {
		return "", err
	}
	return strings.Split(string(id), "\n")[0], nil
}

func getApplicationUrl(name string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json.pin", viper.GetString("DATA"), name))
	if err != nil {
		return "", err
	}
	pin := PinFileResponse{}
	if err := json.Unmarshal(data, &pin); nil != err {
		return "", err
	}
	return fmt.Sprintf("ipfs://%s", pin.IpfsHash), nil
}

func txnBuild(app uint64, file string, acc crypto.Account, params types.SuggestedParams) (uint64, error) {
	tx1, err := txnBuildAssetCreate(file, acc, params)
	if nil != err {
		return 0, fmt.Errorf("asset txn: %s", err)
	}
	tx2, err := txnBuildAppCall(app, file, acc, params)
	if nil != err {
		return 0, fmt.Errorf("app call txn: %s", err)
	}
	txn, err := txnBuildGroup(tx1, tx2, acc)
	if nil != err {
		return 0, fmt.Errorf("app group txn: %s", err)
	}
	txInfo, err := net.SendRawTransaction(txn)
	if nil != err {
		return 0, fmt.Errorf("send raw txn:  %s", err)
	}
	fmt.Println("Asset:", txInfo.AssetIndex)

	return txInfo.AssetIndex, nil
}

func txnBuildGroup(tx1, tx2 types.Transaction, ac crypto.Account) ([]byte, error) {
	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})
	if err != nil {
		return []byte{}, err
	}
	tx1.Group = gid
	tx2.Group = gid

	_, stx1, err := crypto.SignTransaction(ac.PrivateKey, tx1)
	if err != nil {
		return []byte{}, err
	}
	_, stx2, err := crypto.SignTransaction(ac.PrivateKey, tx2)
	if err != nil {
		return []byte{}, err
	}

	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)

	return signedGroup, nil
}

func txnBuildAppCall(app uint64, file string, ac crypto.Account, params types.SuggestedParams) (types.Transaction, error) {
	meta, err := getAssetMeta(file)
	if err != nil {
		return types.Transaction{}, err
	}
	id, err := strconv.ParseUint(meta.UnitName[3:], 10, 64)
	if err != nil {
		return types.Transaction{}, err
	}
	appArgs := [][]byte{
		[]byte("reserve"),
		i32tob(uint32(id)),
	}

	note := []byte(nil)
	accounts := []string{}
	foreignApps := []uint64{}
	foreignAssets := []uint64{}
	group := types.Digest{}
	lease := [32]byte{}
	rekeyTo := types.ZeroAddress

	// Note: Inconsitency in the api, needs a decoded address?
	sender, err := types.DecodeAddress(ac.Address.String())
	if err != nil {
		return types.Transaction{}, err
	}

	appTxn, err := future.MakeApplicationNoOpTx(
		app, appArgs, accounts, foreignApps, foreignAssets,
		params, sender, note, group, lease, rekeyTo,
	)
	if err != nil {
		return types.Transaction{}, err
	}

	return appTxn, nil
}

func txnBuildAssetCreate(file string, ac crypto.Account, params types.SuggestedParams) (types.Transaction, error) {
	pin, err := getPinData(fmt.Sprintf("%s.pin", file))
	if err != nil {
		return types.Transaction{}, err
	}
	hash, err := hashAsaFile(file)
	if err != nil {
		return types.Transaction{}, err
	}
	meta, err := getAssetMeta(file)
	if err != nil {
		return types.Transaction{}, err
	}

	unitName := meta.UnitName
	assetName := meta.Name
	assetPath := fmt.Sprintf("ipfs://%s#arc3", pin.IpfsHash)

	note := []byte(nil)
	total := uint64(1)
	frozen := false
	decimals := uint32(0)

	address := ac.Address.String()
	manager := types.ZeroAddress.String()

	freeze := types.ZeroAddress.String()
	reserve := types.ZeroAddress.String()
	clawback := types.ZeroAddress.String()

	createTxn, err := future.MakeAssetCreateTxn(
		address, note, params, total, decimals,
		frozen, manager, reserve, freeze, clawback,
		unitName, assetName, assetPath, "",
	)
	if nil != err {
		return types.Transaction{}, err
	}

	// Hack, Bug MakeAssetCreateTxnm needs a string
	copy(createTxn.AssetParams.MetadataHash[:], hash[:32])

	return createTxn, nil
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * (3 - i))) & 0xff)
	}
	return r
}
