package ffmpeg

import (
	"context"
	"fmt"
	"github.com/rovn208/ross/pkg/util"
	"os/exec"
	"strings"
)

func ToHLSFormat(ctx context.Context, fileName string) error {
	// ffmpeg -i videos/n1i6bnqzXM0.mp4 -c copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls videos/test.m3u8
	util.Logger.Info(fmt.Sprintf("Converting file %s to hls", fileName))
	fileNameWithoutExtension := strings.Split(fileName, ".")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", fileName, "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", fmt.Sprintf("%s.m3u8", fileNameWithoutExtension[0]))
	return cmd.Run()
}
