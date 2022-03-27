#!/bin/bash

export mkdir bin && \
  GO111MODULE=on && \
  rm -f ./bin/main && \
  go build -o ./bin/main ./cmd/main.go
