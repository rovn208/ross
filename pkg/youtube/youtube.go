package youtube

import (
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"golang.org/x/net/http/httpproxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	rx             = regexp.MustCompile(`https://(.+.youtube.com|youtu.be)/(watch\?v=([^&^>^|]+)|([^&^>^|]+))`)
	ErrInvalidLink = errors.New("invalid youtube link")
)

type Client struct {
	*youtube.Client
}

func NewYoutubeClient() *Client {
	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		// Proxy: http.ProxyFromEnvironment() does not work. Why?
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Client{
		Client: &youtube.Client{
			HTTPClient: &http.Client{Transport: httpTransport},
		},
	}
}

func (c *Client) GetVideoID(url string) (id string, err error) {
	if rx.MatchString(url) {
		sub := rx.FindStringSubmatch(url)
		if len(sub) > 3 {
			// whenever link is youtu.be return sub 4
			if sub[1] == "youtu.be" && len(sub[4]) > 0 {
				return sub[4], nil
			}

			// Assume link prefix is youtube.com, take the id from ?v=... and check the len
			if len(sub[3]) > 0 {
				return sub[3], nil
			}
		}

	}
	return "", fmt.Errorf("%s %w", url, ErrInvalidLink)
}

func (c *Client) DownloadVideo(url string) error {
	videoID, err := c.GetVideoID(url)
	if err != nil {
		return err
	}
	video, err := c.GetVideo(videoID)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := c.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	videoDir := "videos" // TODO: Using env param
	file, err := os.Create(fmt.Sprintf("%s/%s.mp4", videoDir, videoID))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}
	return nil
}
