// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"fmt"
	_ "log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationName    string
	Env                string
	DatabaseName       string
	UserName           string
	Host               string
	Password           string
	Port               string
	SpreadSheetId      string
	RefreshToken       string
	Kind               string
	Team               string
	LogrusLogLevel     string
	BurzumLogLevel     string
	BurzumToken        string
	Schedule           string
	ShowSql            bool
	MaxConcurrentProcs int
}

var (
	Env = GetConfig()
)

// Method to execute and read the environment variables
func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Runnning the application without a .env file.")
	}
}

// Method to get available variables in the environment
func GetConfig() Config {
	var config Config
	LoadEnvs()
	config.DatabaseName = os.Getenv("DATABASE_NAME")
	config.UserName = os.Getenv("DATABASE_USERNAME")
	config.Host = os.Getenv("DATABASE_HOST")
	config.Password = os.Getenv("DATABASE_PASSWORD")
	config.Port = os.Getenv("DATABASE_PORT")
	config.ApplicationName = os.Getenv("APPLICATION")
	config.SpreadSheetId = os.Getenv("SPREAD_SHEET_ID")
	config.RefreshToken = os.Getenv("REFRESH_TOKEN")
	config.Env = os.Getenv("ENV")
	config.Kind = os.Getenv("KIND")
	config.Team = os.Getenv("TEAM")
	config.LogrusLogLevel = os.Getenv("LOGRUS_LOG_LEVEL")
	config.BurzumLogLevel = os.Getenv("BURZUM_LOG_LEVEL")
	config.BurzumToken = os.Getenv("BURZUM_TOKEN")
	config.Schedule = os.Getenv("SCHEDULE")
	config.ShowSql = GetBoolFromEnv("SHOW_SQL", false)
	config.MaxConcurrentProcs = GetIntFromEnv(
		"MAX_CONCURRENT_PROCS",
	)
	return config
}

// get a bool value from a environment variable,
// if an error or the environment variable doesn't exists,
// the default value will be always returned.
func GetBoolFromEnv(key string, defaultValue bool) bool {
	env := os.Getenv(key)
	if env != "" {
		boolEnv, err := strconv.ParseBool(env)
		if err != nil {
			return defaultValue
		}
		return boolEnv
	}
	return defaultValue
}

// get an integer value from an environment variable,
// if an error or the environment variable doesn't exists,
// the default value will be always returned.
func GetInt64FromEnv(key string) int64 {
	env := os.Getenv(key)
	intEnv, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		// TODO: it doesn't make sense to return 0
		// move it out of here later
		return 0
	}
	return int64(intEnv)
}

// get an integer value from an environment variable,
// if an error or the environment variable doesn't exists,
// the default value will be always returned.
func GetIntFromEnv(key string) int {
	env := os.Getenv(key)
	intEnv, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		// TODO: it doesn't make sense to return 0
		// move it out of here later
		return 0
	}
	return int(intEnv)
}
