#!/bin/bash
rm -rf /tmp/fifo-to-badge
rm -rf /tmp/fifo-from-badge

mkfifo /tmp/fifo-to-badge
mkfifo /tmp/fifo-from-badge

go run server.go