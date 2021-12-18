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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	//Pinata.AddCommand(pinataDataCmd)
	Pinata.AddCommand(pinataImageCmd)
}

type PinFileResponse struct {
	PinSize   uint64 `json:"PinSize"`
	IpfsHash  string `json:"IpfsHash"`
	Timestamp string `json:"Timestamp"`
}

type PinTestResponse struct {
	Message string `json:"message"`
}

var Pinata = &cobra.Command{
	Use:   "pinata",
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
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := onInitialize(true); err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}
	},
}

var pinataImageCmd = &cobra.Command{
	Use:   "images",
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

func imageUploadPin(url, path string) error {
	client := &http.Client{}

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

	// Note: clean up the json output for humans
	body, err = json.MarshalIndent(pin, "", "  ")
	if err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s.json", path[:len(path)-4]))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.Write(body); nil != err {
		return err
	}

	return nil
}
