package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/util"
	"os"
	"path/filepath"
	"runtime"
)

// DownloadMutagen gets the mutagen binary and related and puts it into
// ~/.vessel/bin
func DownloadMutagen() error {
	StopMutagenDaemon()
	flavor := runtime.GOOS + "_" + runtime.GOARCH

	mutagenBinDir, err := util.GetBinDir()

	if err != nil {
		return fmt.Errorf("could not get mutagen binary dir: %w", err)
	}

	mutagenBinFile, err := GetMutagenPath()

	if err != nil {
		return fmt.Errorf("could not get mutagen binary path: %w", err)
	}

	globalMutagenDir := filepath.Dir(mutagenBinDir)
	destFile := filepath.FromSlash(globalMutagenDir + "/mutagen.tgz")
	mutagenURL := fmt.Sprintf("https://github.com/mutagen-io/mutagen/releases/download/v%s/mutagen_%s_v%s.tar.gz", "0.15.1", flavor, "0.15.1")

	// Remove the existing binary, if exists
	_ = os.Remove(mutagenBinFile)

	_ = os.MkdirAll(globalMutagenDir, 0777)
	err := util.DownloadFile(destFile, mutagenURL, "true" != os.Getenv("DDEV_NONINTERACTIVE"))
	if err != nil {
		return err
	}

	err = archive.Untar(destFile, globalMutagenDir, "")
	_ = os.Remove(destFile)
	if err != nil {
		return err
	}
	err = os.Chmod(globalconfig.GetMutagenPath(), 0755)
	if err != nil {
		return err
	}

	// Stop daemon in case it was already running somewhere else
	StopMutagenDaemon()
	return nil
}
