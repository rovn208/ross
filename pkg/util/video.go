package util

var (
	SUPPORTED_VIDEO_EXTENSIONS = []string{"mp4", "avi"}
)

func IsSupportedExtensions(ex string) bool {
	for _, extension := range SUPPORTED_VIDEO_EXTENSIONS {
		if extension == ex {
			return true
		}
	}

	return false
}
