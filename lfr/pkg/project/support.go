package project

var supportedJavaVersions = [2]int{8, 11}

func IsSupportedJavaVersion(javaVersion int) bool {
	for _, version := range supportedJavaVersions {
		if javaVersion == version {
			return true
		}
	}
	return false
}
