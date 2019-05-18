#!/bin/bash
mkdir -p /etc/basestation/

echo "serialPort: /dev/ttyACM0
ir: true
serialDebug: true
bwDebug: true" > /etc/basestation/baseconfig.yaml