package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	acc "github.com/vecno-io/go-pyteal/account"
	net "github.com/vecno-io/go-pyteal/network"
)

func init() {
	Assets.AddCommand(assetsMintCmd)
	Assets.AddCommand(assetsMetaCmd)
	Assets.AddCommand(assetsImageCmd)
}

type PinFileResponse struct {
	PinSize   uint64 `json:"PinSize"`
	IpfsHash  string `json:"IpfsHash"`
	Timestamp string `json:"Timestamp"`
}

type PinTestResponse struct {
	Message string `json:"message"`
}

var Assets = &cobra.Command{
	Use:   "assets",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		const url = "https://api.pinata.cloud/data/testAuthentication"

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		req.Header.Set("pinata_api_key", viper.GetString("PINATA_KEY"))
		req.Header.Set("pinata_secret_api_key", viper.GetString("PINATA_SEC"))

		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		msg := PinTestResponse{}
		if err := json.Unmarshal(body, &msg); nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf(":: Testing Pinata api keys on %s \n", url)
		fmt.Printf(">> Response: %s \n\n", msg.Message)

		cmd.HelpFunc()(cmd, args)
	},
}

var assetsMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := fmt.Sprintf("%s/images", viper.GetString("DATA"))
		fmt.Println(":: Building meta data for:", root)

		list, err := getListOfFiles(".pin", root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		// filter down img pins
		pinList := make([]string, 0)
		for _, path := range list {
			file := path[:len(path)-4]
			if filepath.Ext(file) != ".json" {
				pinList = append(pinList, path)
			}
		}

		fmt.Printf(">> Found %d assets to build metadata for.\n", len(pinList))
		for idx, path := range pinList {
			fmt.Printf(">> Build > %s.json ", path[:len(path)-4])
			if err := metaBuildData(root, path, uint32(idx+1)); nil != err {
				fmt.Printf(">> Error: \n>>>> %s\n", err)
			} else {
				fmt.Println(">> Ok")
			}
			// TODO Push Meta to ipfs
			fmt.Printf(">> Push > %s.json ", path[:len(path)-4])
			if err := metaPushData(root, path); nil != err {
				fmt.Printf(">> Error: \n>>>> %s\n", err)
			} else {
				fmt.Println(">> Ok")
			}
		}
	},
}

var assetsImageCmd = &cobra.Command{
	Use:   "image",
	Short: "upload all images to pinata",
	Long:  `Uploads all images to the pinata server and stores the ipfs hashes to disk.`,
	Run: func(cmd *cobra.Command, args []string) {
		const url = "https://api.pinata.cloud/pinning/pinFileToIPFS"

		root := fmt.Sprintf("%s/images", viper.GetString("DATA"))
		fmt.Println(":: Upload from folder to pinata:", root)

		list, err := getListOfFiles(".png", root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf(">> Found %d png images to upload and pin.\n", len(list))
		for _, path := range list {
			fmt.Printf(">> Upload > %s ", path)
			if err := imageUploadPin(url, path); nil != err {
				fmt.Printf(">> Err: \n >>>> %s\n", err)
			} else {
				fmt.Println(">> Ok")
			}
		}
	},
}

var assetsMintCmd = &cobra.Command{
	Use:   "mint",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := fmt.Sprintf("%s/images", viper.GetString("DATA"))
		fmt.Println(":: Build ASA transactions in folder:", root)

		// filter down meta pins
		list, err := getListOfFiles(".pin", root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		pinList := make([]string, 0)
		for _, path := range list {
			path = path[:len(path)-4]
			if filepath.Ext(path) == ".json" {
				pinList = append(pinList, path)
			}
		}

		app, err := getApplicationId("collection")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		id, err := strconv.ParseUint(app, 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		params, err := net.MakeTxnParams()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		manager, err := acc.Load(viper.GetString("APP_MANAGER"), viper.GetString("PASS"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		idList := make([]uint64, 0)
		fmt.Printf(">> Found %d assets to prepair create transactions for.\n", len(pinList))
		for idx, path := range pinList {
			fmt.Printf(">> Build > %d > %s \n", idx+1, path)
			asset, err := txnBuild(id, path, manager, params)
			if nil != err {
				fmt.Printf(">> Err: \n>>>> %s\n", err)
			} else {
				fmt.Println(">> Ok")
				idList = append(idList, asset)
			}
		}

		// TEMP Just dump the asset id's to a file
		out, err := json.MarshalIndent(idList, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		outFile, err := os.Create(fmt.Sprintf("%s/assets.json", viper.GetString("DATA")))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
		defer outFile.Close()
		if _, err := outFile.Write(out); nil != err {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

func imageUploadPin(url, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", viper.GetString("PINATA_KEY"))
	req.Header.Set("pinata_secret_api_key", viper.GetString("PINATA_SEC"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	pin := PinFileResponse{}
	if err := json.Unmarshal(body, &pin); nil != err {
		return err
	}
	if 0 >= len(pin.IpfsHash) {
		return fmt.Errorf("invalid hash returned")
	}

	// Note: cleans up the json output for humans
	body, err = json.MarshalIndent(pin, "", "  ")
	if err != nil {
		return err
	}
	out, err := os.Create(fmt.Sprintf("%s.pin", path[:len(path)-4]))
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := out.Write(body); nil != err {
		return err
	}

	return nil
}
