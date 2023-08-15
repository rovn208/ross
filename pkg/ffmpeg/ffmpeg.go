package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rovn208/ross/pkg/util"
)

// CheckVideoCorrupted checks whether the video is corrupted
func CheckVideoCorrupted(fileName string) error {
	cmd := exec.Command("ffmpeg", "-v", "error", "-i", fileName, "-f", "null", "-")
	err := cmd.Run()
	if err != nil {
		util.Logger.Error("The video is corrupted", "error", err)
		return err
	}
	return nil

}

// ToHLSFormat converts the video to hls
func ToHLSFormat(ctx context.Context, fileName string) error {
	err := CheckVideoCorrupted(fileName)
	if err != nil {
		return err
	}
	// ffmpeg -i videos/n1i6bnqzXM0.mp4 -c copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls videos/test.m3u8
	util.Logger.Info(fmt.Sprintf("Converting file %s to hls", fileName))
	fileNameWithoutExtension := strings.Split(fileName, ".")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", fileName, "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", fmt.Sprintf("%s.m3u8", fileNameWithoutExtension[0]))
	return cmd.Run()
}
