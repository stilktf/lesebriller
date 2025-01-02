FROM golang:1.23.2 AS build

WORKDIR /src
COPY . /src

RUN go mod download
RUN go test -v
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w" -o /src/lesebriller

FROM alpine AS bin

WORKDIR /app
COPY --from=build /src/lesebriller .

EXPOSE 7200

CMD ["/app/lesebriller"]