FROM golang:1.12 as base
RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.3/seabolt-1.7.3-Linux-ubuntu-18.04.deb && \
    dpkg -i seabolt-1.7.3-Linux-ubuntu-18.04.deb

FROM base
COPY . /go/src/github.com/wordnet-world/Conductor
WORKDIR /go/src/github.com/wordnet-world/Conductor
RUN echo $GOPATH && \
    make
CMD ["/go/src/github.com/wordnet-world/Conductor/Conductor"]
