package appstate

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_ENV_PATH = "/home/bao/selfproj/kanban_drag_drop/be/env.yaml"
)

var DEFAULT_QUERY_TIMEOUT int64 = 20

type Appconfig struct {
	Dburl        string `yaml:"dburl"`
	QueryTimeout *int64 `yaml:"query_timeout"`
	ServerUrl    string `yaml:"server_url"`
}

type AppState struct {
	Db     *sql.DB
	Config *Appconfig
}

var appState *AppState

func init() {
	configPath := os.Getenv("KANBAN_CONFIG")
	if _, err := os.Stat(configPath); err != nil {
		if _, err = os.Stat(DEFAULT_ENV_PATH); err != nil {
			log.Fatal("config file not exists")
		}
		configPath = DEFAULT_ENV_PATH
	}
	var config Appconfig
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Failed to parse yaml config: ", err)
	}
	db, err := sql.Open("mysql", config.Dburl)
	if err != nil {
		log.Fatal("Failed to create db instance: ", err)
	}
	if config.QueryTimeout == nil {
		config.QueryTimeout = &DEFAULT_QUERY_TIMEOUT
	}
	appState = &AppState{
		Db:     db,
		Config: &config,
	}
}

func GetAppState() *AppState {
	return appState
}
