FROM golang:1.12

COPY . /go/src/github.com/wordnet-world/Conductor
WORKDIR /go/src/github.com/wordnet-world/Conductor
RUN echo $GOPATH && \
    make build-linux

FROM scratch
COPY --from=0 /go/src/github.com/wordnet-world/Conductor/Conductor_linux /
CMD ["/Conductor_linux"]
