## Steps

### 1. Install Dependencies

To prepare your environment, you need to install the necessary tools:

#### Install `go-winres`

This package is used to add icons to your Go application.

```bash
go install github.com/tc-hib/go-winres@latest
```

### 2. Compile the Go Executable

To compile your Go program, use the following command:

```bash
go build -ldflags "-s -w" -o sky-tailscale.exe main.go
```

### 3. Add an Icon

To add an icon to your application, follow these steps:

#### Add the Icon

Run the following command to add an icon to your application:

```bash
go-winres simply --icon ./img/sky-tailscale-icon.png
```

#### Explanation:
- **`go-winres simply`**: A command to add resources (like icons) to your Go application.
- **`--icon ./img/sky-tailscale-icon.png`**: Specifies the path to your icon file that you want to include in your executable.

### 4. Rebuild the Executable with the Icon

After adding the icon, rebuild your application:

```bash
go build -ldflags "-s -w" -o sky-tailscale.exe
```

### 5. Compress the Executable (Optional)

Optionally, you can compress the final executable using UPX for a smaller file size:

```bash
upx -9 sky-tailscale.exe
```

#### Explanation:
- **`upx -9 sky-tailscale.exe`**: `upx` is a command-line tool to compress executable files. The `-9` flag specifies the maximum compression level.

## Additional References

- [Go Documentation](https://golang.org/doc/)
- [UPX Official Site](https://upx.github.io/)
- [Embedding Resources in Go](https://stackoverflow.com/questions/25602600/how-do-you-set-the-application-icon-in-golang)