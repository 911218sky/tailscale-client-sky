### 1. Install Dependencies

To prepare your environment, you need to install the necessary tools:

#### Install `go-winres`

This package is used to add icons to your Go application.

```bash
go install github.com/tc-hib/go-winres@latest
```

#### Install UPX (Optional)

UPX is used to compress the executable file for smaller size. You can install it as follows:

- **On Linux:**
  ```bash
  sudo apt install upx
  ```

- **On macOS (via Homebrew):**
  ```bash
  brew install upx
  ```

---

### 2. Compile the Go Executable

To compile your Go program, use the following commands for **Windows** and **Linux**.

#### For Windows:

```bash
go build -ldflags "-s -w" -o sky-tailscale.exe main.go
```

#### For Linux:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o sky-tailscale-linux main.go
```

- **`GOOS=linux`**: Specifies the target operating system as Linux.
- **`GOARCH=amd64`**: Specifies the target architecture as 64-bit.

---

### 3. Add an Icon (Windows Only)

To add an icon to your Windows application, follow these steps:

#### Add the Icon

Run the following command to add an icon to your Windows executable:

```bash
go-winres simply --icon ./img/sky-tailscale-icon.png
```

#### Explanation:
- **`go-winres simply`**: A command to add resources (like icons) to your Go application.
- **`--icon ./img/sky-tailscale-icon.png`**: Specifies the path to your icon file that you want to include in your executable.

---

### 4. Rebuild the Executable with the Icon (Windows Only)

After adding the icon, rebuild your Windows application:

```bash
go build -ldflags "-s -w" -o sky-tailscale.exe
```

---

### 5. Compress the Executable (Optional)

Optionally, you can compress the final executable for both **Windows** and **Linux** versions using UPX:

#### Compress Windows Executable:

```bash
upx -9 sky-tailscale.exe
```

#### Compress Linux Executable:

```bash
upx -9 sky-tailscale-linux
```

#### Explanation:
- **`upx -9`**: Compresses the executable file at the maximum compression level.

---

### 6. Cross-Platform Build (Optional)

If you need to build for multiple platforms in one step, you can use the following command:

```bash
gox -os="windows linux" -arch="amd64" -output="sky-tailscale-{{.OS}}-{{.Arch}}"
```

- **`gox`**: A cross-compilation tool for Go.
- **`-os="windows linux"`**: Specifies the target operating systems.
- **`-arch="amd64"`**: Specifies the target architecture.
- **`-output="sky-tailscale-{{.OS}}-{{.Arch}}"`**: Formats the output file names to include the OS and architecture.

Install `gox` if not already installed:

```bash
go install github.com/mitchellh/gox@latest
```

---

## Additional References

- [Go Documentation](https://golang.org/doc/)
- [UPX Official Site](https://upx.github.io/)
- [Cross Compilation in Go](https://golang.org/doc/install/source#environment)  
- [Embedding Resources in Go](https://stackoverflow.com/questions/25602600/how-do-you-set-the-application-icon-in-golang)