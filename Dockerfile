FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

LABEL maintainer="st.fetisov@gmail.com"

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

EXPOSE 8000

ENTRYPOINT [ "sh", "-c", "go build -o app ./foc_api.go && ./app" ]