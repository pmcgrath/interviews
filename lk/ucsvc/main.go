package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	port = flag.Int("port", 8090, "Port")
)

func init() {
	flag.Parse()
}

func createMux() *echo.Echo {
	mux := echo.New()
	// Global error handler
	mux.HTTPErrorHandler(globalErrorHandler)
	// Middleware
	mux.Use(middleware.Logger)

	/* Routes
	Urls				Verb	Body				Description
	------------------------------- ------	------------------------------- ---------------------------
	/users				GET	Empty				Returns all users
	/users				POST	User name			Creates a new user with no connections
	/users/:id			GET	Empty				Gets a user
	/users/:id			DELETE	Empty				Deletes a user
	/users/:id/connections		GET	Empty				Get a user's connections
	/users/:id1/connections/:id2	PUT	Empty				Creates a connection between users
	/users/:id1/connections/:id2	DELETE	Empty				Deletes the connection between users
	*/
	mux.Get("/users", getUsers)
	mux.Post("/users", createUser)
	mux.Get("/users/:id", getUser)
	mux.Delete("/users/:id", deleteUser)
	mux.Get("/users/:id/connections", getUserConnections)
	mux.Put("/users/:id1/connections/:id2", createUserConnection)
	mux.Delete("/users/:id1/connections/:id2", deleteUserConnection)

	return mux
}

func main() {
	userConns = newUserConnections(newMemoryStore())

	e := createMux()

	// Start
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Using runtime %s\n", runtime.Version())
	log.Printf("Commit = %s build @ %s Full commit = %s\n", shortCommitHash, buildDate, commitHash)
	log.Printf("About to listen at %s", addr)
	e.Run(addr)

	os.Exit(0)
}
