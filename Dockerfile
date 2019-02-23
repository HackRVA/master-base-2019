# Start from golang v1.11 base image
FROM golang:1.11

# Docker Maintainer
LABEL maintainer="Dustin Firebaugh <dafirebaugh@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/HackRVA/master-base-2019

# Copy everything from the current directory
COPY . .

RUN go get
# build main.go
RUN go build -o server server.go

# Run the executable
CMD ["./server"]