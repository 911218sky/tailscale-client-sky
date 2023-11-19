package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"tailscale/utilsTermbox"
)

const tailscaleDownloadURL = "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe"

func DownloadTailscale(fileName string) error {
	utilsTermbox.PrintMessage("Downloading tailscale ...")

	resp, err := http.Get(tailscaleDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 获取文件大小
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return fmt.Errorf("Unable to get file size")
	}

	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 设置进度条的初始状态
	percent := 0
	total := int64(0)
	utilsTermbox.ProgressBarInit()
	utilsTermbox.PrintProgressBar(percent)

	// 复制文件内容到文件并手动更新进度条
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

		// 更新进度条
		total += int64(n)
		percent = int(float64(total) / float64(contentLength) * 100)
		utilsTermbox.PrintProgressBar(percent)

		if err == io.EOF {
			break
		}
	}

	utilsTermbox.PrintMessage("Download completed")
	return nil
}

func Install(downloadFileName string) error {
	installCmd := exec.Command(downloadFileName, "--install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	utilsTermbox.PrintMessage("Installing Tailscale...")

	if err := installCmd.Run(); err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Error installing Tailscale: %v\n", err))
		return err
	}

	utilsTermbox.PrintMessage("Tailscale installed successfully.")
	return nil
}
