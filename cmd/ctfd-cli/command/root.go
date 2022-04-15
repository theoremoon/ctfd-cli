package command

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctfd-cli",
	Short: "A command line CTFd client utlitity",
	Long: `A command line CTFd client utlitity.
You can download tasks and get task statuses.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// persistent variable
var (
	username  string
	password  string
	serverURL string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctfd-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	rootCmd.MarkPersistentFlagRequired("username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
	rootCmd.MarkPersistentFlagRequired("password")
	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "", "CTFd Server URL")
	rootCmd.MarkPersistentFlagRequired("url")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
