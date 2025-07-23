package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ShowCommandHistory bool

var rootCmd = &cobra.Command{
	Use:   "gossy",
	Short: "Gossy: A versatile AWS CLI tool for efficient management across AWS services.",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&ShowCommandHistory, "history", "H", false, "Show command history in output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
