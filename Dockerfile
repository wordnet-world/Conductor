FROM golang:1.12

COPY . /root/go/src/github.com/wordnet-world/Conductor
WORKDIR /root/go/src/github.com/wordnet-world/Conductor
RUN echo $GOPATH && \
    make build-linux

FROM scratch
COPY --from=0 /root/go/src/github.com/wordnet-world/Conductor/Conductor_linux /
CMD ["/Conductor_linux"]
