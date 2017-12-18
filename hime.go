package hime

import (
	"context"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/acoshift/middleware"
)

// App is the hime app
type App interface {
	http.Handler

	// Builder

	// TemplateDir sets directory to load template,
	// default is "view"
	TemplateDir(path string) App

	// TemplateRoot sets root layout using t.Lookup,
	// default is "layout"
	TemplateRoot(name string) App

	// TemplateFuncs adds template funcs while load template
	TemplateFuncs(funcs ...template.FuncMap) App

	// Component adds given templates to every templates
	Component(filename ...string) App

	// Template loads template into memory
	Template(name string, filename ...string) App

	// BeforeRender runs given middleware for before render,
	// ex. View, JSON, String, Bytes, CopyForm, etc.
	BeforeRender(m middleware.Middleware) App

	// Minify enables minify when render html, css, js
	Minify() App

	// Handler sets the handler factory
	Handler(factory HandlerFactory) App

	// Routes registers route name and path
	Routes(routes Routes) App

	// GracefulShutdown runs server as graceful shutdown,
	// can works only when start server with app.ListenAndServe
	GracefulShutdown() App

	// ShutdownTimeout sets graceful shutdown timeout
	ShutdownTimeout(d time.Duration) App

	// ListenAndServe starts web server
	ListenAndServe(addr string) error

	// Route gets route path from given name
	Route(name string) string
}

// Routes is the map for route name => path
type Routes map[string]string

// HandlerFactory is the function for create router
type HandlerFactory func(App) http.Handler

// Handler is the hime handler
type Handler func(Context) Result

// Result is the handler result
type Result interface {
	Response(w http.ResponseWriter, r *http.Request)
}

// ResultFunc is the result function
type ResultFunc func(w http.ResponseWriter, r *http.Request)

// Response implements Result interface
func (f ResultFunc) Response(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// Context is the hime context
type Context interface {
	context.Context

	// Request returns http.Request from context
	Request() *http.Request

	// ResponseWrite returns http.ResponseWriter from context
	ResponseWriter() http.ResponseWriter

	// Status sets response status
	Status(code int) Context

	// Request functions
	ParseForm() error
	ParseMultipartForm(maxMemory int64) error
	Form() url.Values
	PostForm() url.Values
	FormValue(key string) string
	PostFormValue(key string) string
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)
	MultipartForm() *multipart.Form
	MultipartReader() (*multipart.Reader, error)
	Method() string

	// Results

	// Redirect redirects to given url
	Redirect(url string) Result

	// SafeRedirect extracts only path from url then redirect
	SafeRedirect(url string) Result

	// RedirectTo redirects to named route
	RedirectTo(name string) Result

	// Error wraps http.Error
	Error(error string) Result

	// View renders template
	View(name string, data interface{}) Result

	// JSON renders json
	JSON(data interface{}) Result

	// String renders string with format
	String(format string, a ...interface{}) Result

	// CopyFrom copies source into response writer
	CopyFrom(src io.Reader) Result

	// Bytes renders bytes
	Bytes(b []byte) Result

	// Handle wrap h with Result
	Handle(h http.Handler) Result
}
