# kthcloud-cli

> [!NOTE]  
> This project is in the very early stages of development. Features are incomplete, and things may change frequently.

## Table of Contents

- [Overview](#overview)
    - [Compose](#compose)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Commands](#commands)

## Overview

`kthcloud-cli` is a command-line interface tool for interacting with kthclouds rest api. It allows you to perform various operations such as logging in, fetching resources, and creating deployments from `docker-compose` files.

![Screencast-from-2024-09-10-17-47-54](https://github.com/user-attachments/assets/ffa9d85d-0974-4a89-a480-3918b4ebb35f)

### Compose

The clis core functionallity is to parse `docker compose` files and create deployments from the content.

For example, if i have this `docker-compose.yaml` file:
```yaml
services:
  testingcompose1:
    image: registry.cloud.cbh.kth.se/waitapp/cicd:latest
    ports:
      - "8080:8080"
  testingcompose:
    image: postgres:15
    environment:
      POSTGRES_USER: supersecretuserhere
      POSTGRES_PASSWORD: supersecretpassword
      POSTGRES_DB: WAIT
    command: ["sleep", "infinity"]
    ports:
      - "5432:5432"
    volumes:
      - dbdata:/var/lib/postgresql/data
```
> [!NOTE] 
> The above example is just an example to showcase what is supported. It does not provide a functional application, the database will just run sleep

The tool will create two deployments and set up their environment variables, port, start commands and persistent storage.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Phillezi/kthcloud-cli.git
   cd kthcloud-cli
   ```

2. Build the application:

   ```bash
   make
   ```

3. Run the application:
   ```bash
   ./bin/kthcloud
   ```
> [!TIP]
> Alternatively you can add it to the PATH to be able to use it globally, TODO: Make static .config location, the current configuration and session files will break otherwise.

<!--
   Alternatively you can move it to a location and add it to the path to be able to use it globally. **(dont do it yet, not ready)**
   ```bash
   sudo cp ./bin/kthcloud /usr/local/bin/
   ```
   Also make sure that `/usr/local/bin` is on the `PATH`.
   ```bash
   echo $PATH | grep /usr/local/bin
   ```
-->

## Usage

### Logging In

To log in to kthcloud using keycload, use the `login` command:

```bash
kthcloud login
```

This will bring up your browser and prompt you to login.

### Compose

To parse a Docker Compose file, and deploy to the cloud. Use the `compose up` command:

```bash
kthcloud compose parse
```

### Commands

#### `login`

Logs in to kthcloud and retrieves an authentication token, the token gets saved to a file named session.json.

**Usage:**

```bash
kthcloud login
```

#### `compose`

Parses a `docker-compose.yaml` or `docker-compose.yml` file and displays the services, environment variables, ports, and volumes.

**Usage:**

```bash
kthcloud compose
```

**Sub-Commands:**

- `parse`: Parses a Docker Compose file and prints the Services, Envs, Ports and Volumes.
- `up`: (TODO)Brings up the services defined in the Docker Compose file.
- `down`: (TODO)Brings down the services defined in the Docker Compose file.

## Configuration

The `kthcloud-cli` uses a configuration file named `config.yaml`. You can specify the following fields:

- `api-url`: The URL of the API endpoint.
- `auth-token`: The authentication token from keycloak to access the API.

Example `config.yaml`:

```yaml
api-url: https://api.example.com
auth-token: your-auth-token-from-keycloak
```
