package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/util"
	"path/filepath"
	"runtime"
)

func GetMutagenCommandPath() (string, error) {
	binDir, err := util.GetBinDir()

	if err != nil {
		return "", fmt.Errorf("could not get vessel bin dir: %w", err)
	}

	mutagenBinary := "mutagen"
	if runtime.GOOS == "windows" {
		mutagenBinary = mutagenBinary + ".exe"
	}

	return filepath.FromSlash(binDir + "/" + mutagenBinary), nil
}
