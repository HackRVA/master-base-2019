# Start from golang v1.11 base image
FROM golang:1.11

# Docker Maintainer
LABEL maintainer="Dustin Firebaugh <dafirebaugh@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/HackRVA/master-base-2019

# Copy everything from the current directory
COPY . .

# build main.go
RUN go build -o op-codes opcodes/*.go
RUN go build -o ir irpacket/*.go
RUN go build -o readfifo fiforeader/*.go

# Run the executable
CMD ["./readfifo"]