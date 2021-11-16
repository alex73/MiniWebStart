FROM golang:1.15 AS build
WORKDIR /src
COPY . .
RUN go get github.com/pbnjay/memory
RUN GOOS=windows GOARCH=386   go build -o /out/mini-win32.exe -ldflags "-s -w" ./src
RUN GOOS=windows GOARCH=amd64 go build -o /out/mini-win64.exe -ldflags "-s -w" ./src
RUN GOOS=darwin  GOARCH=amd64 go build -o /out/mini-mac64     -ldflags "-s -w" ./src
RUN GOOS=linux   GOARCH=amd64 go build -o /out/mini-linux64   -ldflags "-s -w" ./src
RUN GOOS=linux   GOARCH=386   go build -o /out/mini-linux32   -ldflags "-s -w" ./src

FROM scratch AS bin
COPY --from=build /out/mini-* /
