package cmd

import (
	"fmt"
	"go-workers/database"
	extractor "go-workers/extractor/google"
	"go-workers/structs"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
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

	println("Truncating temporary table")
	truncate := TruncateTempTable(transaction, temp)
	if truncate.Error != nil {
		transaction.Rollback()
		panic(truncate.Error)
	}
	println("Done truncating temporary table")

	println("Inserting values into temporary table")
	insert := shouldInsertPeopleTempValues(transaction, temp)
	if insert.Error != nil {
		transaction.Rollback()
		panic(insert.Error)
	}
	println("Done inserting values into temporary table")

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

// Method to clear all the content of the temporary table
func TruncateTempTable(db *gorm.DB, table string) *gorm.DB {
	return db.Exec(fmt.Sprintf("TRUNCATE %s", table))
}

// Scrap function the read values from the spreadsheet
// and insert all the values into the temporary table. These
// values must be the same of the sheet
func shouldInsertPeopleTempValues(db *gorm.DB, temp string) *gorm.DB {
	query := getInsertPeopleQuery(temp)
	return db.Exec(query)
}

// method to build the query correctly with the values to insert
// according with the spreadsheet in the google drive.
func getInsertPeopleQuery(tempTable string) string {
	pplSheetValues := extractor.GetPeopleSheetValues()
	people := formatPeopleValues(pplSheetValues)
	insertQuery :=
		`INSERT INTO %s (id, description)
		 VALUES %s`
	query := fmt.Sprintf(insertQuery, tempTable, people)
	return query
}

// Method to format the people values scraped from the spreadsheet
// and put it in the insert query formated
func formatPeopleValues(
	people *sheets.ValueRange,
) string {
	var s = []string{}
	for _, line := range people.Values {
		person := structs.Person{}
		person.Name = line[0].(string)
		person.Surname = line[1].(string)

		s = append(s, fmt.Sprintf(
			"('%s', '%s')",
			person.Name,
			person.Surname,
		))
	}
	v := strconv.Quote(strings.Join(s, ", "))
	v = v[1 : len(v)-1]
	return v
}
