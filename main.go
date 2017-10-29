package main

import (
	"github.com/disiqueira/InstagramStoriesDownloader/pkg/config"
	"github.com/disiqueira/InstagramStoriesDownloader/pkg/provider"
	"github.com/disiqueira/InstagramStoriesDownloader/pkg/storage"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	configs := config.NewSpecification()
	if err := envconfig.Process("isd", configs); err != nil {
		panic(err.Error())
	}

	instagram, err := provider.NewInstagram(configs.Username(), configs.Password())
	if err != nil {
		panic(err.Error())
	}

	mediaChannel := make(chan *provider.Media, 100000)

	startWorkers(configs.Workers(), mediaChannel)

	for {
		stories, err := instagram.Stories()
		if err != nil {
			panic(err.Error())
		}

		for _, m := range stories {
			mediaChannel <- m
		}
	}
}

func startWorkers(numWorkers int, mediaChannel chan *provider.Media) {
	workers := make([]storage.Downloader, numWorkers)

	for i := 0; i <= numWorkers; i++ {
		workers[i] = storage.NewDownloader(mediaChannel)
		workers[i].Start()
	}
}
