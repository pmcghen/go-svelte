package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strconv"

  "github.com/go-chi/chi/v5"
)

const (
  JSONFile = "./tasks.json"
)

type Task struct {
  ID int `json:"id"`
  Name string `json:"name"`
  Done bool `json:"done"`
}

type CreateTaskBody struct {
  Name string `json:"name"`
}

type UpdateTaskBody struct {
  Name *string `json:"name"`
  Done *bool `json:"done"`
}

func tasks(w http.ResponseWriter, r *http.Request) {
  jsonFile, err := ioutil.ReadFile(JSONFile)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error reading tasks from jsonfile %v", err)
    return
  }

  w.Write(jsonFile)
}

func createTask(w http.ResponseWriter, r *http.Request) {
  body := CreateTaskBody{}
  
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error decoding body %v", err)
    return
  }

  jsonFile, err := os.Open(JSONFile)
  
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error decoding body %v", err)
    return
  }

  tasks := []Task{}

  if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error decoding json file %v", err)
    return
  }

  newTask := Task{
    Name: body.Name,
    Done: false,
    ID: len(tasks) + 1,
  }

  tasks = append(tasks, newTask)

  j, err := json.Marshal(tasks)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error marshaling tasks %v", err)
    return
  }

  err = ioutil.WriteFile(JSONFile, j, 0755)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error writing JSON %v", err)
    return
  }

  j, err = json.Marshal(newTask)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error marshaling tasks %v", err)
    return
  }

  w.WriteHeader(http.StatusCreated)
  w.Write(j)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
  body := UpdateTaskBody{}

  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error updating task %v", err)
    return
  }

  jsonFile, err := os.Open(JSONFile)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error reading tasks %v", err)
    return
  }

  tasks := []Task{}

  if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error decoding json %v", err)
    return
  }

  taskId, err := strconv.Atoi(chi.URLParam(r, "taskID"))

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error converting string to int %v", err)
    return
  }

  for i, task := range tasks {
    if task.ID == taskId {
      if body.Name != nil {
        task.Name = *body.Name
      }

      if body.Done != nil {
        task.Done = *body.Done
      }

      tasks[i] = task
    }
  }

  j, err := json.Marshal(tasks)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error marshaling tasks %v", err)
    return
  }

  err = ioutil.WriteFile(JSONFile, j, 0755)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error writing task %v", err)
    return
  }

  w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
  jsonFile, err := os.Open(JSONFile)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error opening tasks %v", err)
    return
  }

  tasks := []Task{}

  if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error decoding json %v", err)
    return
  }

  prevLength := len(tasks)

  taskId, err := strconv.Atoi(chi.URLParam(r, "taskID"))

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error converting string to int %v", err)
    return
  }

  for i, task := range tasks {
    if task.ID == taskId {
      tasks = append(tasks[:i], tasks[i+1:]... )
    }
  }

  currLength := len(tasks)

  if currLength == prevLength {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  j, err := json.Marshal(tasks)

  if err != nil {
    return
  }

  err = ioutil.WriteFile(JSONFile, j, 0755)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Error writing json %v", err)
    return
  }

  w.WriteHeader(http.StatusOK)
}

