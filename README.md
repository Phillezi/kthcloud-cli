# kthcloud-cli

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Commands](#commands)
- [License](#license)

## Overview

`kthcloud-cli` is a command-line interface tool for interacting with kthclouds rest api. It allows you to perform various operations such as logging in, fetching resources, and hopefully in the future read docker-compose configurations and deploy them as projects to the cloud.

![Screencast from 2024-09-09 08-23-28](https://github.com/user-attachments/assets/ec040ce9-d11c-436b-9a00-42a0f8de0a1b)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/kthcloud-cli.git
   cd kthcloud-cli
   ```

2. Build the application:

   ```bash
   make
   ```

3. Run the application:
   ```bash
   kthcloud
   ```

## Usage

### Logging In

To log in to kthcloud using keycload, use the `login` command:

```bash
kthcloud login
```

### Parsing Docker Compose Files

To parse a Docker Compose file, use the `compose parse` command:

```bash
kthcloud compose parse
```

### Commands

#### `login`

Logs in to kthcloud and retrieves an authentication token, the token gets saved to the config.yaml.

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
