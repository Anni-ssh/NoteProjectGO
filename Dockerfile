FROM golang:1.22-alpine AS byilder
LABEL authors="AnniSSH"

WORKDIR /user/local/src

RUN apk --no-cache add bash git gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./cmd/NoteProject/main.go

FROM alpine AS runner

COPY --from=byilder /user/local/src/bin/app /
COPY config/config.yaml /config.yaml

CMD ["/app"]

