FROM golang:1.11

ARG binary
ADD ${binary} /${binary}

RUN touch /tmp/healthy
EXPOSE 54/udp