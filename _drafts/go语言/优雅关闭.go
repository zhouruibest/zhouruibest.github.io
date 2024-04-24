package main

import (
  "context"
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "time"
)

// Server is an interface that will used by app to serve and shutdown http server
type Server interface {
  ListenAndServe() error
  Shutdown(context.Context) error
}

// Logger is an interface that used by app to log something
type Logger interface {
  Fatalf(format string, v ...interface{})
  Printf(format string, v ...interface{})
}

func handler() http.Handler {
  mux := http.NewServeMux()
  mux.Handle("/healthz", http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "ok")
    },
  ))

  return mux
}

func serve(ctx context.Context, server Server, logger Logger) error {
  var err error

  go func() {
    if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      logger.Fatalf("ListenAndServe()=%+s", err)
    }
  }()

  <-ctx.Done()

  logger.Printf("server stopped")

  ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer func() {
    cancel()
  }()

  if err = server.Shutdown(ctxShutDown); err != nil {
    logger.Fatalf("server.Shutdown(ctxShutdown)=%+s", err)
  }
  
  logger.Printf("server exited gracefully")
  
  if err == http.ErrServerClosed {
    err = nil
  }
  
  return err
}

func main() {
  server := &http.Server{
    Addr:    ":8080",
    Handler: handler(),
  }

  stdLog := log.New(os.Stderr, "", log.LstdFlags)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  ctx, cancel := context.WithCancel(context.Background())
  go func() {
    osCall := <-c
    log.Printf("system call: %+v", osCall)
    cancel()
  }()

  log.Printf("starting server at %s", server.Addr)
  if err := serve(ctx, server, stdLog); err != nil {
    log.Printf("serve(ctx, server, stdLog)=%+s", err)
  }
}