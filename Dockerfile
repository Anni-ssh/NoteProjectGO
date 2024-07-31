FROM golang:1.22-alpine AS builder
LABEL author="AnniSSH"

WORKDIR /user/local/src

RUN apk --no-cache add bash git gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./cmd/NoteProject/main.go

FROM alpine AS runner

RUN apk --no-cache add bash

COPY --from=builder /user/local/src/bin/app /
COPY --from=builder /user/local/src/migrations /migrations

COPY wait-for-it.sh /wait-for-it.sh

CMD ["./wait-for-it.sh", "postgres:5432", "--", "/app"]

EXPOSE 8080

