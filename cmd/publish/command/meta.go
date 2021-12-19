package command

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Meta.AddCommand(metaBuildCmd)
}

// Arc3App ASA Metadata Info (SA042-ASA)
type Arc3App struct {
	App      string `json:"app"`
	Name     string `json:"name"`
	AsaMax   string `json:"asamax"`
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
	AsaMax   string `json:"asamax"`
	AsaMask  string `json:"asamask"`
	Metadata string `json:"metadata"`
}

var Meta = &cobra.Command{
	Use:   "meta",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var metaBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := fmt.Sprintf("%s/images", viper.GetString("DATA"))
		dest := fmt.Sprintf("%s/metadata", viper.GetString("DATA"))
		fmt.Println(":: Building meta data for:", root)

		list, err := getListOfFiles(".pin", root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(1)
		}

		fmt.Printf(">> Found %d assets to build to metadata for.\n", len(list))
		for idx, path := range list {
			fmt.Printf(">> Build > %s ", path)
			if err := metaBuildfData(dest, path, uint32(idx+1)); nil != err {
				fmt.Printf(">> Error: \n>>>> %s\n", err)
			} else {
				fmt.Println(">> Ok")
			}
		}
	},
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

func metaBuildfData(dest, path string, idx uint32) error {
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
	out, err := json.MarshalIndent(Arc3Asset{
		Name:     file[:len(file)-5],
		UnitName: getUnitName(idx),

		Image:           fmt.Sprintf("ipfs://%s", pin.IpfsHash),
		Image_integrity: hash,
		Image_mimetype:  "image/png",

		Properties: SAAsset{
			// Traits: SATraits{
			// 	"background": "purple",
			// },
			// Attributes: SAAttributes{
			// },
			Collection: SACollection{
				ID:       appId,
				Name:     viper.GetString("META_COLLECT"),
				AsaMax:   viper.GetString("META_COLLECT_MAXCOUNT"),
				AsaMask:  getUnitName(0),
				Metadata: appUrl,
			},
		},
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	fl, err := os.Create(fmt.Sprintf("%s/%s", dest, file))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer fl.Close()
	fl.Write(out)

	return nil
}

func metaBuildfContract(dest, path string) error {
	appId, err := getApplicationId("collection")
	if err != nil {
		return fmt.Errorf("failed to collection id: %s", err)
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
	out, err := json.MarshalIndent(Arc3App{
		App:      appId,
		Name:     viper.GetString("META_COLLECT"),
		AsaMax:   viper.GetString("META_COLLECT_MAXCOUNT"),
		UnitName: getUnitName(0),

		Image:           fmt.Sprintf("ipfs://%s", pin.IpfsHash),
		Image_integrity: hash,
		Image_mimetype:  "image/png",
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	file = file[:len(file)-4] + ".json"
	fl, err := os.Create(fmt.Sprintf("%s/%s", dest, file))
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}
	defer fl.Close()
	fl.Write(out)

	return nil
}
