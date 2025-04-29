FROM golang:1.24.2 AS builder

ENV WORKDIR=${GOPATH}/src/serviceRH

COPY . ${WORKDIR}
WORKDIR ${WORKDIR}

# copies build to ${GOPATH}/bin
# RUN go install -ldflags="-X main.VERSION=${VERSION}" 
RUN go install

FROM golang:1.24.2

EXPOSE 80

RUN mkdir /serviceRH

WORKDIR /serviceRH

COPY --from=builder ${GOPATH}/bin/serviceReqHandler .

ENTRYPOINT ["/serviceRH/serviceReqHandler"]