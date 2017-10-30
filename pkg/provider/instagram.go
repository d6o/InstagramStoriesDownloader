package provider

import (
	"fmt"
	"github.com/ahmdrz/goinsta"
	"github.com/ahmdrz/goinsta/response"
	"time"
)

type (
	Instagram interface {
		Login() error
		Logout() error
		Stories() ([]*Media, error)
	}

	instagram struct {
		provider *goinsta.Instagram
	}

	tray response.TrayResponse

	Media struct {
		URL       string
		Username  string
		Timestamp int64
		Path      string
	}
)

const (
	Images = "images"
	Videos = "videos"
)

func NewInstagram(login, password string) (Instagram, error) {
	i := &instagram{
		provider: goinsta.New(login, password),
	}
	return i, i.Login()
}

func (i *instagram) Login() error {
	return i.provider.Login()
}

func (i *instagram) Logout() error {
	return i.provider.Logout()
}

func (i *instagram) Stories() ([]*Media, error) {
	resp, err := i.provider.GetReelsTrayFeed()
	if err != nil {
		return nil, err
	}

	tray := tray(resp)
	media := tray.Media()

	return media, nil
}

func (t *tray) Images() []*Media {
	mediaList := []*Media{}

	for _, val := range t.Tray {
		for _, me := range val.Media {
			for _, image := range me.ImageVersions2.Candidates {
				timePath := time.Unix(me.DeviceTimestamp, 0).Format("2006/01/02")

				path := fmt.Sprintf("%s/%s/%d/%d", Images, timePath, image.Width, image.Height)

				mediaList = append(mediaList, &Media{
					URL:       image.URL,
					Username:  me.User.Username,
					Timestamp: me.DeviceTimestamp,
					Path:      path,
				})
			}
		}
	}
	return mediaList
}

func (t *tray) Videos() []*Media {
	mediaList := []*Media{}

	for _, val := range t.Tray {
		for _, me := range val.Media {
			for _, video := range me.VideoVersions {
				timePath := time.Unix(me.DeviceTimestamp, 0).Format("2006/01/02")

				mediaList = append(mediaList, &Media{
					URL:       video.URL,
					Username:  me.User.Username,
					Timestamp: me.DeviceTimestamp,
					Path:      fmt.Sprintf("%s/%s/%d/%d/%d", Videos, timePath, video.Width, video.Height, video.Type),
				})
			}
		}
	}
	return mediaList
}

func (t *tray) Media() []*Media {
	mediaList := []*Media{}

	mediaList = append(mediaList, t.Images()...)
	mediaList = append(mediaList, t.Videos()...)

	return mediaList
}
