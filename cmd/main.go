package main

import (
    "os"
    "os/signal"
    "time"
    "context"

    "golang.org/x/sys/unix"

    "github.com/rchnmy/xupery/log"
    "github.com/rchnmy/xupery/pkg/handler"
    "github.com/rchnmy/xupery/pkg/controller"
)

func main() {
    ctx, done := signal.NotifyContext(context.Background(),
        unix.SIGINT,
        unix.SIGHUP,
        unix.SIGQUIT,
        unix.SIGTERM,
    )

    db := controller.Connect("postgres", os.Getenv("DB_DSN"))
    defer db.Close()

    cl := controller.NewKafkaClient()
    defer cl.Close()

    c := controller.NewController(db, cl)
    go c.Consume(ctx)

    t := handler.Token{}
    t.Encode(os.Getenv("TOKEN_BASE"))

    r := handler.NewRouter()
    r.Use(t.Validate, handler.JsonHandler)
    r.HandleFunc("/stat", handler.GetStatistic(c)).Methods("GET")
    r.HandleFunc("/send", handler.PostMessage(c)).Methods("POST")

    s := handler.NewServer(":9900", r)
    go handler.Serve(s)

    <- ctx.Done()
    ct, do := context.WithTimeout(context.Background(), 5 * time.Second)
    if err := s.Shutdown(ct); err != nil {
        log.Error().Err(err).Send()
    }
    do()
    done()
}

