package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"tailscale/utils/drawer"
)

const tailscaleDownloadURL = "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe"

// DownloadTailscale downloads the Tailscale installer with the specified fileName.
func DownloadTailscale(fileName string) error {
	drawer.Print("Downloading Tailscale...", drawer.DefaultOption)

	resp, err := http.Get(tailscaleDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s", resp.Status)
	}

	// Get the file size
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return fmt.Errorf("unable to get file size")
	}

	// Create the destination file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize progress bar state
	percent := 0
	total := int64(0)
	y := drawer.GetY()
	drawer.DrawProgressBar(y, percent, drawer.DefaultOption)

	// Buffer for reading response body
	buffer := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Write to the file
		if _, err := file.Write(buffer[:n]); err != nil {
			return err
		}

		// Update progress
		total += int64(n)
		percent = int(float64(total) / float64(contentLength) * 100)
		drawer.DrawProgressBar(y, percent, drawer.DefaultOption)

		if err == io.EOF {
			break
		}
	}

	drawer.NextLine()
	drawer.Print("Download completed", drawer.DefaultOption)
	return nil
}

// Install installs Tailscale using the specified download file name.
func Install(downloadFileName string) error {
	installCmd := exec.Command(downloadFileName, "--install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	drawer.Print("Installing Tailscale...", drawer.DefaultOption)

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf(fmt.Sprintf("Error installing Tailscale: %v\n", err))
	}

	drawer.Print("Tailscale installed successfully.", drawer.DefaultOption)
	return nil
}
