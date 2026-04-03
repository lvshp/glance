package lib

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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

type UpdateInstallError struct {
	Message string
	TempDir string
}

func (e *UpdateInstallError) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.TempDir) == "" {
		return e.Message
	}
	return fmt.Sprintf("%s。临时文件保留在：%s", e.Message, e.TempDir)
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
	suffix := releaseArchiveSuffix(goos)
	for i := range release.Assets {
		asset := &release.Assets[i]
		if strings.HasPrefix(asset.Name, prefix) && strings.HasSuffix(asset.Name, suffix) {
			return asset
		}
	}
	return nil
}

func InstallLatestReleaseAsset(version, url, executablePath string) error {
	goos := runtime.GOOS
	tempDir, err := os.MkdirTemp("", "readcli-update-*")
	if err != nil {
		return err
	}
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.RemoveAll(tempDir)
		}
	}()

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
	archivePath := filepath.Join(tempDir, "readcli-update"+releaseArchiveSuffix(goos))
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(archiveFile, resp.Body); err != nil {
		archiveFile.Close()
		cleanup = false
		return &UpdateInstallError{Message: "写入更新包失败: " + err.Error(), TempDir: tempDir}
	}
	if err := archiveFile.Close(); err != nil {
		cleanup = false
		return &UpdateInstallError{Message: "保存更新包失败: " + err.Error(), TempDir: tempDir}
	}

	info, statErr := os.Stat(executablePath)
	fileMode := os.FileMode(0755)
	if statErr == nil {
		fileMode = info.Mode()
	}
	binaryPath := filepath.Join(tempDir, releaseBinaryName(goos))
	if err := extractBinaryFromArchive(archivePath, binaryPath, goos, fileMode); err != nil {
		cleanup = false
		return &UpdateInstallError{Message: err.Error(), TempDir: tempDir}
	}

	if err := verifyExecutableDirWritable(execDir); err != nil {
		cleanup = false
		return &UpdateInstallError{Message: err.Error(), TempDir: tempDir}
	}

	if goos == "windows" {
		if err := stageWindowsReplacement(binaryPath, executablePath, tempDir); err != nil {
			cleanup = false
			return &UpdateInstallError{Message: err.Error(), TempDir: tempDir}
		}
		cleanup = false
		return nil
	}

	if err := os.Rename(binaryPath, executablePath); err != nil {
		cleanup = false
		return &UpdateInstallError{Message: "覆盖当前二进制失败: " + err.Error(), TempDir: tempDir}
	}
	return nil
}

func CurrentExecutablePath() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func CurrentPlatformSupported() bool {
	asset := SelectReleaseAsset(&ReleaseInfo{
		Assets: []ReleaseAsset{
			{Name: "readcli-darwin-amd64-v0.0.0.tar.gz"},
			{Name: "readcli-darwin-arm64-v0.0.0.tar.gz"},
			{Name: "readcli-linux-amd64-v0.0.0.tar.gz"},
			{Name: "readcli-windows-amd64-v0.0.0.zip"},
		},
	}, runtime.GOOS, runtime.GOARCH)
	return asset != nil
}

func releaseBinaryName(goos string) string {
	if goos == "windows" {
		return "readcli.exe"
	}
	return "readcli"
}

func releaseArchiveSuffix(goos string) string {
	if goos == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}

func extractBinaryFromArchive(archivePath, targetPath, goos string, fileMode os.FileMode) error {
	if goos == "windows" {
		return extractBinaryFromZip(archivePath, targetPath, releaseBinaryName(goos))
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建临时二进制失败: %w", err)
	}
	defer out.Close()

	archiveReader, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("打开更新包失败: %w", err)
	}
	defer archiveReader.Close()

	if err := extractBinaryFromTarGz(archiveReader, out); err != nil {
		return err
	}
	if err := out.Chmod(fileMode); err != nil {
		return fmt.Errorf("设置更新文件权限失败: %w", err)
	}
	return nil
}

func verifyExecutableDirWritable(execDir string) error {
	if probe, err := os.CreateTemp(execDir, "readcli-write-test-*"); err != nil {
		return errors.New("当前二进制目录不可写，无法自动覆盖")
	} else {
		probePath := probe.Name()
		probe.Close()
		_ = os.Remove(probePath)
	}
	return nil
}

func stageWindowsReplacement(binaryPath, executablePath, tempDir string) error {
	scriptPath := filepath.Join(tempDir, "replace.cmd")
	script := fmt.Sprintf(`@echo off
setlocal
set "TARGET=%s"
set "SOURCE=%s"
:retry
move /Y "%%SOURCE%%" "%%TARGET%%" >nul 2>&1
if errorlevel 1 (
  timeout /t 1 /nobreak >nul
  goto retry
)
endlocal
`, executablePath, binaryPath)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("写入更新脚本失败: %w", err)
	}

	cmd := exec.Command("cmd.exe", "/C", "start", "", "/MIN", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动更新脚本失败: %w", err)
	}
	return nil
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

func extractBinaryFromZip(archivePath, targetPath, binaryName string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if filepath.Base(file.Name) != binaryName {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		out, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("创建临时二进制失败: %w", err)
		}
		defer out.Close()

		if _, err := io.Copy(out, rc); err != nil {
			return err
		}
		return nil
	}
	return errors.New("更新包中未找到 readcli 可执行文件")
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
