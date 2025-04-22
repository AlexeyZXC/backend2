FROM golang:1.24.2 AS builder

ENV WORKDIR=${GOPATH}/src/handlers

ARG VERSION
ENV VERSION=$VERSION

COPY . ${WORKDIR}
WORKDIR ${WORKDIR}

# copies build to ${GOPATH}/bin
RUN go install -ldflags="-X main.VERSION=${VERSION}" 

FROM golang:1.24.2

EXPOSE 8080

RUN mkdir /handlers

WORKDIR /handlers

COPY --from=builder ${GOPATH}/bin/handlers .

ENTRYPOINT ["/handlers/handlers"]