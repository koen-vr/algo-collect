package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Arc3App ASA Metadata Info (SA042-ASA)
type Arc3App struct {
	App      string `json:"app"`
	Name     string `json:"name"`
	UnitMax  string `json:"unitmax"`
	UnitName string `json:"unitname"`

	Image           string `json:"image"`
	Image_integrity string `json:"image_integrity"`
	Image_mimetype  string `json:"image_mimetype"`
}

// Arc3Asset ASA Metadata Info (SA076-ASA)
type Arc3Asset struct {
	Name     string `json:"name"`
	UnitName string `json:"unitname"`

	Image           string `json:"image"`
	Image_integrity string `json:"image_integrity"`
	Image_mimetype  string `json:"image_mimetype"`

	Properties SAAsset `json:"properties"`
}

// SAAsset Asset Properies (SA101)
type SAAsset struct {
	// Traits     SATraits     `json:"traits"`
	Collection SACollection `json:"collection"`
}

// SATraits Market Traits (SA069)
// type SATraits map[string]string

// SACollection ASA Collection (SA042)
type SACollection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Metadata string `json:"metadata"`
}

func hashAsaFile(path string) ([]byte, error) {
	hash := sha256.New()
	data, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer data.Close()
	if _, err := io.Copy(hash, data); err != nil {
		return []byte{}, err
	}
	return hash.Sum(nil), nil
}

func hashImageFile(path string) (string, error) {
	hash := sha256.New()
	data, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer data.Close()
	if _, err := io.Copy(hash, data); err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("sha256-%s", str), nil
}

func metaPushData(root, path string) error {
	_, file := filepath.Split(path)
	file = fmt.Sprintf("%s.json", file[:len(file)-4])
	return metaPushContract(root, file)
}

func metaBuildData(dest, path string, idx uint32) error {
	appId, err := getApplicationId("collection")
	if err != nil {
		return fmt.Errorf("failed to collection id: %s", err)
	}
	appUrl, err := getApplicationUrl("collection")
	if err != nil {
		return fmt.Errorf("failed to collection url: %s", err)
	}
	pin, err := getPinData(path)
	if err != nil {
		return fmt.Errorf("failed to load pin: %s", err)
	}
	hash, err := hashImageFile(path[:len(path)-3] + "png")
	if err != nil {
		return fmt.Errorf("failed to hash image: %s", err)
	}

	_, file := filepath.Split(path)
	file = file[:len(file)-4]
	out, err := json.MarshalIndent(Arc3Asset{
		Name:     file,
		UnitName: getUnitName(idx),

		Image:           fmt.Sprintf("ipfs://%s", pin.IpfsHash),
		Image_integrity: hash,
		Image_mimetype:  "image/png",

		Properties: SAAsset{
			Collection: SACollection{
				ID:       appId,
				Name:     viper.GetString("META_COLLECT"),
				Metadata: appUrl,
			},
		},
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	fl, err := os.Create(fmt.Sprintf("%s/%s.json", dest, file))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer fl.Close()
	fl.Write(out)

	return nil
}

func metaPushContract(path, file string) error {
	const url = "https://api.pinata.cloud/pinning/pinFileToIPFS"

	meta, err := os.Open(fmt.Sprintf("%s/%s", path, file))
	if err != nil {
		return err
	}
	defer meta.Close()

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	part, err := writer.CreateFormFile("file", filepath.Base(meta.Name()))
	if err != nil {
		return err
	}
	io.Copy(part, meta)
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
	out, err := os.Create(fmt.Sprintf("%s/%s.json.pin", path, file[:len(file)-5]))
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := out.Write(body); nil != err {
		return err
	}

	return nil
}

func metaBuildContract(path, file, src string) error {
	appId, err := getApplicationId("collection")
	if err != nil {
		return fmt.Errorf("failed to collection id: %s", err)
	}

	pin, err := getPinData(src)
	if err != nil {
		return fmt.Errorf("failed to load pin: %s", err)
	}
	hash, err := hashImageFile(src[:len(src)-3] + "png")
	if err != nil {
		return fmt.Errorf("failed to hash image: %s", err)
	}

	out, err := json.MarshalIndent(Arc3App{
		App:      appId,
		Name:     viper.GetString("META_COLLECT"),
		UnitMax:  viper.GetString("META_COLLECT_MAXCOUNT"),
		UnitName: getUnitName(0),

		Image:           fmt.Sprintf("ipfs://%s", pin.IpfsHash),
		Image_integrity: hash,
		Image_mimetype:  "image/png",
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	fl, err := os.Create(fmt.Sprintf("%s/%s", path, file))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer fl.Close()
	fl.Write(out)

	return nil
}
