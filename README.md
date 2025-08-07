# FOC API

### Introduction
Hello! This is just a little project to learn how to make a REST API. I haven't made one before, but I thought it seems interesting and useful so here I am. The Festival of Creativity (FOC) is an annual event held at my school during which students have the chance to form bands and perform songs or dances to the school community. Here is a little list of the tools/techniques used for this project:

- REST API design
- Golang
- SQLite (from modernc.org)
- Golang standard library's `net/http` module (I initially wanted to use gin, but decided to use this to make a more lightweight app)
- Docker/Docker-Compose
- Complete unit testing using Go's `testing` tool and `testify`

### Running This Application
If for whatever reason you decide to run this application, you can easily do so with Docker Desktop. Simply open the root directory of this project in your terminal and run the following command:
```bash
docker compose up --build
```
Enjoy!

### Running Tests
To run the unit tests and ensure everything works, you can run the following command:
```bash
go test ./internal/ -v 
```

## API Map

| Path                       | Description                          |
| -------------------------- | ------------------------------------ |
| `GET /performers`          | Returns all the performers           |
| `GET /performers/:id`      | Returns the performer with id `id`   |
| `GET /performers/:id/performances` | Returns the performances of performer with id `id` |
| `GET /performances`        | Returns all the performances         |
| `GET /performances/:id`    | Returns the performance with id `id` |
| `GET /performances/:id/performers` | Returns the performers of performance with id `id` |
| `POST /performers`         | Creates a new performer              |
| `POST /performances`       | Creates a new performance            |
| `POST /junctions`          | Creates a performer:performance pair |
| `PUT /performers/:id`      | Updates the performer with id `id`   |
| `PUT /performances/:id`    | Updates the performance with id `id` |
| `DELETE /performers/:id`   | Deletes the performance with id `id` |
| `DELETE /performances/:id` | Deletes the performance with id `id` |
| `DELETE /junctions/:id1/:id2` | Deletes the performer:performance pair with ids `id1:id2` |