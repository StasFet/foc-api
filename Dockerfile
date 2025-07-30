FROM golang:1.24-alpine

LABEL maintainer="st.fetisov@gmail.com"

WORKDIR /app

COPY go.* .


RUN go mod download

COPY . .

EXPOSE 8000

ENTRYPOINT [ "sh", "-c", "go build -o app ./src/main.go && ./app" ]

# build : "docker build -t go-foc-api-docker ."

# run : "docker run go-foc-api-docker"