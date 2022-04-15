package command

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theoremoon/ctfd-cli/ctfd"
	"golang.org/x/xerrors"
)

var (
	category string = ""
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download open tasks into your local machine",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ctfd.NewClient(serverURL, &ctfd.Credential{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Printf("%v\n", err)
			os.Exit(1)
		}

		chals, err := client.ListChallenges()
		if err != nil {
			log.Printf("%v\n", err)
			os.Exit(1)
		}

		cat := strings.ToLower(category)
		for _, c := range chals {
			if cat != "" {
				// filter by category
				if strings.Index(strings.ToLower(c.Category), cat) == -1 {
					continue
				}
			}

			// download challenge
			log.Printf("downloading challenge %s\n", c.Name)
			cdata, err := client.GetChallenge(c.ID)
			if err != nil {
				log.Printf("%v\n", err)
				continue
			}

			// mkdir
			dir := cleanPath(c.Name)
			if err := os.Mkdir(dir, 0766); err != nil && !os.IsExist(err) {
				// ここに到達したときはpermissionとかdisk fullとかで次も失敗するだろうから終わっとくか
				log.Printf("%v\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(filepath.Join(dir, "description.md"), []byte(formatDescription(cdata)), 0755); err != nil {
				log.Printf("%v\n", err)
				os.Exit(1)
			}

			// download attachments
			// このあたり並列化すると多少速度稼げそうだけどまあいいか
			for _, f := range cdata.Files {
				data, name, err := downloadAttachment(serverURL, f)
				if err != nil {
					log.Printf("failed to download attachment: %v\n", err)
					continue
				}

				if err := os.WriteFile(filepath.Join(dir, name), data, 0755); err != nil {
					log.Printf("%v\n", err)
					os.Exit(1)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	downloadCmd.Flags().StringVar(&category, "category", "", "Specify category to download")
}

/// attachmentをダウンロードして、ファイル名とデータ(bytes)を返す
func downloadAttachment(base string, subpath string) ([]byte, string, error) {
	p := base + subpath
	if !strings.HasSuffix(p, "/") {
		p = base + "/" + subpath
	}
	u, err := url.Parse(p)
	if err != nil {
		return nil, "", xerrors.Errorf(": %w", err)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, "", xerrors.Errorf(": %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", xerrors.Errorf(": %w", err)
	}

	parts := strings.Split(u.Path, "/")
	return data, parts[len(parts)-1], nil
}

func formatDescription(chal *ctfd.Challenge) string {
	// TODO: HINTとかも入れといたほうがよくない？ wowow
	return chal.Description
}

func cleanPath(p string) string {
	return strings.ReplaceAll(filepath.Clean(p), "/", "_")
}
