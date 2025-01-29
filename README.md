# tcp-chat

## Overview

Simple TCP-based chat.

It allows:
  - create a server
  - connect to a server
  - [TODO] register on the server
  - [TODO] login on the server
  - [TODO] create a room (later: private rooms with a password)
  - [TODO] list rooms (including private ones)
  - [TODO] join room (and entering the password for the private ones)
  - [TODO] send messages in the room
  - [TODO] exit the room
  - [TODO] disconnect from the server
  - [TODO] [OPTIONAL] DM functionality

## Build

You need to have `protoc-gen-go` module to be installed:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Build protobuf files into `./internal/proto`:
```bash
protoc -I=./proto --go_out=./internal/proto ./proto/*
```

# Usage

Start the server:
```bash
./run
```

