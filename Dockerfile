FROM golang:1.17.1-alpine3.14 as builder
# Run update and install git
RUN apk update && apk add --no-cache git

# Set Workdir inside container
WORKDIR /go/src/app

#Copy Files into container
COPY . .

# Get Dependency
RUN go mod download

# Install Air for hot reload
RUN go get -u github.com/cosmtrek/air

# The ENTRYPOINT defines the command that will be ran when the container starts up
ENTRYPOINT air