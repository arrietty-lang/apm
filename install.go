package apm

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//func DownloadTarGz(url, host, repoName, repoVersion string) (*os.File, error) {
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	path := filepath.Join(GetApmPath(), "packages", host, fmt.Sprintf("%s@%s", repoName, repoVersion))
//
//	out, err := os.Create(path)
//	if err != nil {
//		return "", err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, resp.Body)
//	if err != nil {
//		return "", err
//	}
//	return filepath.Abs(out.Name())
//}

// todo : simpleRelease消せゴミ
func InstallTarGz(tarGzUrl, host, author, repoName, repoVersion string) error {
	resp, err := http.Get(tarGzUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tarGzReader := resp.Body

	gzReader, err := gzip.NewReader(tarGzReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	var header *tar.Header

	tempInstallDir := filepath.Join(GetApmPath(), "packages", host, author)

	var installedDir string

	for {
		header, err = tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to Next(): %v", err)
		}

		fileName := filepath.Join(tempInstallDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fileName, 0755); err != nil {
				return fmt.Errorf("extractTarGz: Mkdir() failed: %v", err)
			}
			installedDir = header.Name
		case tar.TypeReg:
			outFile, err := os.Create(fileName)
			if err != nil {
				return fmt.Errorf("extractTarGz: Create() failed: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return fmt.Errorf("extractTarGz: Copy() failed: %v", err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("extractTarGz: Close() failed: %v", err)
			}

		case tar.TypeXGlobalHeader:
			continue

		default:
			return fmt.Errorf("extractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	realName := filepath.Join(tempInstallDir, fmt.Sprintf("%s@%s", repoName, repoVersion))
	err = os.Rename(filepath.Join(tempInstallDir, installedDir), realName)
	if err != nil {
		return err
	}
	return nil
}
