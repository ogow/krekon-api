FROM golang:1.24.1-bookworm

WORKDIR /krekonapi

COPY ./src/ .
# COPY ./src/go.mod ./src/go.sum ./

RUN go mod download

# COPY ./src/*.go ./

# RUN go build cmd/main.go -o /api

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

EXPOSE 3000

CMD [ "/api" ]
