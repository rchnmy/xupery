package controller

import (
    "fmt"
    "errors"
    "context"
    "database/sql"

    _ "github.com/lib/pq"
    "github.com/twmb/franz-go/pkg/kgo"

    "github.com/rchnmy/xupery/log"
)

const (
    broker = "kafka:9092"
    topic  = "messages"
)

func Connect(driver, dsn string) *sql.DB {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        log.Error().Err(err).Send()
    }
    if err := db.Ping(); err != nil {
        log.Fatal().Err(err).Send()
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id SERIAL PRIMARY KEY, text VARCHAR (50) NOT NULL)")
    if err != nil {
        log.Fatal().Err(err).Send()
    }
    return db
}

func NewKafkaClient() *kgo.Client {
    cl, err := kgo.NewClient(
        kgo.SeedBrokers(broker),
        kgo.DefaultProduceTopic(topic),
        kgo.ConsumeTopics(topic),
        kgo.AllowAutoTopicCreation(),
    )
    if err != nil {
        log.Fatal().Err(err).Send()
    }
    return cl
}

type Controller struct {
    db *sql.DB
    cl *kgo.Client
}

func NewController(db *sql.DB, cl *kgo.Client) *Controller {
    return &Controller {
       db: db,
       cl: cl,
    }
}

func(c *Controller) Consume(ctx context.Context) {
    for {
        select {
        case <- ctx.Done():
            return
        default:
            fetch := c.cl.PollFetches(ctx)
            if errs := fetch.Errors(); 0 < len(errs) {
                err := errors.New(fmt.Sprint(errs))
                log.Error().Err(err).Send()
            }
            itr := fetch.RecordIter()
            for !itr.Done() {
                rec := itr.Next()
                _, err := c.db.Query("INSERT INTO messages (text) VALUES ($1)", rec.Value)
                if err != nil {
                    log.Error().Err(err).Send()
                    continue
                }
            }
        }
    }
}

func(c *Controller) ProduceRecord(text string) ([]byte, error) {
    rec, err := c.cl.ProduceSync(context.Background(), kgo.StringRecord(text)).First()
    if err != nil {
        return nil, err
    }
    var b []byte
    recf, _ := kgo.NewRecordFormatter(`%{"topic":"%t", "record":"%v"%}\n`)
    return recf.AppendRecord(b, rec), nil
}

func(c *Controller) CountRows() (int, error) {
    var count int
    if err := c.db.QueryRow("SELECT COUNT (*) FROM messages").Scan(&count); err != nil {
        return 0, errors.New("count error")
    }
    return count, nil
}

