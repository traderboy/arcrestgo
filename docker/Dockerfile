#FROM postgres:9.5
#FROM golang:onbuild
#$ docker build –t traderboy/arcrestgo -f docker/Dockerfile .
#$  docker run -d -p 8000:8000  --name arcrestgo traderboy/arcrestgo
#   docker logs arcrestgo
#   docker rm arcrestgo

FROM golang:1.7
#RUN mkdir /app
#RUN mkdir /usr/local/go/src/github.com/traderboy
RUN mkdir -p /usr/local/go/src/github.com/traderboy/arcrestgo
RUN go get github.com/gorilla/handlers
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/twinj/uuid

ADD docker/server.go /usr/local/go/src/github.com/traderboy/arcrestgo 
ADD structs /usr/local/go/src/github.com/traderboy/arcrestgo/structs
ADD routes /usr/local/go/src/github.com/traderboy/arcrestgo/routes
ADD config /usr/local/go/src/github.com/traderboy/arcrestgo/config
ADD docker/config.go /usr/local/go/src/github.com/traderboy/arcrestgo/config/
ADD catalogs /usr/local/go/src/github.com/traderboy/arcrestgo/catalogs
ADD dist /usr/local/go/src/github.com/traderboy/arcrestgo/dist
ADD sha256 /usr/local/go/src/github.com/traderboy/arcrestgo/sha256
ADD arcrest.sqlite /usr/local/go/src/github.com/traderboy/arcrestgo/
RUN ls -l /usr/local/go/src/github.com/traderboy/arcrestgo/

WORKDIR /usr/local/go/src/github.com/traderboy/arcrestgo
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install github.com/traderboy/arcrestgo
RUN go build /usr/local/go/src/github.com/traderboy/arcrestgo/server.go

CMD ["./server"]
#ADD myapp /app/
#ENTRYPOINT ["./server"]

EXPOSE 8080