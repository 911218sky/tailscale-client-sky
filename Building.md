## Prerequisites

Make sure you have the following installed:
- **Go programming language**: Follow the [official installation guide](https://golang.org/doc/install) to set it up.
- **UPX** (Optional): For compressing the executable, you can download it from the [UPX official site](https://upx.github.io/).

Ensure that Go is added to your system's PATH environment variable.

## Steps

### 1. Install Dependencies

To prepare your environment, you need to install the necessary tools:

#### Install `rsrc`

`rsrc` is used to embed Windows resources into Go executables.

```bash
go install github.com/akavel/rsrc@latest
```

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

### 3. Compile the Resource File (`*.syso`)

To create a resource file from a manifest, use the following command:

```bash
rsrc -manifest ./.manifest -o test.syso
```

#### Explanation:
- **`rsrc`**: A tool used to embed Windows resources into Go executables.
- **`-manifest ./.manifest`**: Specifies the manifest file to use. This file defines various settings for your executable, including compatibility and privileges.
- **`-o test.syso`**: Specifies the output name of the resource file, which will be automatically linked with your Go executable when you build it.

### 4. Add an Icon

To add an icon to your application, follow these steps:

#### Add the Icon

Run the following command to add an icon to your application:

```bash
go-winres simply --icon ./img/sky-tailscale-icon.png
```

#### Explanation:
- **`go-winres simply`**: A command to add resources (like icons) to your Go application.
- **`--icon ./img/sky-tailscale-icon.png`**: Specifies the path to your icon file that you want to include in your executable.

### 5. Rebuild the Executable with the Icon

After adding the icon and resource file, rebuild your application:

```bash
go build -ldflags "-s -w" -o sky-tailscale.exe
```

### 6. Compress the Executable (Optional)

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
