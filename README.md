# FOC API

### Introduction
Hello! This is just a little project to learn how to make a REST API. I haven't made one before, but I thought it seems interesting and useful so here I am. The Festival of Creativity (FOC) is an annual event held at my school during which students have the chance to form bands and perform songs or dances to the school community. Here is a little list of the tools/techniques used for this project:

- REST API design
- Golang
- SQLite3
- Golang standard library's `net/http` module (I initially wanted to use gin, but decided to use this to make a more lightweight app)
- Docker/Docker-Compose

### Running This Application
If for whatever reason you decide to run this application, you can easily do so with Docker Desktop. Simply open the root directory of this project in your terminal and run. 
```bash
docker-compose up --build
``` 
Or alternatively use the `start.sh` script which will also prune dangling images for your convenience. Ensure you have docker-compose version >=2.22.0, or remove the `--watch` flag will not work.
```bash
sh ./start.sh
```
Enjoy!
## API Map

| Path                       | Description                          |
| -------------------------- | ------------------------------------ |
| `GET /performers`          | Returns all the performers           |
| `GET /performers/:id`      | Returns the performer with id `id`   |
| `GET /performances`        | Returns all the performances         |
| `GET /performances/:id`    | Returns the performance with id `id` |
| `POST /performers`         | Creates a new performer              |
| `POST /performances`       | Creates a new performance            |
| `PUT /performers/:id`      | Updates the performer with id `id`   |
| `PUT /performances/:id`    | Updates the performance with id `id` |
| `DELETE /performers/:id`   | Deletes the performance with id `id` |
| `DELETE /performances/:id` | Deletes the performance with id `id` |