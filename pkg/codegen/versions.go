package codegen

// SupportedVersions describe the asyncapi-codegen supported versions.
var SupportedVersions = []string{
	"2.0.0",
	"2.1.0",
	"2.2.0",
	"2.3.0",
	"2.4.0",
	"2.5.0",
	"2.6.0",

	"3.0.0",
}

// IsVersionSupported checks that the version is supported.
func IsVersionSupported(version string) bool {
	for _, v := range SupportedVersions {
		if v == version {
			return true
		}
	}

	return false
}
