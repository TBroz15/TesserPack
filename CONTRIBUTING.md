# Contributing

TesserPack is kinda bare bones right now. PR contributions and how to develop on TesserPack is not completely considered until v1 is being released. Feel free to report any issues!

> [!WARNING]
> **ALL SETUPS FOR DEVELOPMENT ARE NOT TESTED AND CONSIDERED!**

## Before Sending Issues...

- Please check if your issue is already resolved, so your issue won't be closed as duplicate.
- If you are making a bug report, please make sure the reproduction process is accurate as possible.
- If your bug report is related to pack compilation, please upload and include your pack into your issue for better reproduction.
- TesserPack's intention was to optimize Minecraft packs for production. So make sure to stay on topic!

## Windows Setup (for development & PR)

> [!NOTE]
> Windows Defender and other anti-viruses tend to slow down the compilation process,
> since TesserPack isn't a well recognized program and think it's a sus program. (we are not sus ඞ)
>
> There are 3 solutions to this, 
> 1. Add Go's binaries and TesserPack binary to Windows Defender exclusions.
> 2. Temporarily disable the anti-virus. (not recommended)
> 3. Just simply use Linux. (windows bad, linux good -TuxeBro)
>
> [For more information...](https://go.dev/doc/faq#virus)

0.  Preparations
    1. [Download and install Go](https://go.dev/dl/) `<=` 1.25 if you haven't. Read the [guide](https://go.dev/doc/tutorial/getting-started#prerequisites).
    2. [Download Libvips](https://github.com/libvips/build-win64-mxe/releases/download/v8.17.1/vips-dev-w64-web-8.17.1.zip) and extract it.
    3. Go to the extracted directory and put `vips-dev-8.17` into `C:\`.
    4. Add Environmental Variables in Powershell. (you don't need epic admin powers btw...)
        ```powershell
        # Adds libvips to PATH
        [Environment]::SetEnvironmentVariable("Path", "$([Environment]::GetEnvironmentVariable("Path", "User"));C:\vips-dev-8.17\bin", "User")

        # Add libvips's pkgconfig dir to PKG_CONFIG_PATH, so govips will recognize it.
        [Environment]::SetEnvironmentVariable("PKG_CONFIG_PATH", "$([Environment]::GetEnvironmentVariable("PKG_CONFIG_PATH", "User"));C:\vips-dev-8.17\lib\pkgconfig", "User")
        
        # Set jsonv2 to GOEXPERIMENT, so TesserPack can use encoding/json/v2
        [Environment]::SetEnvironmentVariable("GOEXPERIMENT", "jsonv2", "User")
        ```

1. Fork this repository
2. In your repository, create a new branch with any name.
3. Clone your fork to your computer.
4. Run `go get -u ./...` to install the required dependencies.
5. Make some changes for your PR.
6. Run `go run ./cmd/tesserpack` or `go build ./cmd/tesserpack` then run `./tesserpack` to test it out.
7. Run `git commit` if you are done with your changes, like the usual.
8. Push your branch to GitHub via `git push origin (your branch name)`
9. Finally send your pull request to main. _It is recommended to draft your PR first if you are not sure with your code._


##