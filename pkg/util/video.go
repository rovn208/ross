package util

var (
	SupportedVideoExtensions = []string{"mp4", "mov"}
)

// IsSupportedExtensions checks if the extension is supported by walking through SupportedVideoExtensions
func IsSupportedExtensions(ex string) bool {
	for _, extension := range SupportedVideoExtensions {
		if extension == ex {
			return true
		}
	}

	return false
}
