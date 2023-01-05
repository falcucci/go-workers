package cmd

import (
	"fmt"
	"go-workers/database"
	extractor "go-workers/extractor/google"
	"go-workers/structs"
	"log"
	"math/rand"
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

	temp := "temp_people"
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
	insert_temp := shouldInsertPeopleTempValues(transaction, temp)
	if insert_temp.Error != nil {
		transaction.Rollback()
		panic(insert_temp.Error)
	}
	fmt.Printf("Inserted %d records in the temporary table\n\n", insert_temp.RowsAffected)
	println("Done inserting values into temporary table")

	println("Inserting values into people table")
	insert := shouldInsertPeopleValues(transaction)
	if insert.Error != nil {
		transaction.Rollback()
		panic(insert.Error)
	}

	fmt.Printf("Inserted %d records in the people table\n\n", insert.RowsAffected)
	println("Done inserting values into people table")

	fmt.Println(`Checking new seller values to update`)
	update := shouldUpdatePeopleValues(transaction)
	if update.Error != nil {
		// rollback the transaction in case of error
		transaction.Rollback()
		log.Fatal(update.Error)
	}
	fmt.Printf(
		"Updated %d records in the people table\n\n",
		update.RowsAffected,
	)

	fmt.Println(`Checking possible people to remove`)
	delete := shouldDeletePeopleRow(transaction)
	if delete.Error != nil {
		// rollback the transaction in case of error
		transaction.Rollback()
		log.Fatal(delete.Error)
	}
	fmt.Printf(
		"Removing %d seller's values in the people table\n\n",
		delete.RowsAffected,
	)

	transaction.Commit()
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
		`INSERT INTO %s (id, name, surname) VALUES %s`
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
		// generate random number in golang to use as id
		// format float64 to int to use as id
		person.Id = int(GenerateRandomNumber())
		person.Name = line[0].(string)
		person.Surname = line[1].(string)

		s = append(s, fmt.Sprintf(
			"(%d, '%s', '%s')",
			person.Id,
			person.Name,
			person.Surname,
		))
	}
	v := strconv.Quote(strings.Join(s, ", "))
	v = v[1 : len(v)-1]
	return v
}

// This Method will generate a random int value
// with no specific upper or lower limits.
func GenerateRandomNumber() int {
	// Seed the generator with the current time
	return rand.Intn(100000)
}

// Method to compare values with the temporary
// table and insert people if the people exists
// in the sheet and not exists in the database.
func shouldInsertPeopleValues(db *gorm.DB) *gorm.DB {
	insertQuery := `
		INSERT INTO people
		(
			id,
			name,
			surname
		) (
		SELECT
			tp.id,
			tp.name,
			tp.surname
		FROM
			temp_people tp
		LEFT JOIN
			people p
		ON (tp.id = p.id)
		WHERE
			p.id is null)`
	return db.Exec(insertQuery)
}

// Method to compare values with the temporary
// table and update people if the people exists
// in the sheet and exists in the database.
func shouldUpdatePeopleValues(db *gorm.DB) *gorm.DB {
	updateQuery := `
		UPDATE people p
		SET
			id = tp.id,
			name = tp.name,
			surname = tp.surname
		FROM
			temp_people tp
		WHERE
			p.id = tp.id AND ( p.name != tp.name OR p.surname != tp.surname )`
	return db.Exec(updateQuery)
}

// Method to compare values with the temporary
// table and delete people if the people exists
// in the database and not exists in the sheet.
func shouldDeletePeopleRow(db *gorm.DB) *gorm.DB {
	deleteQuery := `
		DELETE FROM people p
		WHERE NOT EXISTS (
			SELECT 1
			FROM temp_people tp
			WHERE (
				p.id = tp.id 
				AND p.name = tp.name AND p.surname = tp.surname
			)
		)`
	return db.Exec(deleteQuery)
}
