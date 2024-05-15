# Providers Endpoints

This repository contains the source code for a Go application that interacts with various cloud service providers.

## Structure

The main application code is in `main.go`. The `service/` directory contains individual Go files for each service provider, such as AWS, DigitalOcean, and Google Cloud. Each of these files contains functions for interacting with the respective service.

## Dependencies

This project uses several dependencies, including:

- `github.com/PuerkitoBio/goquery` for parsing HTML
- `github.com/go-ping/ping` for pinging servers
- `github.com/tbxark/g4vercel` for Vercel integration

## Building and Running

To build and run this project, you need to have Go installed. Then, you can use the `go build` and `go run` commands.
