# Two-stage build:
#    first  FROM prepares a binary file in full environment ~780MB
#    second FROM takes only binary file ~10MB

FROM golang:1.16-alpine AS builder

WORKDIR "/go/src/github.com/your-login/your-project"

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /your-app

CMD ["/your-app"]




#########
# second stage to obtain a very small image
FROM scratch

COPY --from=builder /your-app .

CMD ["/your-app"]
