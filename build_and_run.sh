#!/bin/bash
rm test.log
touch test.log
go build -o basestation server.go
./basestation
