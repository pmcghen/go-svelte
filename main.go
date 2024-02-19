package main

import (
  "net/http"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := chi.NewRouter()

  r.Use(middleware.Logger)
  r.Use(JSONMiddleware)

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world"))
  })

  r.Get("/tasks", tasks)
  r.Post("/tasks", createTask)
  r.Patch("/tasks/{taskID}", updateTask)
  r.Delete("/tasks/{taskID}", deleteTask)
  http.ListenAndServe(":8080", r)
}

func JSONMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    next.ServeHTTP(w, r)
  })
}

