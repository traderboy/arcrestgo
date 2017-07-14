#FROM postgres:9.5
#kartoza/postgis:9.5-2.2
#FROM golang:onbuild

#   docker stop arcrestgo;docker rm arcrestgo;docker build -t traderboy/arcrestgo -f docker/Dockerfile .;docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 -e ROOT_PATH=catalogs -e DB_SOURCE=PG -e DB_NAME="user=postgres dbname=gis host=172.17.0.5" --name arcrestgo traderboy/arcrestgo;docker logs arcrestgo -f
#   docker build –t traderboy/arcrestgo -f docker/Dockerfile .
#   docker rm arcrestgo
#   docker run -d -p 80:80 -p 443:443 -e HTTPS_PORT=443 -e HTTP_PORT=80 --name arcrestgo traderboy/arcrestgo
#   docker run -d -p 80:8000 --name arcrestgo --link db:postgres traderboy/arcrestgo
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


ADD structs /usr/local/go/src/github.com/traderboy/arcrestgo/structs
ADD routes /usr/local/go/src/github.com/traderboy/arcrestgo/routes
ADD config /usr/local/go/src/github.com/traderboy/arcrestgo/config
ADD catalogs /usr/local/go/src/github.com/traderboy/arcrestgo/catalogs
ADD dist /usr/local/go/src/github.com/traderboy/arcrestgo/dist
ADD ssl /usr/local/go/src/github.com/traderboy/arcrestgo/ssl
ADD sha256 /usr/local/go/src/github.com/traderboy/arcrestgo/sha256
ADD arcrest.sqlite /usr/local/go/src/github.com/traderboy/arcrestgo/

ADD server.go /usr/local/go/src/github.com/traderboy/arcrestgo 
#ADD docker/server.go /usr/local/go/src/github.com/traderboy/arcrestgo 
#ADD docker/config.go /usr/local/go/src/github.com/traderboy/arcrestgo/config
ADD index.html /usr/local/go/src/github.com/traderboy/arcrestgo

RUN chmod 775 /usr/local/go/src/github.com/traderboy/arcrestgo/index.html
RUN chmod 775 /usr/local/go/src/github.com/traderboy/arcrestgo/arcrest.sqlite
RUN chmod -R 775 /usr/local/go/src/github.com/traderboy/arcrestgo/dist

#RUN ls -l /usr/local/go/src/github.com/traderboy/arcrestgo
#RUN ls -l /usr/local/go/src/github.com/traderboy/arcrestgo/config

WORKDIR /usr/local/go/src/github.com/traderboy/arcrestgo
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install github.com/traderboy/arcrestgo
RUN go build /usr/local/go/src/github.com/traderboy/arcrestgo/server.go

#CMD ["./server","-root catalogs","-file"]
#CMD ["./server","-root","catalogs","-sqlite","arcrest.sqlite"]
#CMD ["./server","-root","catalogs","-pgsql","user=postgres dbname=gis host=172.17.0.5"]
CMD ["./server"]

#ADD myapp /app/
#ENTRYPOINT ["./server"]

#EXPOSE 8080
