package storage

import (
	"fmt"
	"github.com/disiqueira/InstagramStoriesDownloader/pkg/provider"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type (
	Downloader interface {
		Start()
		Stop()
	}

	downloader struct {
		mediaChannel chan *provider.Media
		stopSignal   bool
	}
)

func NewDownloader(mediaChannel chan *provider.Media) Downloader {
	return &downloader{
		mediaChannel: mediaChannel,
		stopSignal:   false,
	}
}

func (d *downloader) Start() {
	d.stopSignal = false

	go func() {
		for m := range d.mediaChannel {
			err := d.download(m)
			if err != nil {
				panic(err)
			}
			if d.stopSignal {
				break
			}
		}
	}()
}

func (d *downloader) Stop() {
	d.stopSignal = true
}

func (d *downloader) download(media *provider.Media) error {
	path := fmt.Sprintf("%s/%s", media.Username, media.Path)
	os.MkdirAll(path, os.ModePerm)

	u, err := url.Parse(media.URL)
	if err != nil {
		return err
	}

	urlParts := strings.Split(u.Path, "/")

	filename := urlParts[len(urlParts)-1]
	tm := time.Unix(media.Timestamp, 0)

	filePath := fmt.Sprintf("%s/%s_%s", path, tm.Format("20060102_150405"), filename)

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return nil
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	resp, err := http.Get(media.URL)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	if err = out.Close(); err != nil {
		return err
	}

	return nil
}
