package internal

type Resolution string

const (
	FileResolution    Resolution = "file"
	PackageResolution Resolution = "package"
)

var validResolutions = map[Resolution]bool{
	FileResolution:    true,
	PackageResolution: true,
}

func IsValidResolution(resolution Resolution) bool {
	_, ok := validResolutions[resolution]
	return ok
}

func ValidResolutions() []string {
	resolutions := make([]string, 0, len(validResolutions))
	for k, _ := range validResolutions {
		resolutions = append(resolutions, string(k))
	}
	return resolutions
}
