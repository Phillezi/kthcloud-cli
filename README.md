<div align="center">
<pre>
   __    __    __         __                __             __   _ 
  / /__ / /_  / /  ____  / / ___  __ __ ___/ / ____ ____  / /  (_)
 /  '_// __/ / _ \/ __/ / / / _ \/ // // _  / /___// __/ / /  / / 
/_/\_\ \__/ /_//_/\__/ /_/  \___/\_,_/ \_,_/       \__/ /_/  /_/  
                                                                  
</pre>

[![Windows](https://img.shields.io/badge/Windows-FFFFFF?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiBoZWlnaHQ9IjgwMHB4IiB3aWR0aD0iODAwcHgiIHZlcnNpb249IjEuMSIgaWQ9IkNhcGFfMSIgdmlld0JveD0iMCAwIDE5LjEzMiAxOS4xMzIiIHhtbDpzcGFjZT0icHJlc2VydmUiPgo8Zz4KCTxnPgoJCTxwYXRoIHN0eWxlPSJmaWxsOiMwMzAxMDQ7IiBkPSJNOS4xNzIsOS4xNzlWMC4xNDZIMHY5LjAzM0g5LjE3MnoiLz4KCQk8cGF0aCBzdHlsZT0iZmlsbDojMDMwMTA0OyIgZD0iTTE5LjEzMiw5LjE3OVYwLjE0Nkg5Ljk1OXY5LjAzM0gxOS4xMzJ6Ii8+CgkJPHBhdGggc3R5bGU9ImZpbGw6IzAzMDEwNDsiIGQ9Ik0xOS4xMzIsMTguOTg2VjkuOTU1SDkuOTU5djkuMDMySDE5LjEzMnoiLz4KCQk8cGF0aCBzdHlsZT0iZmlsbDojMDMwMTA0OyIgZD0iTTkuMTcyLDE4Ljk4NlY5Ljk1NUgwdjkuMDMySDkuMTcyeiIvPgoJPC9nPgo8L2c+Cjwvc3ZnPg==)](#windows)
[![Linux](https://img.shields.io/badge/Linux-FFFFFF?style=for-the-badge&logo=linux&logoColor=black)](#mac-and-linux)
[![macOS](https://img.shields.io/badge/mac%20os-FFFFFF?style=for-the-badge&logo=apple&logoColor=black)](#mac-and-linux)

[![Go Report Card](https://goreportcard.com/badge/github.com/Phillezi/kthcloud-cli?style=social)](https://goreportcard.com/report/github.com/Phillezi/kthcloud-cli)

</div>

# kthcloud-cli

> [!NOTE]  
> This project is in the very early stages of development. Features are incomplete, and things may change frequently.

## Table of Contents

- [Overview](#overview)
  - [Compose](#compose)
- [Installation](#installation)
  - [Download binary](#download-and-install-binary)
    - [Mac and Linux](#mac-and-linux)
    - [Windows](#windows)
  - [Build](#build-it-yourself)
- [Commands](#commands)
  - [Login](#login-command)
  - [Compose](#compose-command)
    - [Up](#compose-up-command)
    - [Down](#compose-down-command)
    - [Parse](#compose-parse-command)
  - [Update](#update-command)
  - [Version](#version-command)
- [Configuration](#configuration)

## Overview

`kthcloud-cli` is a command-line interface tool for interacting with kthcloud’s API. It allows you to perform various operations such as listin deployments, creating api keys, and creating deployments from `docker-compose` files.

<div align="center">
    
![example](https://github.com/user-attachments/assets/9482fba7-a50d-4502-8d80-1319b932dfe1)

</div>

### Compose

The clis core functionality is to parse `docker compose` files and create deployments from the content.

For example, if i have this `docker-compose.yaml` file:

```yaml
services:
  file-server:
    image: phillezi/tinyhttpfileserver:latest
    environment:
      KTHCLOUD_CORES: 0.1
      KTHCLOUD_RAM: 0.1
      KTHCLOUD_VISIBILITY: auth
      KTHCLOUD_HEALTH_PATH: "/"
    ports:
      - "8080:8080"
    volumes:
      - "./testpath:/static"

  log-test-app1:
    image: phillezi/litelogger:latest
    environment:
      KTHCLOUD_CORES: 0.1
      KTHCLOUD_RAM: 0.1
      KTHCLOUD_VISIBILITY: private
    depends_on:
      - file-server

  log-test-app2:
    image: phillezi/litelogger:latest
    environment:
      KTHCLOUD_CORES: 0.1
      KTHCLOUD_RAM: 0.1
      KTHCLOUD_VISIBILITY: private
    depends_on: [log-test-app1]

  log-test-app3:
    image: phillezi/litelogger:latest
    environment:
      KTHCLOUD_CORES: 0.1
      KTHCLOUD_RAM: 0.1
      KTHCLOUD_VISIBILITY: private

```

> [!NOTE]  
> The above example showcases what is supported but does not provide a functional application. You need to have the ./testpath in your `cwd`.

The tool will create four deployments and set up their environment variables, port, start commands and persistent storage.

## Installation

### Download and install binary

#### Mac and Linux

For Mac and Linux, there is an installation script that can be run to install the CLI.

##### Prerequisites

- bash
- curl

```bash
curl -fsSL https://raw.githubusercontent.com/Phillezi/kthcloud-cli/main/scripts/install.sh | bash

```

Check out what the script does [here](https://github.com/Phillezi/kthcloud-cli/blob/main/scripts/install.sh).

#### Windows

There is a PowerShell installation script that can be run to install the CLI.

```powershell
powershell -c "irm https://raw.githubusercontent.com/Phillezi/kthcloud-cli/main/scripts/install.ps1 | iex"

```

Check out what the script does [here](https://github.com/Phillezi/kthcloud-cli/blob/main/scripts/install.ps1).

### Installing using `go install`

If you have `go` installed and with `GOBIN` set and added to `PATH` you can install the cli using:

```bash
go install github.com/Phillezi/kthcloud-cli@latest
``` 
> [!NOTE]  
> The cli executable will be named `kthcloud-cli` from the module / repo name instead of `kthcloud`.

### Build it yourself

If your OS and architecture combo isnt available as a pre-built binary and you dont want to use `go install` you can build the cli from source.

#### Prerequisites

- [![Git](https://img.shields.io/badge/Git-FFFFFF?style=for-the-badge&logo=Git&logoColor=black)](https://git-scm.com/downloads)
- [![Go >= 1.23.1](https://img.shields.io/badge/Go%20%3E%3D%201.23.1-FFFFFF?style=for-the-badge&logo=go&logoColor=black)](https://go.dev/dl/)
- [![Gnu Make](https://img.shields.io/badge/GNU%20Make-FFFFFF?style=for-the-badge&logo=GNU&logoColor=black)](https://www.gnu.org/software/make/)

1. Clone the repository:

   ```bash
   git clone https://github.com/Phillezi/kthcloud-cli.git
   cd kthcloud-cli
   ```

2. Build the application:

   ```bash
   make
   # or
   # make install
   # make install adds it to the same dir as the install script does, but you will manually have to add it to your PATH.
   ```

3. Run the application:

   ```bash
   ./bin/kthcloud
   ```

> [!TIP]
> Alternatively you can add it to the PATH to be able to use it globally. The installation script automatically does this.

### Commands

#### Login command

Logs in to kthcloud and retrieves an authentication token, the token gets saved to a file named `session.json` inside the configuration path. It opens a browser window to let you login through the kthcloud keycloak login page.

##### Usage of the login command

```bash
kthcloud login
```

#### Compose command

Parses a `docker-compose.yaml`, `docker-compose.yml` or `*.docker-compose.yaml` file (will prefer `kthcloud.docker-compose.yaml/yml`) and gives the ability to bring up these services with the specified configuration on [`kthcloud`](https://cloud.cbh.kth.se).

##### Usage of the compose command

```bash
kthcloud compose # lists all options
```

##### Compose up command

Brings up the services defined in the Docker Compose file.

##### Compose down command

Brings down the services defined in the Docker Compose file.
> ![NOTE]
> This will not remove the volumes created on the storagemanager.

##### Compose parse command

Parses a Docker Compose file and prints the Services, Envs, Ports, Commands, Depends on and Volumes. And prints out the resulting deployments (the json used with the REST API).

#### Update command

Checks for newer releases than the release of the binary running the command. If a newer release is found it will prompt you to install it, (can be bypassed wit the `-y` flag).

Versions can be selected by passing the `-i` flag.

> [!WARNING]
> This currently doesnt work as expected on Windows.

> [!WARNING]
> This does not verify against a hash to confirm the integrity of the bibary (yet).

##### Usage of the update command

```bash
kthcloud update
```

#### Version command

Displays the version of the binary.

##### Usage of the version command

```bash
kthcloud version
```

## Configuration

The `kthcloud-cli` uses a configuration file named `config.yaml` it is located in the configuration directory. You can specify the following fields:

- `api-url`: The URL of the API endpoint.
- `api-token`: The api token from kthcloud.
- `loglevel`: The logging level (info, warn, error, debug) (default "info")
- `resource-cache-duration duration`: How long resources should be cached when possible (default 1m0s)
- `session-path`: The filepath where the session should be loaded and saved to (default "~/.config/.kthcloud/session.json")
- `zone`: The preferred kthcloud zone to use, will use `se-flem2` by default

Example `config.yaml`:

```yaml
api-url: https://api.example.com
api-token: your-api-key-from-kthcloud
loglevel: error
```
