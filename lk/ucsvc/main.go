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
	certFile = flag.String("certFile", "", "Certificate file, if both certFile and keyFile present will use TLS")
	keyFile  = flag.String("keyFile", "", "Private key file, if both certFile and keyFile present will use TLS")
	port     = flag.Int("port", 8090, "Port")
	useTLS   = false
)

func init() {
	flag.Parse()

	if (len(*certFile) > 0 && len(*keyFile) == 0) || (len(*certFile) == 0 && len(*keyFile) > 0) {
		log.Println("You must supply both a certificate file and key file")
		os.Exit(1)
	}
	if (len(*certFile) + len(*keyFile)) > 0 {
		if !fileExists(*certFile) {
			log.Println("Certificate file [%s] does not exist", *certFile)
			os.Exit(1)
		}
		if !fileExists(*keyFile) {
			log.Println("Key file [%s] does not exist", *certFile)
			os.Exit(1)
		}

		useTLS = true
	}
}

func createMux() *echo.Echo {
	mux := echo.New()
	// Global error handler
	mux.HTTPErrorHandler(globalErrorHandler)
	// Middleware
	mux.Use(middleware.Logger)
	mux.Use(basicAuth)

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
	log.Printf("About to listen at %s using TLS: %t", addr, useTLS)
	if useTLS {
		e.RunTLS(addr, *certFile, *keyFile)
	} else {
		e.Run(addr)
	}

	os.Exit(0)
}
