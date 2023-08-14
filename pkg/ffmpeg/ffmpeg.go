package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ToHLSFormat(ctx context.Context, file *os.File) error {
	// ffmpeg -i videos/n1i6bnqzXM0.mp4 -c copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls videos/test.m3u8
	arr := strings.Split(file.Name(), ".")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", file.Name(), "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", fmt.Sprintf("%s.m3u8", arr[0]))
	return cmd.Run()
}
