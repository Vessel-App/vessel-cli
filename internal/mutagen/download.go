package mutagen

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/util"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// InstallMutagen gets the mutagen binary and related and puts it into
// ~/.vessel/bin
func InstallMutagen() error {
	flavor := runtime.GOOS + "_" + runtime.GOARCH

	mutagenBinDir, err := util.GetBinDir()

	if err != nil {
		return fmt.Errorf("could not get mutagen binary dir: %w", err)
	}

	mutagenBinFile, err := GetMutagenCommandPath()

	if err != nil {
		return fmt.Errorf("could not get mutagen binary path: %w", err)
	}

	// If mutagen is already installed
	if util.FileExists(mutagenBinFile) {
		return nil
	}

	StopMutagenDaemon()
	destFile := filepath.FromSlash(mutagenBinDir + "/mutagen.tgz")
	mutagenURL := fmt.Sprintf("https://github.com/mutagen-io/mutagen/releases/download/v%s/mutagen_%s_v%s.tar.gz", "0.15.1", flavor, "0.15.1")

	// Remove the existing binary, if exists
	_ = os.Remove(mutagenBinFile)

	_ = os.MkdirAll(mutagenBinDir, 0777)
	err = downloadMutagen(destFile, mutagenURL)

	if err != nil {
		return fmt.Errorf("could not download mutagen: %w", err)
	}

	err = untar(destFile, mutagenBinDir)

	if err != nil {
		return err
	}

	_ = os.Remove(destFile)

	err = os.Chmod(mutagenBinFile, 0755)

	if err != nil {
		return err
	}

	// Stop daemon in case it was already running somewhere else
	StopMutagenDaemon()
	return nil
}

func downloadMutagen(destination, url string) error {
	out, err := os.Create(destination)

	if err != nil {
		return fmt.Errorf("could not create mutagen destination file: %w", err)
	}

	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("could not download mutagen binary: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download link %s returned wrong status code: got %v want %v", url, resp.StatusCode, http.StatusOK)
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return fmt.Errorf("could not copy mutagen file to its destination: %w", err)
	}

	return nil
}

func untar(src, dest string) error {
	var tf *tar.Reader

	f, err := os.Open(src)

	if err != nil {
		return fmt.Errorf("could not open tar file for reading: %w", err)
	}

	defer f.Close()

	gf, err := gzip.NewReader(f)

	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}

	gf.Close()

	tf = tar.NewReader(gf)

	for {
		file, err := tf.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error during read of tar archive %v, err: %v", src, err)
		}

		// If file.Name is now empty this is the root directory we want to extract, and need not do anything.
		if file.Name == "" && file.Typeflag == tar.TypeDir {
			continue
		}

		fullPath := filepath.Join(dest, file.Name)

		// At this point only directories and block-files are handled. Symlinks and the like are ignored.
		switch file.Typeflag {
		case tar.TypeDir:
			// For a directory, if it doesn't exist, we create it.
			finfo, err := os.Stat(fullPath)
			if err == nil && finfo.IsDir() {
				continue
			}

			err = os.MkdirAll(fullPath, 0755)
			if err != nil {
				return err
			}

			err = os.Chmod(fullPath, fs.FileMode(file.Mode))
			if err != nil {
				return fmt.Errorf("failed to chmod %v dir %v, err: %v", fs.FileMode(file.Mode), fullPath, err)
			}

		case tar.TypeReg:
			fallthrough
		case tar.TypeRegA:
			// Always ensure the directory is created before trying to move the file.
			fullPathDir := filepath.Dir(fullPath)
			err = os.MkdirAll(fullPathDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create the directory %s, err: %v", fullPathDir, err)
			}

			// For a regular file, create and copy the file.
			exFile, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("failed to create file %v, err: %v", fullPath, err)
			}
			_, err = io.Copy(exFile, tf)
			_ = exFile.Close()
			if err != nil {
				return fmt.Errorf("failed to copy to file %v, err: %v", fullPath, err)
			}
			err = os.Chmod(fullPath, fs.FileMode(file.Mode))
			if err != nil {
				return fmt.Errorf("failed to chmod %v file %v, err: %v", fs.FileMode(file.Mode), fullPath, err)
			}

		}
	}

	return nil
}
