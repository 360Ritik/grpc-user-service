# gRPC User Service

This repository contains a gRPC-based user service implementation using Go.

[Watch API Demonstration Video] [grpc-user-service.webm](https://github.com/360Ritik/grpc-user-service/assets/93071300/67c79e05-7172-4ea7-aa14-a38e363b3f99)

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Building and Running](#building-and-running)
  - [Accessing gRPC Endpoints](#accessing-grpc-endpoints)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project implements a user service with the following features:
- Retrieve a user by ID
- Retrieve multiple users by their IDs
- Add a new user
- Search users based on criteria (city, phone, marital status, etc.)

The service is implemented using gRPC in Go and uses Protocol Buffers for defining service contracts.

## Prerequisites

Before running this application, ensure you have the following installed:

- Go (version 1.22.5 or later)
- Protocol Buffers (protobuf) compiler
- grpc-go library

## Install Dependencies
To install dependencies for this project, run the following command:
- go mod download

## Running the Application
To run the application directly on your terminal :
- go run server.go


## Running with Docker
I have tested the application using Docker, ensuring smooth operation.

If you prefer to run the application using Docker, you can build the Docker image and run it.

- docker build -t grpc-user-service:latest .
- docker run -p 50051:50051 grpc-user-service



## Clone the Repository

Clone the repository and navigate into the project directory:

```bash
git clone https://github.com/360Ritik/grpc-user-service.git
cd grpc-user-service
