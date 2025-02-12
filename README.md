# NexusLink

NexusLink is a project with an attempt to create yet another p2p/decentralized app/system stack
with well-known technologies (like `Golang` with its stdlib), as well as experimental tech (i.e. `yggdrasil`).
Unlike some existing admired p2p stacks (like `libp2p`), the project's goal is to be as simple as possible (KISS),
which includes extensive comments and learning materials (i.e. tutorials on how some of the parts are made),
while not compromizing security (i.e. with strong encryption).

First app implementations of the NexusLink would be
the following projects (all made with strong pub-key encryption):
- p2p TCP chat
- p2p TCP messenger
- p2p TCP "hole punching"
- p2p TCP file sharing

Later app implementations could be something like:
- p2p TCP audio/video calls
- p2p TCP tunneling
- p2p Decentralized DNS for pub keys

Or even something like:
- p2p custom AI Assistants/Agents running on user's computer and data (via something like `ollama` with a knowledge RAG)
  and their communication in p2p-manner with another AIs and users
- p2p decentralized blockchain networks, building DeFi platforms to serve local communities
- the list goes on

Yes, I know, there's existing tech for all of that but this world deserves diversity in tech and ideas!
And mainly, I'm doing all of this for myself in the first place (with learning/practice goal in mind)
rather to raise a lot of $$$ in funding or to implement features needed for some specific companies.

# Later

Nodes (clients as well as servers) later could be identified not only by `<ip/domain>:<port>`,
but also by 256-bit hash of their public keys
in a Decentralized hash-based DNS system.

# Why?

Initially my dream was to build a new internet (as they called it in Silicon Valley HBO)
built on a decetralized LoRa-based radio mesh network (like `Reticulum`).
But then I realized that there's almost zero radio enthusiast in my local area,
so I felt that I would be the only guy in such a network...

So then I came up with an idea - why no use the existing networking tech (like the Internet)
to build a network of inter-connected devices? It's when I came up with the idea of a decentralized p2p `tcp-chat`.

As the NexusLink project evolves, my goal is to continue adding the packages/software in this monorepo
for creation of decentralized p2p apps.

# TODO/WIP/DONE

## tcp-chat
- [x] initial PoC TCP setup / networking code
- [ ] finalize the TCP-chat functionality without p2p parts
  - [ ] create/join/leave rooms and send messages
  - [ ] messages history
  - [ ] authn and authz
- [ ] make tcp tunneling and allow to run something like:
      `./build/server -p 5000 --domain=tcp-chat-1.sergeycooper.com`
      or tunnel the local TCP port on a public server/domain
- [ ] create a basic PoC p2p functionality and allow clients to become chat servers


# NexusLink Projects

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

## Usage

Start the server:
```bash
./run
```

