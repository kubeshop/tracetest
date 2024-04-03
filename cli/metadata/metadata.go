package metadata

import cienvironment "github.com/cucumber/ci-environment/go"

var (
	tracetestSource     = "tracetest.source"
	tracetestCliVersion = "tracetest.cli.version"
	gitRemote           = "git.GitRemote"
	gitBranch           = "git.branch"
	gitTag              = "git.tag"
	gitSha              = "git.sha"
	cIBuildNumber       = "ci.build.number"
	cIProvider          = "ci.provider"
	cIBuildUrl          = "ci.build.url"
)

type Metadata map[string]string

func (m Metadata) Merge(other Metadata) Metadata {
	for k, v := range other {
		m[k] = v
	}

	return m
}

func GetMetadata() Metadata {
	// TODO: add more metadata after getting the response from the k6 team
	// https://github.com/grafana/k6/issues/1320#issuecomment-2032734378
	metadata := Metadata{}
	metadata[tracetestSource] = "cli"
	metadata[tracetestCliVersion] = "0.1.8"

	ci := cienvironment.DetectCIEnvironment()
	if ci == nil {
		return metadata
	}

	metadata[cIProvider] = ci.Name
	metadata[cIBuildUrl] = ci.URL
	metadata[cIBuildNumber] = ci.BuildNumber
	metadata[gitBranch] = ci.Git.Branch

	if ci.Git != nil {
		metadata[gitTag] = ci.Git.Tag
		metadata[gitSha] = ci.Git.Revision
		metadata[gitRemote] = ci.Git.Remote
	}

	return metadata
}
