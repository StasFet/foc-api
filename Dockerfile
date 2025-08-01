# build stage
FROM golang:1.24-alpine AS build-stage

LABEL maintainer="st.fetisov@gmail.com"

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

ENV CGO_ENABLED = 0
ENV GOOS=linux

RUN go build -ldflags="-s -w" -o app .


# Run stage
FROM alpine:3.20 AS run-stage

WORKDIR /app
COPY --from=build-stage /app/app .

CMD ["./app"]