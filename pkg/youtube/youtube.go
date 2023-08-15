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
	"time"

	"github.com/rovn208/ross/pkg/util"

	"github.com/kkdai/youtube/v2"
	"github.com/rovn208/ross/pkg/configure"
	"golang.org/x/net/http/httpproxy"
)

var (
	rx             = regexp.MustCompile(`https://(.+.youtube.com|youtu.be)/(watch\?v=([^&^>^|]+)|([^&^>^|]+))`)
	ErrInvalidLink = errors.New("invalid youtube link")
)

type Client struct {
	*youtube.Client
	config configure.Config
}

type VideoYoutube struct {
	*os.File
	Video *youtube.Video
}

// NewYoutubeClient returns a new youtube client
func NewYoutubeClient(config configure.Config) *Client {
	util.Logger.Info("Creating youtube client")
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

// GetVideoID returns the video id from the url
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

// DownloadVideo downloads the video from youtube url
func (c *Client) DownloadVideo(url string) (*VideoYoutube, error) {
	util.Logger.Info("Downloading video from youtube", "url", url)
	videoID, err := c.GetVideoID(url)
	if err != nil {
		return nil, err
	}

	video, err := c.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	util.Logger.Info("Getting video from youtube", "videoID", videoID, "format", formats[0])
	fileReader, _, err := c.GetStream(video, &formats[0])
	if err != nil {
		return nil, err
	}

	util.Logger.Info("Creating youtube file", "videoID", videoID)
	dir := fmt.Sprintf("%s/%s", c.config.VideoDir, videoID)
	if err = util.CreateDirectory(dir); err != nil {
		return nil, fmt.Errorf("error when creating directory %s", dir)
	}
	file, err := os.Create(getFileName(c.config, videoID))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		return nil, err
	}

	return &VideoYoutube{
		File:  file,
		Video: video,
	}, nil
}

func getFileName(config configure.Config, videoId string) string {
	return fmt.Sprintf("%s/%s/%s.mp4", config.VideoDir, videoId, videoId)
}

// GetStreamFile returns the stream file in a format of m3u8. e.g {videoDir}/{videoId}/{videoId}.m3u8
func GetStreamFile(config configure.Config, videoId string) string {
	return fmt.Sprintf("%s/%s/%s.m3u8", config.VideoDir, videoId, videoId)
}
