FROM golang:1.11 AS build-env


ENV GO111MODULE=on
ADD . /src

# Compile the binary, we don't want to run the cgo resolver
RUN cd /src && GOOS=linux GOARCH=386 go build -o goapp

# final stage
FROM alpine

WORKDIR /

ENV PORT=$PORT
WORKDIR /app
COPY --from=build-env /src/goapp /app/
CMD ["./goapp"]