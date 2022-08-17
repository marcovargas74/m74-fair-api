FROM golang:1.17.0-alpine3.14 AS builder

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache build-base git 
 

# Set the working directory
WORKDIR $GOPATH/src/fair-app

COPY . .

# Build and strip the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/fair-app src/fair/main.go

FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/fair-app /go/bin/fair-app

# Run the fair-app binary.
CMD ["/go/bin/fair-app"]


