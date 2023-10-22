## Translation API

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

https://github.com/deepmap/oapi-codegen/tree/master/examples/petstore-expanded/fiber


#### Migrations

`docker run --rm -it --network=host -v "$(pwd)/db:/db" ghcr.io/amacneil/dbmate new initial_schema`

`docker run --rm -it --network=host -v "$(pwd)/db:/db" ghcr.io/amacneil/dbmate up`