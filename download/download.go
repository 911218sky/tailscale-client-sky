package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"tailscale/utils/utilsTermbox"
)

const tailscaleDownloadURL = "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe"

var pm = utilsTermbox.Td.PrintMessage
var cm = utilsTermbox.Td.ClearMessage

// DownloadTailscale downloads the Tailscale installer with the specified fileName.
func DownloadTailscale(fileName string) error {
	pm("Downloading Tailscale...")

	resp, err := http.Get(tailscaleDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Get the file size
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return fmt.Errorf("Unable to get file size")
	}

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set the initial state of the progress bar
	percent := 0
	total := int64(0)
	printProgressBar := utilsTermbox.Td.ProgressBarAtY()
	printProgressBar(percent)

	// Copy file content to the file and manually update the progress bar
	buffer := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}

		// Update the progress bar
		total += int64(n)
		percent = int(float64(total) / float64(contentLength) * 100)
		printProgressBar(percent)

		if err == io.EOF {
			break
		}
	}

	pm("Download completed")
	return nil
}

// Install installs Tailscale using the specified downloadFileName.
func Install(downloadFileName string) error {
	installCmd := exec.Command(downloadFileName, "--install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	pm("Installing Tailscale...")

	if err := installCmd.Run(); err != nil {
		pm(fmt.Sprintf("Error installing Tailscale: %v\n", err))
		return err
	}

	pm("Tailscale installed successfully.")
	return nil
}
