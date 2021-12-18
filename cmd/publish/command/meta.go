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

// Arc3 NFT Metadata (ASA Spec)
type Arc3 struct {
	Name        string `json:"name"`
	UnitName    string `json:"unitname"`
	Description string `json:"description"`

	Image           string `json:"image"`
	Image_integrity string `json:"image_integrity"`
	Image_mimetype  string `json:"image_mimetype"`

	Properties SAProps `json:"properties"`
}

// SAProps NFT Properies (SA007)
type SAProps struct {
	// Traits     SATraits  `json:"traits"`
	// Attributes SAAttributes `json:"attributes"`
	Collection SACollection `json:"collection"`
}

// SATraits NFT Traits (SA069)
// type SATraits map[string]string

// SAAttributes RPG Attributes (SA076)
// type SAAttributes struct {}

// SACollection NFT Collections (SA0042)
type SACollection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Masks    string `json:"masks"`
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

		list, err := getListOfFiles(".json", root)
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
	appId, err := getapplicationId("collection")
	if err != nil {
		return fmt.Errorf("failed to collection id: %s", err)
	}
	appUrl, err := getapplicationUrl("collection")
	if err != nil {
		return fmt.Errorf("failed to collection url: %s", err)
	}
	pin, err := getPinData(path)
	if err != nil {
		return fmt.Errorf("failed to load pin: %s", err)
	}
	hash, err := hashImageFile(path[:len(path)-4] + "png")
	if err != nil {
		return fmt.Errorf("failed to hash image: %s", err)
	}

	_, file := filepath.Split(path)
	out, err := json.MarshalIndent(Arc3{
		Name:     file[:len(file)-5],
		UnitName: getUnitName(idx),

		Image:           fmt.Sprintf("ipfs://%s", pin.IpfsHash),
		Image_integrity: hash,
		Image_mimetype:  "image/png",

		Properties: SAProps{
			// Traits: SATraits{
			// 	"background": "purple",
			// },
			// Attributes: SAAttributes{
			// },
			Collection: SACollection{
				ID:       appId,
				Name:     viper.GetString("META_COLLECT"),
				Masks:    getUnitName(0),
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
