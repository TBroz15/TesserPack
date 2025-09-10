# Contributing

TesserPack is kinda bare bones right now. PR contributions and how to develop on TesserPack is not completely considered until v1 is being released. Feel free to report any issues!

## Before Sending Issues...

- Please check if your issue is already resolved, so your issue won't be closed as duplicate.
- If you are making a bug report, please make sure the reproduction process is accurate as possible.
- If your bug report is related to pack compilation, please upload and include your pack into your issue for better reproduction.
- TesserPack's intention was to optimize Minecraft packs for production. So make sure to stay on topic!

## Development (PR)

These setup instructions assumes that you don't have Go, Libvips, and GCC installed. 

This guide will definitely help you save some time without countless trials and errors trying to figure out how to setup TesserPack.

### Windows Setup

> [!NOTE]
> Windows Defender and other anti-viruses tend to slow down the compilation process,
> since TesserPack isn't a well recognized program and think it's a suspicious program. (we are not sus ඞ)
>
> There are 3 solutions to this, 
> 1. Add Go's binaries and TesserPack binary to Windows Defender exclusions.
> 2. Temporarily disable the anti-virus. (not recommended)
> 3. Just simply use Linux. (windows bad, linux good -TuxeBro)
>
> [For more information...](https://go.dev/doc/faq#virus)

1. [Download and install Go](https://go.dev/dl/) `>=` 1.25 if you haven't. Read the [guide](https://go.dev/doc/tutorial/getting-started#prerequisites) if you also haven't.
2. [Download Libvips](https://github.com/libvips/build-win64-mxe/releases/download/v8.17.1/vips-dev-w64-web-8.17.1.zip) and extract it.
3. Go to the extracted directory and put `vips-dev-8.17` into `C:\`.
4. [Download and install MSYS2](https://www.msys2.org/#installation) and follow installation wizard if not installed. Make sure that `Run MSYS2 now` is checked before finishing.
5. In MSYS2, install some packages by `pacman -S --needed base-devel mingw-w64-ucrt-x86_64-toolchain`. Just press enter and then enter `Y` to install.
6. Add Environmental Variables in Powershell and restart your terminal. (you don't need epic admin powers btw...)
    ```powershell
    # Note: This will make the variables globally available on your every terminal.

    # Adds gcc to PATH
    [Environment]::SetEnvironmentVariable("Path", "$([Environment]::GetEnvironmentVariable("Path", "User"));C:\msys64\ucrt64\bin", "User")

    # Adds libvips to PATH
    [Environment]::SetEnvironmentVariable("Path", "$([Environment]::GetEnvironmentVariable("Path", "User"));C:\vips-dev-8.17\bin", "User")

    # Add libvips's pkgconfig dir to PKG_CONFIG_PATH, so govips will recognize it.
    [Environment]::SetEnvironmentVariable("PKG_CONFIG_PATH", "$([Environment]::GetEnvironmentVariable("PKG_CONFIG_PATH", "User"));C:\vips-dev-8.17\lib\pkgconfig", "User")
    
    # Set jsonv2 to GOEXPERIMENT, so TesserPack can use encoding/json/v2
    [Environment]::SetEnvironmentVariable("GOEXPERIMENT", "jsonv2", "User")

    # Enable CGO so TesserPack can use govips/libvips
    [Environment]::SetEnvironmentVariable("CGO_ENABLED", "1", "User")
    ```

### Linux Setup

These setups are relatively similar for every other distributions that are not included here.

#### Fedora 42+

1. [Download Go](https://go.dev/dl/) `>=` 1.25 and [install](https://go.dev/doc/install#install) it if you haven't. Read the [guide](https://go.dev/doc/tutorial/getting-started#prerequisites) if you also haven't.

    *Not recommended* but alternatively, you can install the latest version of Go from Rawhide.
    ```bash
    sudo dnf install fedora-repos-rawhide
    sudo dnf install golang --enablerepo=rawhide
    ```

2. Uninstall Fedora's Libvips (Optional)
    ```bash
    sudo dnf remove vips vips-devel
    ```

3. Install Build Dependencies
    ```bash
    sudo dnf group install "development-tools"

    sudo dnf install meson ninja-build pkg-config gcc gcc-c++ \
        glib2-devel expat-devel libjpeg-turbo-devel libpng-devel \
        libwebp-devel libtiff-devel giflib-devel lcms2-devel \
        orc-devel poppler-glib-devel librsvg2-devel \
        ImageMagick-c++-devel libheif-devel libexif-devel
    ```

4. Download Libvips
    ```bash
    wget https://github.com/libvips/libvips/releases/download/v8.17.1/vips-8.17.1.tar.xz
    tar -xf vips-8.17.1.tar.xz
    cd vips-8.17.1
    ```

5. Setup & Install
    ```bash
    meson setup builddir --prefix=/usr
    cd builddir
    ninja
    sudo ninja install
    sudo ldconfig
    ```

6. Verify if Libvips was installed
    ```bash
    vips --version
    ```

7. Set GOEXPERIMENT to jsonv2
    ```bash
    export GOEXPERIMENT=jsonv2
    ```
    You could add this into `.bashrc` or any shell configuration file.

#### Ubuntu 24+

1. [Download Go](https://go.dev/dl/) `>=` 1.25 and [install](https://go.dev/doc/install#install) it if you haven't. Read the [guide](https://go.dev/doc/tutorial/getting-started#prerequisites) if you also haven't.

2. Uninstall Ubuntu's Libvips (Optional)
    ```bash
    sudo apt remove libvips libvips-dev
    ```

3. Install Build Dependencies
    ```bash
    sudo apt update
    
    sudo apt install -y build-essential meson ninja-build pkg-config \
        libglib2.0-dev libexpat1-dev libjpeg-turbo8-dev libpng-dev \
        libwebp-dev libtiff-dev libgif-dev liblcms2-dev liborc-0.4-dev \
        librsvg2-dev libheif-dev libexif-dev libpoppler-glib-dev \
        libimagequant-dev
    ```

4. Download Libvips
    ```bash
    wget https://github.com/libvips/libvips/releases/download/v8.17.1/vips-8.17.1.tar.xz
    tar -xf vips-8.17.1.tar.xz
    cd vips-8.17.1
    ```

5. Setup & Install
    ```bash
    meson setup builddir --prefix=/usr/local
    cd builddir
    ninja
    sudo ninja install
    sudo ldconfig
    ```

6. Verify if Libvips was installed
    ```bash
    vips --version
    ```

7. Set GOEXPERIMENT to jsonv2
    ```bash
    export GOEXPERIMENT=jsonv2
    ```
    You could add this into `.bashrc` or any shell configuration file.

### After Setup, Steps:

1. Fork this repository
2. In your repository, create a new branch with any name.
3. Clone your fork to your computer.
4. Run `go get -u ./...` to install the required dependencies.
5. Make some changes for your PR.
6. Run `go run ./cmd/tesserpack` or `go build ./cmd/tesserpack` then run `./tesserpack` to test it out.
7. Run `git commit` if you are done with your changes, like the usual. [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) are followed.
8. Push your branch to GitHub via `git push origin (your branch name)`
9. Finally send your pull request to main. _It is recommended to draft your PR first if you are not sure with your code._