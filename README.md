# NoteProject

## Структура проекта

- cmd/
    - noteProject/
        - main.go

- config/
    - config.yaml

- internal/
    - config/
        - config.go
    - entities/
        - note.go
        - user.go
    - service/
        - authorization.go
        - service.go
    - storage/
        - postgres/
          - authorization.go
          - error.go
          - postgres.go
        - redis/
        - storage.go
    - transport/
        - http-server/
          - handler/
              createNote.go
              deleteNote.go
              handler.go
              home.go
              middleware.go
              readNotes.go
              response.go
              saveNote.go
              sign_in.go
              sign_up.go
              static.go
              updateNote.go
          - server/
              server.go
- pkg/
    logger/
          - slog.go
- schema/
    - 000001_init.down.sql
    - 000001_init.up.sql
- web/
- .gitignore
- go.mod
- go.sum
- README.md
