package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "copy-paste-notes",
	Short: "copy-paste-notes is a command-line note taking app",
	Long: `Manage all your notes using the command-line, including the ability to
	copy directly into your system clipboard.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	initConfig()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.copy-paste-notes.yaml)")

	if err := setupCommands(); err != nil {
		fmt.Fprintln(os.Stderr, "unable to setup commands: ", err)
		os.Exit(1)
	}
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to determine home directory:", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".copy-paste-notes")

		// Custom cofig file mapped as a volume when using Docker
		viper.AddConfigPath("/config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment Variables
	handleBindEnvErr(viper.BindEnv("db.file", "CPN_DB_FILE"))

	// Merge config
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore as we use defaults/environment variables
			// and if anything required isn't set properly (e.g. db file) we'll error later on
		} else {
			// Config file was found but another error was produced
			fmt.Fprintln(os.Stderr, "unable to merge in config: ", err)
			os.Exit(1)
		}
	}

	viper.SetDefault("db.file", path.Join(home, "cpn.db"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore as we use defaults/environment variables
			// and if anything required isn't set properly (e.g. db file) we'll error later on
		} else {
			// Config file was found but another error was produced
			fmt.Fprintln(os.Stderr, "unable to read in config: ", err)
			os.Exit(1)
		}
	}
}

func handleBindEnvErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to bind viper key to environment variable: ", err)
		os.Exit(1)
	}
}

func setupCommands() error {
	client, err := sqlite.New(viper.GetString("db.file"))
	if err != nil {
		return err
	}

	addCmd := newAddCommand(client)
	listCmd := newListCommand(client)
	copyCmd := newCopyCommand(client)
	deleteCmd := newDeleteCommand(client)

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(deleteCmd)

	return nil
}
