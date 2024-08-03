package handler

import (
    "io"
    "time"
    "errors"
    "net/http"
    "encoding/json"
    "encoding/base64"

    "github.com/gorilla/mux"

    "github.com/rchnmy/xupery/log"
    "github.com/rchnmy/xupery/pkg/controller"
)

type Token struct { t string }

func(t *Token) Encode(s string) {
    *t = Token { t: base64.URLEncoding.EncodeToString([]byte(s)) }
    log.Info().Msgf("X-Upery-Token: %s", t.t)
}

func NewRouter() *mux.Router {
    return mux.NewRouter()
}

func(t *Token) Validate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-Upery-Token") != t.t {
            w.WriteHeader(403)
            log.Error().Err(errors.New("missing token")).Send()
            return
        }
        next.ServeHTTP(w, r)
    })
}

func JsonHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

type report struct { MessagesTotal int `json:"messages_total"` }

func GetStatistic(c *controller.Controller) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        count, err := c.CountRows()
        if err != nil {
            w.WriteHeader(500)
            log.Error().Err(err).Send()
            return
        }
        rep := report{ MessagesTotal: count }
        if err := json.NewEncoder(w).Encode(&rep); err != nil {
            w.WriteHeader(500)
            log.Error().Err(err).Send()
            return
        }
    }
}

type message struct { Text string `json:"text"` }

func PostMessage(c *controller.Controller) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var msg message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            w.WriteHeader(400)
            log.Error().Err(err).Send()
            return
        }
        defer drop(r.Body)
        if msg.Text == "" || 50 < len(msg.Text) {
            w.WriteHeader(400)
            log.Error().Err(errors.New("message not valid")).Send()
            return
        }
        rec, err := c.ProduceRecord(msg.Text)
        if err != nil {
            w.WriteHeader(500)
            log.Error().Err(err).Send()
            return
        }
        w.Write(rec)
    }
}

func NewServer(a string, r *mux.Router) *http.Server {
    return &http.Server {
        Addr:         a,
        Handler:      r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 20 * time.Second,
    }
}

func Serve(s *http.Server) {
    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal().Err(err).Send()
    }
}

func drop(rc io.ReadCloser) {
    rc.Close()
    io.Copy(io.Discard, rc)
}

