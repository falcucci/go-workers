package cmd

import "github.com/spf13/cobra"

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		println("hello")
	},
}

func init() {
	RootCmd.AddCommand(helloCmd)
}
