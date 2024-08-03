package log

import (
    "os"

    "github.com/rs/zerolog"
)

var l zerolog.Logger

func init() {
    l = zerolog.New(os.Stderr).
        Level(zerolog.InfoLevel).
        With().
        Timestamp().
        Logger()
}

func Info() *zerolog.Event {
    return l.Info()
}

func Error() *zerolog.Event {
    return l.Error()
}

func Fatal() *zerolog.Event {
    return l.Fatal()
}

