package master

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func RunScheduler() {
	ticker := time.NewTicker(time.Second * 1)
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			err := doSchedule(db)
			if err != nil {
				log.Printf("schedule error: %v", err)
			}
		}
	}
}

func doSchedule(db *gorm.DB) error {

	return nil
}
