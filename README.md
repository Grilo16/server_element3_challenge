# Go Server App

## Introduction

This is a simple Go server application that connects to a SQL Server database running in a Docker container. It provides a basic API that interacts with the database, allows for the upload and download of files and handles authentication through JWT tokens.

## Prerequisites

Before running this application, make sure you have the following installed on your system:

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Go Programming Language: [Install Go](https://golang.org/doc/install)

## Setup

1. Clone this repository to your local machine:

    ```bash
    git clone https://github.com/Grilo16/server_element3_challenge
    ```

2. Navigate to the project directory:

    ```bash
    cd server_element3_challenge
    ```

3. Run the setup script to build the Docker image, create the network and start the server:

    ```bash
    ./setup.sh
    ```

