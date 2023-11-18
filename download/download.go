package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/cheggaaa/pb/v3"
)

const tailscaleDownloadURL = "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe"

func DownloadTailscale(fileName string) error {
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

	// 创建进度条
	bar := pb.Start64(contentLength)
	bar.SetTemplate(pb.Default)
	bar.SetMaxWidth(80)

	// 复制文件内容到文件并手动更新进度条
	buffer := make([]byte, 1024)
	total := int64(0)

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
		total += int64(n)
		bar.SetCurrent(total)
	}

	bar.Finish()
	return nil
}

func Install(downloadFileName string) error {
	// 创建安装命令
	installCmd := exec.Command(downloadFileName, "--install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	fmt.Println("Installing Tailscale...")

	// 执行安装命令并返回可能的错误
	if err := installCmd.Run(); err != nil {
		fmt.Printf("Error installing Tailscale: %v\n", err)
		return err
	}

	fmt.Println("Tailscale installed successfully.")
	return nil
}
