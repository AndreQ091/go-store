FROM okteto/golang.1.17 as build

COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs


FROM scratch

COPY --from=build /src/kvs .
COPY --from=build /src/*.pem .

EXPOSE 8080
CMD ["/kvs"]

# CGO_ENABLED=0 GOOS=linux go build -a -o kvs
