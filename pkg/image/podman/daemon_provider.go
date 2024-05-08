package podman

import (
	"github.com/docker/docker/client"

	"github.com/anchore/stereoscope/internal/podman"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/stereoscope/pkg/image/docker"
	"github.com/anchore/stereoscope/pkg/pathfilter"
)

const Daemon image.Source = image.PodmanDaemonSource

func NewDaemonProvider(tmpDirGen *file.TempDirGenerator, imageStr string, platform *image.Platform, pathFilterFunc pathfilter.PathFilterFunc) image.Provider {
	return docker.NewAPIClientProvider(Daemon, tmpDirGen, imageStr, platform, pathFilterFunc,
		func() (client.APIClient, error) {
			return podman.GetClient()
		})
}
