FROM golang:1.21

ARG USER_UID
ARG USER_GID

RUN groupadd -g $USER_GID go
RUN useradd -m -u $USER_UID -g $USER_GID go

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 $GOPATH/bin/migrate

RUN mkdir /app && chown -R go:go /app
RUN mkdir /.cache && chown -R go:go /.cache

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod download

USER go

CMD ["air"]
