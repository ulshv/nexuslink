#!/bin/bash

protoc -I=./proto --go_out=./ ./proto/*
