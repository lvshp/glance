package lib

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const readcliRepo = "lvshp/ReadCLI"

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseInfo struct {
	TagName string         `json:"tag_name"`
	Name    string         `json:"name"`
	Body    string         `json:"body"`
	Assets  []ReleaseAsset `json:"assets"`
}

func FetchLatestRelease(version string) (*ReleaseInfo, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+readcliRepo+"/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ReadCLI/"+strings.TrimSpace(version))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("检查更新失败: %s", resp.Status)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func ShouldOfferUpdate(currentVersion, latestVersion string) bool {
	currentVersion = strings.TrimSpace(currentVersion)
	latestVersion = strings.TrimSpace(latestVersion)
	if currentVersion == "" || latestVersion == "" || currentVersion == latestVersion {
		return false
	}

	currentSemver, currentOK := parseSemverLike(currentVersion)
	latestSemver, latestOK := parseSemverLike(latestVersion)
	if currentOK && latestOK {
		for i := 0; i < 3; i++ {
			if latestSemver[i] > currentSemver[i] {
				return true
			}
			if latestSemver[i] < currentSemver[i] {
				return false
			}
		}
		return false
	}
	if latestOK && !currentOK {
		return false
	}

	return latestVersion != currentVersion
}

func SelectReleaseAsset(release *ReleaseInfo, goos, goarch string) *ReleaseAsset {
	if release == nil {
		return nil
	}
	prefix := fmt.Sprintf("readcli-%s-%s-", goos, goarch)
	for i := range release.Assets {
		asset := &release.Assets[i]
		if strings.HasPrefix(asset.Name, prefix) && strings.HasSuffix(asset.Name, ".tar.gz") {
			return asset
		}
	}
	return nil
}

func InstallLatestReleaseAsset(version, url, executablePath string) error {
	client := &http.Client{Timeout: 90 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("User-Agent", "ReadCLI/"+strings.TrimSpace(version))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载更新失败: %s", resp.Status)
	}

	execDir := filepath.Dir(executablePath)
	tempFile, err := os.CreateTemp(execDir, "readcli-update-*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		_ = os.Remove(tempPath)
	}()

	info, statErr := os.Stat(executablePath)
	fileMode := os.FileMode(0755)
	if statErr == nil {
		fileMode = info.Mode()
	}
	if err := extractBinaryFromTarGz(resp.Body, tempFile); err != nil {
		return err
	}
	if err := tempFile.Chmod(fileMode); err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, executablePath)
}

func CurrentPlatformSupported() bool {
	asset := SelectReleaseAsset(&ReleaseInfo{
		Assets: []ReleaseAsset{
			{Name: "readcli-darwin-amd64-v0.0.0.tar.gz"},
			{Name: "readcli-darwin-arm64-v0.0.0.tar.gz"},
			{Name: "readcli-linux-amd64-v0.0.0.tar.gz"},
		},
	}, runtime.GOOS, runtime.GOARCH)
	return asset != nil
}

func extractBinaryFromTarGz(source io.Reader, target io.Writer) error {
	gzReader, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return errors.New("更新包中未找到 readcli 可执行文件")
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if filepath.Base(header.Name) != "readcli" {
			continue
		}
		_, err = io.Copy(target, tarReader)
		return err
	}
}

func parseSemverLike(version string) ([3]int, bool) {
	version = strings.TrimSpace(strings.TrimPrefix(version, "v"))
	version = strings.TrimSpace(strings.TrimPrefix(version, "V"))
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return [3]int{}, false
	}

	var out [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		value, err := strconv.Atoi(parts[i])
		if err != nil || value < 0 {
			return [3]int{}, false
		}
		out[i] = value
	}
	return out, true
}
