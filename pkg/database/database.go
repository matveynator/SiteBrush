package database

import (
	"log"
	"fmt"
	"time"
	"strings"

	"database/sql"

	"sitebrush/pkg/mylog"
	"sitebrush/pkg/data"
	"sitebrush/pkg/config"
)

// Saving post data to database
var DatabaseSavePostTask chan Data.Post

var respawnLock chan int
//по умолчанию оставляем только один процесс который будет брать задачи и записывать их в базу данных
var databaseWorkersMaxCount int = 1

func Run(config Config.Settings) {

	//initialise channel with 1000000 tasks capacity:
	DatabaseSavePostTask = make(chan Data.Post, 1000000)

	//initialize unblocking channel to guard respawn tasks
	respawnLock = make(chan int, databaseWorkersMaxCount)

	go func() {
		for {
			// will block if there is databaseWorkersMaxCount in respawnLock
			respawnLock <- 1 
			//sleep 1 second
			time.Sleep(1 * time.Second)
			go databaseWorkerRun(len(respawnLock), config)
		}
	}()
}

//close dbConnection on programm exit
func deferCleanup(db *sql.DB) {
	<-respawnLock
	err := db.Close() 
	if err != nil {
		log.Println("Error closing database connection:", err)
	}
}

func databaseWorkerRun(workerId int, config Config.Settings ) {


	dbConnection, err := connectToDb(config)

	defer deferCleanup(dbConnection)

	if err != nil  {
		MyLog.Printonce(fmt.Sprintf("Database %s is unreachable. Error: %s",  config.DB_TYPE, err))
		return

	} else {
		MyLog.Println(fmt.Sprintf("Database worker #%d connected to %s database", workerId, config.DB_TYPE))
	}

	//initialise dbConnection error channel
	databaseErrorChannel := make(chan error)

	go func() {
		for {
			//do watchdog operations only if channel with database tasks is empty (channel length equals zero):
			if len(DatabaseSavePostTask) == 0 {
				_, err = dbConnection.Exec("UPDATE DBWatchDog SET UnixTime = ? WHERE ID = 1", time.Now().UnixMilli())
				if err != nil {
					//skip busy SQLITE database errors:
					if !strings.Contains(err.Error(), "database is locked (5) (SQLITE_BUSY)") {
						log.Println("Database watchdog error:", err)
						databaseErrorChannel <- err
						return
					} else {
						//sleep some time to calm down database operations:
						log.Println("Watchdog notice: Database is busy - sleeping to calm down operations.")
						time.Sleep(config.DB_SAVE_INTERVAL_DURATION)
					}
				} else {
					//sleep some time to calm down database operations:
					time.Sleep(config.DB_SAVE_INTERVAL_DURATION)
				}
			}
		}
	}()

	// Run the main logic:
	for {
		select {

      //в случае если есть задание в канале DatabaseSavePostTask
    case <- DatabaseSavePostTask :
      //sleep some time to calm down disk operations:
      time.Sleep(config.DB_SAVE_INTERVAL_DURATION)
      //пробежать во всем доступным данным в канале заданий для бд и сохранить их в базе данных:
      for currentDatabaseSavePostTask := range DatabaseSavePostTask {
        err := SavePostDataInDB(dbConnection, currentDatabaseSavePostTask)
        if err != nil {
          //skip busy SQLITE database errors:
          if strings.Contains(err.Error(), "database is locked (5) (SQLITE_BUSY)") {
            log.Println("Saving data to disk notice: Database is busy - sleeping to calm down operations.")
            //return task to channel (this may lead to post data id misorder in database):
            DatabaseSavePostTask <- currentDatabaseSavePostTask
            //sleep some time to calm down disk operations:
            time.Sleep(config.DB_SAVE_INTERVAL_DURATION)
          }  else {
            //return task to channel (this may lead to post data id misorder in database):
            DatabaseSavePostTask <- currentDatabaseSavePostTask

            log.Printf("Database worker %d exited due to critical processing error: %s\n", workerId, err)
            return
          }
        }
      }

		case databaseError := <-databaseErrorChannel :
			//обнаружена критическая ошибка базы данных - завершаем гоурутину:
			log.Printf("Database worker %d exited due to critical error: %s\n", workerId, databaseError)
			return
		default:
			//set non blocking case
			time.Sleep(config.DB_SAVE_INTERVAL_DURATION)
		}
	}
}

