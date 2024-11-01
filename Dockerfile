FROM golang:latest

WORKDIR /usr/src/app
COPY . .
VOLUME [ ".:/usr/src/app" ]
CMD [ "go", "run", "cmd/main.go", "-b", "0.0.0.0" ]
RUN  go mod tidy
EXPOSE  8181:8181
