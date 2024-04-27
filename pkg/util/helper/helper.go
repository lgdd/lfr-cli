package helper

var supportedJavaVersions = [2]int{8, 11}

// Checks if the Java version is supported by Liferay
func IsSupportedJavaVersion(javaVersion int) bool {
	for _, version := range supportedJavaVersions {
		if javaVersion == version {
			return true
		}
	}
	return false
}
