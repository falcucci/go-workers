package cmd

import (
	"fmt"
	"go-workers/database"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

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
	var DB = database.DB
	transaction := DB.Begin()

	temp := "temp"
	table := "people"
	execution := ShouldCreateTempTable(transaction, temp, table)
	if execution.Error != nil {
		transaction.Rollback()
		panic(execution.Error)
	}

	fmt.Printf("Creating temporary %s table\n", temp)

	println("Done updating rows")
}

// Method to create a non-existent temporary table with the default schema from an
// existent table
func ShouldCreateTempTable(db *gorm.DB, table string, from string) *gorm.DB {
	createTempTable := fmt.Sprintf(
		`CREATE TEMP TABLE IF NOT EXISTS %s AS SELECT * FROM %s LIMIT 0`,
		table,
		from,
	)
	return db.Exec(createTempTable)
}
