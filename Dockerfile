FROM golang:1.21

ARG USER_UID
ARG USER_GID

RUN groupadd -g $USER_GID go
RUN useradd -m -u $USER_UID -g $USER_GID go

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN mkdir /app && chown -R go:go /app
RUN mkdir /.cache && chown -R go:go /.cache

WORKDIR /app

USER go

CMD ["air"]
