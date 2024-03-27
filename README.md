# go-mongo-microservice

### MongoDB setup 

docker-compose.yml
```
version: '3'
services:
  mongodb:
    image: mongo
    container_name: catalogdb
    ports:
      - "27017:27017"
    volumes:
      - "mongodata:/data/db"
    networks:
      - network1

volumes:
   mongodata:

networks:
   network1:
```
```
$ docker-compose up
```

### verify network
```
$ docker network ls
NETWORK ID          NAME                DRIVER              SCOPE
d9772c6b6972        bridge              bridge              local
cf3d62f44f7c        host                host                local
b90b47a74a9d        mongodb_network1    bridge              local
cc1c1ed12914        none                null                local

```

### build golang microservice image

Dockerfile
```
FROM golang:1.9.2 as builder
ARG SOURCE_LOCATION=/
WORKDIR ${SOURCE_LOCATION}
RUN go get -d -v github.com/gorilla/mux \
	&& go get -d -v gopkg.in/mgo.v2/bson \
	&& go get -d -v gopkg.in/mgo.v2
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
ARG SOURCE_LOCATION=/
RUN apk --no-cache add curl
EXPOSE 9090
WORKDIR /root/
COPY --from=builder ${SOURCE_LOCATION} .
CMD ["./app"]  

```

```
$ docker build --build-arg SOURCE_LOCATION=<GO_FILE_DIR> --no-cache -t go-mongo-microservice:latest .
$ docker images
REPOSITORY              TAG                 IMAGE ID            CREATED             SIZE
go-mongo-microservice   latest              51246ed02f92        29 minutes ago      13MB
mongo                   latest              43099507792a        38 hours ago        366MB
mysql                   latest              5d4d51c57ea8        38 hours ago        374MB
wordpress               latest              80a6fca6cc6a        10 days ago         407MB
golang                  latest              1c1309ff8e0d        10 days ago         779MB
ubuntu                  latest              0458a4468cbc        4 weeks ago         112MB
alpine                  latest              3fd9065eaf02        7 weeks ago         4.15MB
golang                  1.9.2               138bd936fa29        2 months ago        733MB
```

### run docker microservice image


```
$ docker run --name go-mongo-microservice --link catalogdb:mongo  --net mongodb_network1 -p 9090:9090  go-mongo-microservice
```

```
$ curl -v -d '{"name":"iPhone","company":"Apple"}'  -X POST http://localhost:9090/catalogs
$ curl -v -d '{"name":"Note","company":"Samsung"}'  -X POST http://localhost:9090/catalogs
```

```
$ curl http://localhost:9090/catalogs | json_pp

[
   {
      "name" : "iPhone",
      "company" : "Apple"
   },
   {
      "company" : "Samsung",
      "name" : "Note"
   }
]
```
