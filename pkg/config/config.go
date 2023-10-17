package Config

import (
	"os"
	"fmt"
	"flag"
	"time"
	"hash/fnv"

)
var CompileVersion string

type Settings struct {
	APP_NAME, VERSION, DAEMON_LISTENER_ADDRESS, DAEMON_LISTENER_ADDRESS_HASH, WEB_LISTENER_ADDRESS, PROXY_ADDRESS, DB_TYPE, DB_FILE_PATH, DB_FULL_FILE_PATH, PG_HOST, PG_USER, PG_PASS, PG_DB_NAME, PG_SSL, TIME_ZONE, RACE_TYPE string
	PG_PORT int
	AVERAGE_RESULTS, VARIABLE_DISTANCE_RACE bool
	DB_SAVE_INTERVAL_DURATION time.Duration
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return fmt.Sprint(hash.Sum32())
}

func ParseFlags() (config Settings)  { 
	config.APP_NAME = "sitebrush"
	flagVersion := flag.Bool("version", false, "Output version information")


	flag.StringVar(&config.WEB_LISTENER_ADDRESS, "web", "0.0.0.0:80", "Please specify the IP address and port number on which HTTP web interface will be running.")
	flag.StringVar(&config.TIME_ZONE, "timezone", "UTC", "Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")

	//db
	flag.StringVar(&config.DB_FILE_PATH, "db-path", ".", "Provide path to writable directory to store database data.")
	flag.StringVar(&config.DB_TYPE, "db-type", "genji", "Select db type: sqlite / genji / postgres")
	flag.DurationVar(&config.DB_SAVE_INTERVAL_DURATION, "db-save-interval", 30000*time.Millisecond, "Duration to save data from memory to database (disk). Setting duration too low may cause unpredictable performance results." )

	//PostgreSQL related start
	flag.StringVar(&config.PG_HOST, "pg-host", "127.0.0.1", "PostgreSQL DB host.")
	flag.IntVar(&config.PG_PORT, "pg-port", 5432, "PostgreSQL DB port.")
	flag.StringVar(&config.PG_USER, "pg-user", "postgres", "PostgreSQL DB user.")
	flag.StringVar(&config.PG_PASS, "pg-pass", "", "PostgreSQL DB password.")
	flag.StringVar(&config.PG_DB_NAME, "pg-db-name", "chicha", "PostgreSQL DB name.")
	flag.StringVar(&config.PG_SSL, "pg-ssl", "prefer", "disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html")

	//process all flags
	flag.Parse()

	//делаем хеш от порта коллектора чтобы использовать в уникальном названии файла бд
	config.DAEMON_LISTENER_ADDRESS_HASH = hash(config.DAEMON_LISTENER_ADDRESS)

	//путь к файлу бд
	config.DB_FULL_FILE_PATH = fmt.Sprintf(config.DB_FILE_PATH+"/"+config.APP_NAME+"."+config.DAEMON_LISTENER_ADDRESS_HASH+".db."+config.DB_TYPE)

	//set version from CompileVersion variable at build time
	config.VERSION = CompileVersion 

	if *flagVersion  {
		if config.VERSION != "" {
			fmt.Println("Version:", config.VERSION)
		}
		os.Exit(0)
	}

	// Startup banner START:
	fmt.Printf("Starting %s ", config.APP_NAME)
	if config.VERSION != "" {
		fmt.Printf("version %s ", config.VERSION)
	}
	fmt.Printf("at %s \n", config.WEB_LISTENER_ADDRESS)


	// END.


	return
}
