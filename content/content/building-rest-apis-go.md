# Building REST APIs with Go

Go's standard library provides excellent support for building HTTP servers and REST APIs. Here's how to create efficient, scalable APIs.

## Basic HTTP Server

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    http.HandleFunc("/users", handleUsers)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "John Doe"},
        {ID: 2, Name: "Jane Smith"},
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
```

## Using Gorilla Mux for Routing

For more complex routing, Gorilla Mux is a popular choice:

```go
import "github.com/gorilla/mux"

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/users", getUsers).Methods("GET")
    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    
    http.ListenAndServe(":8080", r)
}
```

## Key Considerations

- **Middleware**: Use middleware for authentication, logging, and CORS
- **Error Handling**: Implement consistent error responses
- **Validation**: Validate input data thoroughly
- **Testing**: Write comprehensive tests for your endpoints

Go's performance and simplicity make it ideal for building production-ready APIs.