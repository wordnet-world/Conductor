FROM golang:1.12

COPY . /root/Conductor
WORKDIR /root/Conductor
ENV GOPATH=/root/Conductor
RUN echo $GOPATH && \
    make build-linux

FROM scratch
COPY --from=0 Conductor_linux /
CMD ["/Conductor_linux"]
