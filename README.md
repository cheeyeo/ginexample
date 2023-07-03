### API + MONGO DB

Example of building a simple API backend that uses mongodb as storage.

It was based on articles:
https://dev.to/deeshath/rest-api-using-gogin-mongo-33i3

https://go.dev/doc/tutorial/web-service-gin


I extended the above by adding:

* Running the service and mongodb in containers

* Created Dockerfile and docker compose template to run both concurrently

* Custom Dockerfile with multistage build to compile the final web artifact

Note that the compiled artifact is a Gin server running in dev mode so as such its not suitable for production use only for learning.


### Running it

To run locally:
```
go mod download
```

To use docker compose:
```
docker compose build web

docker compose up
```

If running the service and mongodb container separately, need to export the mongodb url as an environment variable:

To run mongodb container:
```
docker run --rm --name some-mongo -v ./data:/data/db mongo
```

```
MONGO_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' some-mongo)

export MONGO_URL=mongodb://${MONGO_IP}:27017
```

### How to interact with API

The api is a simple CRUD service for books.

Use either curl or postman to make requests.

To create a new book:
```
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"Title": "Atomic Habits","Author": "James Clear"}'
```

To view a book with id:
```
curl http://localhost:8080/books/64a1be64fb71ba1f1a00b5d5
```

To view all books:
```
curl http://localhost:8080/books
```

To update book with id:
```
curl http://localhost:8080/books/64a3249454fbd70d922ba3f8 \
    --include \
    --header "Content-Type: application/json" \
    --request "PATCH" \
    --data '{"Title": "Mickey 7","Author": "Ashton Edward"}'
```

To delete book with id:
```
curl http://localhost:8080/books/64a1be64fb71ba1f1a00b5d5 \
    --request "DELETE"
```