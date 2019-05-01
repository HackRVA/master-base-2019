#!/bin/bash
mkdir -p ~/etc/

echo "serialPort: /dev/ttyACM0
ir: true
serialDebug: true
bwDebug: true
leaderBoard_API: \"http://10.200.200.161:5000/api/\"" > ~/etc/baseconfig.yaml
go run server.go