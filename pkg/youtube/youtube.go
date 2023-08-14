package youtube

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/rovn208/ross/pkg/configure"
	"golang.org/x/net/http/httpproxy"
)

var (
	rx              = regexp.MustCompile(`https://(.+.youtube.com|youtu.be)/(watch\?v=([^&^>^|]+)|([^&^>^|]+))`)
	ErrInvalidLink  = errors.New("invalid youtube link")
	ErrCreatingFile = errors.New("error when creating file")
	ErrRemovingFile = errors.New("error when removing file")
)

type Client struct {
	*youtube.Client
	config configure.Config
}

func NewYoutubeClient(config configure.Config) *Client {
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
		config: config,
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

func (c *Client) DownloadVideo(url string) (*os.File, error) {
	videoID, err := c.GetVideoID(url)
	if err != nil {
		return nil, err
	}
	video, err := c.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	fileReader, _, err := c.GetStream(video, &formats[0])
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.mp4", c.config.VideoDir, videoID))
	if err != nil {
		return nil, ErrCreatingFile
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getAudioWebmFormat(v *youtube.Video) (*youtube.Format, error) {
	formats := v.Formats

	audioFormats := formats.Type("audio")
	audioFormats.Sort()
	for _, fm := range formats {
		if strings.HasPrefix(fm.MimeType, "audio/webm") {
			return &fm, nil
		}
	}
	// no webm, take first format
	return &formats[0], nil
}
