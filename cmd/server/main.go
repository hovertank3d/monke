package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hovertank3d/monke"
	"github.com/hovertank3d/monke/env"
	"github.com/hovertank3d/monke/parser"
	"github.com/hovertank3d/monke/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func scan(db *gorm.DB, pages int) {
	p := parser.New()
	db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&monke.Season{}).Where("1 = 1").Delete(&monke.Season{})

		allAnimes := map[string]struct{}{}
		for i := 0; i < pages; i++ {
			println("page ", i)
			animes, err := p.Page(i)
			if err != nil || len(animes) == 0 {
				break
			}

			for _, a := range animes {
				allAnimes[a] = struct{}{}
			}
		}

		for link := range allAnimes {
			var anime monke.Anime

			anime, err := p.Anime(link)
			if err != nil {
				log.Fatal(err)
			}

			tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "original_name", "tags", "description"}),
			}).Create(&anime)
		}
		return nil
	})
}

func main() {
	db, err := gorm.Open(postgres.Open(env.DSL))
	if err != nil {
		log.Fatal(err)
	}

	host := ":8080"

	db.AutoMigrate(&monke.Anime{})
	db.AutoMigrate(&monke.Season{})

	if len(os.Args) > 2 && os.Args[1] == "rescan" {
		pages, _ := strconv.Atoi(os.Args[2])
		scan(db, pages)
		return
	} else if len(os.Args) == 2 {
		host = os.Args[1]
	}

	s := server.New(&monke.DB{DB: db})
	s.Logger.Fatal(s.Start(host))
}
