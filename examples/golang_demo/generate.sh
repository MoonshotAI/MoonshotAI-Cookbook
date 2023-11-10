#!/bin/bash

URL="https://api.moonshot.cn/v1"
KEY="MOONSHOT_API_KEY"

GOPACKAGE=main MOONSHOT_BASE_URL=$URL MOONSHOT_API_KEY=$KEY go run -tags=generate_models_file . -file=constants.gen.go
go generate client.go
go mod tidy