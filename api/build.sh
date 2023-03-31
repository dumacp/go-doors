#!/bin/bash

protoc -I=./ -I=$GOROOT/src --go_out=. --go_opt=paths=source_relative ./messages.proto
