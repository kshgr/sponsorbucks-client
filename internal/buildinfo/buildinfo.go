package buildinfo

const DefaultVersion = "1.0.0-preview"
const DefaultBuildID = "dev"
const DefaultBuildChannel = "dev"

var Version = DefaultVersion
var BuildID = DefaultBuildID
var BuildChannel = DefaultBuildChannel

type Info struct {
	ClientVersion string
	BuildID       string
	BuildChannel  string
}

func Current() Info {
	return Info{
		ClientVersion: Version,
		BuildID:       BuildID,
		BuildChannel:  BuildChannel,
	}
}

func IsDev() bool {
	return BuildChannel == DefaultBuildChannel || BuildID == DefaultBuildID
}
