package jobs

import (
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/skatekrak/scribe/clients/feedly"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/services"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) {
	s := gocron.NewScheduler(time.UTC)

	if db == nil {
		log.Fatalln("Cannot start jobs, missing DB")
	}

	// At midnight every day
	if _, err := s.Cron("0 0 * * *").Do(refreshFeedly(db)); err != nil {
		log.Fatalf("Cannot start refreshFeedly job: %s", err.Error())
	}
	if _, err := s.Cron("0 0 * * *").Do(refreshVideos(db)); err != nil {
		log.Fatalf("Cannot start refreshVideos job: %s", err.Error())
	}

	s.StartAsync()
	log.Println("scheduler started")
}

func refreshFeedly(db *gorm.DB) func() {
	return func() {
		feedlyCategoryID := os.Getenv("FEEDLY_FETCH_CATEGORY_ID")

		feedlyClient := feedly.New(os.Getenv("FEEDLY_API_KEY"))
		fetcher := fetchers.New(nil, nil, feedlyClient)

		refreshService := services.NewRefreshService(db, fetcher, feedlyCategoryID)

		if _, err := refreshService.RefreshFeedlySource(); err != nil {
			log.Printf("Error refreshing feedly sources: %s", err.Error())
		} else {
			log.Println("Feedly sources refreshed")
		}

		if _, err := refreshService.RefreshByTypes([]string{"rss"}); err != nil {
			log.Printf("Error fetching feedly contents: %s", err.Error())
		} else {
			log.Println("Feedly contents refreshed")
		}
	}
}

func refreshVideos(db *gorm.DB) func() {
	return func() {
		youtubeClient := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
		vimeoClient := vimeo.New(os.Getenv("VIMEO_API_KEY"))

		fetcher := fetchers.New(vimeoClient, youtubeClient, nil)

		refreshService := services.NewRefreshService(db, fetcher, "")

		if _, err := refreshService.RefreshByTypes([]string{"vimeo", "youtube"}); err != nil {
			log.Printf("Error refreshing videos: %s", err.Error())
		} else {
			log.Println("Videos refreshed")
		}
	}
}
