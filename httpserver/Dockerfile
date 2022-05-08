FROM golang:1.18 as httpserver-build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o app httpserver.go

FROM busybox

COPY --from=httpserver-build /usr/src/app/app /usr/local/bin/app
EXPOSE 8099

ENTRYPOINT ["app"]
