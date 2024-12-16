package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"tailscale/utils/drawer"
)

const (
	// TailscaleWindowsURL is the URL to download the latest stable version of Tailscale for Windows.
	// Exported as it might be useful for other packages to know the download URL.
	TailscaleWindowsURL = "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe"
	bufferSize          = 1024
)

// DownloadTailscaleWindows downloads the Tailscale installer for Windows and saves it to the specified file.
// It displays a progress bar during download.
func DownloadTailscaleWindows(fileName string) error {
	drawer.Print("Downloading Tailscale...", drawer.DefaultOption)

	resp, err := http.Get(TailscaleWindowsURL)
	if err != nil {
		return fmt.Errorf("failed to download Tailscale: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return fmt.Errorf("invalid content length received")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := downloadWithProgress(resp.Body, file, contentLength); err != nil {
		return fmt.Errorf("failed during download: %w", err)
	}

	drawer.NextLine()
	drawer.Print("Download completed", drawer.DefaultOption)
	return nil
}

// downloadWithProgress copies data from src to dst while updating a progress bar.
func downloadWithProgress(src io.Reader, dst io.Writer, total int64) error {
	y := drawer.GetY()
	buffer := make([]byte, bufferSize)
	var downloaded int64

	for {
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dst.Write(buffer[:n]); err != nil {
			return err
		}

		downloaded += int64(n)
		percent := int(float64(downloaded) / float64(total) * 100)
		drawer.DrawProgressBar(y, percent, drawer.DefaultOption)

		if err == io.EOF {
			break
		}
	}
	return nil
}

// DownloadTailscaleLinux downloads and installs Tailscale for Linux using the official install script.
func DownloadTailscaleLinux() error {
	drawer.Print("Downloading and installing Tailscale for Linux...", drawer.DefaultOption)

	cmd := exec.Command("sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Tailscale: %w", err)
	}

	drawer.Print("Tailscale downloaded and installed successfully.", drawer.DefaultOption)
	return nil
}

// Install runs the Tailscale installer executable on Windows.
func Install(downloadFileName string) error {
	drawer.Print("Installing Tailscale...", drawer.DefaultOption)

	cmd := exec.Command(downloadFileName, "--install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run installer: %w", err)
	}

	drawer.Print("Tailscale installed successfully.", drawer.DefaultOption)
	return nil
}
