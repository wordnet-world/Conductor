FROM golang:1.12

COPY . /go/src/github.com/wordnet-world/Conductor
WORKDIR /go/src/github.com/wordnet-world/Conductor
RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install
RUN echo $GOPATH && \
    make

FROM scratch
COPY --from=0 /go/src/github.com/wordnet-world/Conductor/Conductor /
COPY ./config/conductor-conf.json /config/conductor-conf.json
CMD ["/Conductor"]
