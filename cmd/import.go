package cmd

import "github.com/spf13/cobra"

var importCmd = &cobra.Command{
	Use:   "import runs the importer and update the values in the database",
	Short: "import Go code from other languages",
	Long:  `Import all the values from a spreadsheet into a database.`,
	Run: func(cmd *cobra.Command, args []string) {
		UpdateRows()
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}

func UpdateRows() {
	println("Updating rows")
	println("Done updating rows")
}
