FROM golang:1.24.2 AS builder

ENV WORKDIR=${GOPATH}/src/pipo

ARG VERSION
ENV VERSION=$VERSION

COPY . ${WORKDIR}
WORKDIR ${WORKDIR}

# RUN go build -ldflags="-X /home/user/sources/geekBrains/backend2/L2/pipo/main.version=${VERSION}" -o ${GOPATH}/bin/pipo ./
# go install -ldflags="-X /home/user/sources/geekBrains/backend2/L2/pipo/version.Version=1.0.a" 
RUN go install -ldflags="-X pipo/version.Version=${VERSION}" 

FROM golang:1.24.2

EXPOSE 8080

# ARG PIPOMASTER
# ENV PIPOMASTER=$PIPOMASTER PIPOMASTERNAME=pipoMaster PIPOSLAVENAME=pipoSlave

RUN mkdir /pipo

WORKDIR /pipo

COPY --from=builder ${GOPATH}/bin/pipo .

ENTRYPOINT ["/pipo/pipo"]