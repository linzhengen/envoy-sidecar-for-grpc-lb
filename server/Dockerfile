FROM golang:1.13-alpine as build

WORKDIR /src

ADD . /src/

RUN GOPATH="/go" GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -installsuffix "static" && \
    echo 'nobody:x:65534:' > /src/group.nobody && \
    echo 'nobody:x:65534:65534::/:' > /src/passwd.nobody

FROM gcr.io/distroless/static

COPY --from=build /src/group.nobody /etc/group
COPY --from=build /src/passwd.nobody /etc/passwd
USER nobody:nobody

COPY --from=build /go/bin/server /bin/server

ENTRYPOINT ["/bin/server"]
