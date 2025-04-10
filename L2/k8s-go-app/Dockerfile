ARG GIT_COMMIT
ARG VERSION
ARG PROJECT

FROM golang:1.24.2 AS builder
ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

ARG VERSION
ENV VERSION=$VERSION

ARG PROJECT
ENV PROJECT=$PROJECT

ENV GOSUMDB=off
ENV GO111MODULE=on
ENV WORKDIR=${GOPATH}/src/k8s-go-app

COPY . ${WORKDIR}
WORKDIR ${WORKDIR}

RUN set -xe ;\
    go install -ldflags="-X ${PROJECT}/version.Version=${VERSION} -X ${PROJECT}/version.Commit=${GIT_COMMIT}" ;\
    ls -lhtr /go/bin/

FROM golang:1.24.2

EXPOSE 8080

WORKDIR /go/bin

COPY --from=builder /go/bin/k8s-go-app .
COPY --from=builder ${GOPATH}/src/k8s-go-app/config/*.env ./config/

ENTRYPOINT ["/go/bin/k8s-go-app"]