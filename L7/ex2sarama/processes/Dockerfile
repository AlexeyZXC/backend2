FROM golang:1.24.2 AS builder

# ENV WORKDIR=${GOPATH}/src/processes
ENV WORKDIR=${GOPATH}/src

ARG VERSION
ENV VERSION=$VERSION

COPY handlers ${WORKDIR}/handlers
COPY processes ${WORKDIR}/processes

WORKDIR ${WORKDIR}/processes

# RUN mkdir handlers
# COPY ../handlers ../handlers
# COPY ../handlers/* handlers/
# ${GOPATH}/src/processes

# WORKDIR ${WORKDIR}

# copies build to ${GOPATH}/bin
RUN go install -ldflags="-X main.VERSION=${VERSION}" 

FROM golang:1.24.2

RUN mkdir /processes

WORKDIR /processes

COPY --from=builder ${GOPATH}/bin/processes .

ENTRYPOINT ["/processes/processes"]