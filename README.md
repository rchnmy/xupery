`xupery` is a mock REST API integrated with PostgreSQL and Kafka.

![Untitled Diagram drawio1](https://github.com/user-attachments/assets/8b304052-c753-41ef-ba39-7b968ef3b7f3)


### Build
```go
docker-compose build && \
docker-compose pull postgres kafka && \
docker-compose up xupery -d
```
### Usage
Find `X-Upery-Token` in the first log message of a `xupery` container.
```go
docker logs --tail 1 xupery
{"level":"info","time":"2024-07-29T16:43:21Z","message":"X-Upery-Token: 🦊"}
```
Use it to address the `/send` and `/stat` endpoints.
```go
curl -H 'X-Upery-Token: 🦊' http://localhost:8080/send -d '{"text": "A goal without a plan is just a wish."}'
{"topic":"messages", "record":"A goal without a plan is just a wish."}
```

```go
curl -H 'X-Upery-Token: 🦊' http://localhost:8080/stat
{"messages_total":1}
```
