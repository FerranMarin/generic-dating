FROM golang:1.22-alpine

ENV PORT 3000
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD ["go", "run", "."]
EXPOSE ${PORT}