FROM golang:1.24.2 AS builder

ENV WORKDIR=${GOPATH}/src/acl

COPY . ${WORKDIR}
WORKDIR ${WORKDIR}/acl

# copies build to ${GOPATH}/bin
RUN go install

FROM golang:1.24.2

EXPOSE 80

RUN mkdir /acl

WORKDIR /acl

COPY --from=builder ${GOPATH}/bin/acl .

ENTRYPOINT ["/acl/acl"]