#!/bin/bash
mkfifo /tmp/fifo-to-badge
mkfifo /tmp/fifo-from-badge

rm -rf admin/build
rm -rf client/build

cd admin/
npm run build
cd ..

docker-compose up --build