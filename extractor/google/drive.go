// Copyright 2023 @falcucci
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extractor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"go-workers/config"
)

var (
	Env = config.Env
)

// Method to scrap every value from people spreadsheet
func GetPeopleSheetValues() *sheets.ValueRange {
	client := GetGoogleDriveClient()
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	print("Getting people sheet values...")
	readRangePeople := fmt.Sprintf("!A1:B")
	people, err := srv.Spreadsheets.Values.Get(
		Env.SpreadSheetId,
		readRangePeople,
	).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve sheet: %v", err)
	}
	return people
}

// Method that returns drive client based on the configs
func GetGoogleDriveClient() *http.Client {
	credentials, err := ioutil.ReadFile("service.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	driveConfig, err := google.JWTConfigFromJSON(
		credentials,
		sheets.DriveScope,
		sheets.DriveFileScope,
		sheets.SpreadsheetsScope,
		sheets.SpreadsheetsReadonlyScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return driveConfig.Client(oauth2.NoContext)
}
