package interfaces

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/takashabe/go-router"
)

// printDebugf behaves like log.Printf only in the debug env
func printDebugf(format string, args ...interface{}) {
	if env := os.Getenv("GO_SERVER_DEBUG"); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// ErrorResponse is Error response template
type ErrorResponse struct {
	Message string `json:"reason"`
	Error   error  `json:"-"`
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("reason: %s, error: %#v", e.Message, e.Error)
}

// Respond is response write to ResponseWriter
func Respond(w http.ResponseWriter, code int, src interface{}) {
	var body []byte
	var err error

	switch s := src.(type) {
	case []byte:
		if !json.Valid(s) {
			Error(w, http.StatusInternalServerError, err, "invalid json")
			return
		}
		body = s
	case string:
		body = []byte(s)
	case *ErrorResponse, ErrorResponse:
		// avoid infinite loop
		if body, err = json.Marshal(src); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"reason\":\"failed to parse json\"}"))
			return
		}
	default:
		if body, err = json.Marshal(src); err != nil {
			Error(w, http.StatusInternalServerError, err, "failed to parse json")
			return
		}
	}
	w.WriteHeader(code)
	w.Write(body)
}

// Error is wrapped Respond when error response
func Error(w http.ResponseWriter, code int, err error, msg string) {
	e := &ErrorResponse{
		Message: msg,
		Error:   err,
	}
	printDebugf("%s", e.String())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, e)
}

// JSON is wrapped Respond when success response
func JSON(w http.ResponseWriter, code int, src interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, src)
}

// Server supply HTTP server
type Server struct {
	Entry *EntryHandler
	Token *TokenHandler
}

// Routes returns router
func (s *Server) Routes() *router.Router {
	r := router.NewRouter()

	// For entries
	r.Post("/api/entry/", s.Entry.Post)
	r.Get("/api/entry/:id", s.Entry.Get)
	r.Put("/api/entry/:id", s.Entry.Edit)
	r.Delete("/api/entry/:id", s.Entry.Delete)

	// For tokens
	r.Get("/api/token/:id", s.Token.Get)
	r.Get("/api/token/value/:value", s.Token.FindByValue)

	// Routing of the frontend
	// TODO(takashabe): Want to proxy SPA traffic using a web server.
	r.ServeFile("/", "./public/index.html")
	r.ServeFile("/bundle.js", "./public/bundle.js")
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./public/index.html")
	})

	return r
}

// Run start server
func (s *Server) Run(port int) error {
	log.Printf("Lumber server running at http://localhost:%d/", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.Routes())
}
