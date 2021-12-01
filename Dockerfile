# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.16 as builder

# Add Maintainer Info
LABEL maintainer="wade <wadejet.work@gmail.com>"

ARG GO_OPTS
# Set the Current Working Directory inside the container
WORKDIR /build

#新增生成swag文檔的套件
RUN go get -u github.com/swaggo/swag/cmd/swag

# Copy go mod and sum files
COPY . .
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
#COPY . .

#生成 swag 文檔
RUN swag init

# Build the Go app
# RUN go test -v ./repositories/...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
#RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:3.14

WORKDIR /app/

RUN echo "http://alpine.ccns.ncku.edu.tw/alpine/v3.14/main/" > /etc/apk/repositories && apk update && apk add ca-certificates && apk add curl

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/main ./
COPY --from=builder /build/resources ./resources

RUN ls && cd ./resources && ls && cd ../

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ./main ${GO_OPTS}
