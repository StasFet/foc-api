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
docker-compose up --build --watch
``` 
Or alternatively use the `start.sh` script which will also prune dangling images for your convenience.
```bash
sh ./start.sh
```
Enjoy!
## API Map

! UNDER CONSTRUCTION !