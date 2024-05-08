package stereoscope

import (
	"github.com/anchore/go-collections"
	containerdClient "github.com/anchore/stereoscope/internal/containerd"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/stereoscope/pkg/image/containerd"
	"github.com/anchore/stereoscope/pkg/image/docker"
	"github.com/anchore/stereoscope/pkg/image/oci"
	"github.com/anchore/stereoscope/pkg/image/podman"
	"github.com/anchore/stereoscope/pkg/image/sif"
	"github.com/anchore/stereoscope/pkg/pathfilter"
)

const (
	FileTag     = "file"
	DirTag      = "dir"
	DaemonTag   = "daemon"
	PullTag     = "pull"
	RegistryTag = "registry"
)

// ImageProviderConfig is the uber-configuration containing all configuration needed by stereoscope image providers
type ImageProviderConfig struct {
	UserInput      string
	Platform       *image.Platform
	Registry       image.RegistryOptions
	PathFilterFunc pathfilter.PathFilterFunc
}

func ImageProviders(cfg ImageProviderConfig) []collections.TaggedValue[image.Provider] {
	tempDirGenerator := rootTempDirGenerator.NewGenerator()
	return []collections.TaggedValue[image.Provider]{
		// file providers
		taggedProvider(docker.NewArchiveProvider(tempDirGenerator, cfg.UserInput, cfg.PathFilterFunc), FileTag),
		taggedProvider(oci.NewArchiveProvider(tempDirGenerator, cfg.UserInput, cfg.PathFilterFunc), FileTag),
		taggedProvider(oci.NewDirectoryProvider(tempDirGenerator, cfg.UserInput, cfg.PathFilterFunc), FileTag, DirTag),
		taggedProvider(sif.NewArchiveProvider(tempDirGenerator, cfg.UserInput, cfg.PathFilterFunc), FileTag),

		// daemon providers
		taggedProvider(docker.NewDaemonProvider(tempDirGenerator, cfg.UserInput, cfg.Platform, cfg.PathFilterFunc), DaemonTag, PullTag),
		taggedProvider(podman.NewDaemonProvider(tempDirGenerator, cfg.UserInput, cfg.Platform, cfg.PathFilterFunc), DaemonTag, PullTag),
		taggedProvider(containerd.NewDaemonProvider(tempDirGenerator, cfg.Registry, containerdClient.Namespace(), cfg.UserInput, cfg.Platform, cfg.PathFilterFunc), DaemonTag, PullTag),

		// registry providers
		taggedProvider(oci.NewRegistryProvider(tempDirGenerator, cfg.Registry, cfg.UserInput, cfg.Platform, cfg.PathFilterFunc), RegistryTag, PullTag),
	}
}

func taggedProvider(provider image.Provider, tags ...string) collections.TaggedValue[image.Provider] {
	return collections.NewTaggedValue[image.Provider](provider, append([]string{provider.Name()}, tags...)...)
}

func allProviderTags() []string {
	return collections.TaggedValueSet[image.Provider]{}.Join(ImageProviders(ImageProviderConfig{})...).Tags()
}
