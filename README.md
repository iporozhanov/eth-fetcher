# eth-fetcher


## Table of Contents

- [Usage](#usage)
  - [Running Locally](#running-locally)
  - [Running with Docker](#running-with-docker)
  - [Run tests](#run-tests)
  - [Authenticate](#authenticate)
- [ApiDoc](openapi.yaml)
## Usage

### Running Locally
1. Setup .env file

   ```shell
   cp .env.example .env
   ```
And fill it with the correct settings

2. Run it

   ```shell
   make run
   ```

### Running with Docker

1. Build the Docker image:

   ```shell
   make build-docker
   ```

2. Run the Docker container with docker compose:

   ```shell
   make run-docker
   ```

### Run tests

   ```shell
   make test
   ```

### Authenticate
The following usernames and passwords are available
- `alice`/ `alice`
- `bob`/ `bob`
- `carol` / `carol`
- `dave`/ `dave`