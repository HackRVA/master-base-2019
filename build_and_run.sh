#!/bin/bash
mkfifo /tmp/fifo-to-badge
mkfifo /tmp/fifo-from-badge

docker-compose up --build