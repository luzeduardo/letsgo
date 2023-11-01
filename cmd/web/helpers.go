package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) decodePostForm(r *http.Request, destination any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(destination, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}
	return nil
}

// returns a pointer to a templateData struct initialized with the current year
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r), //adds the auth status to template data
		CSRFToken:       nosurf.Token(r),        //adds the token every time a pages get rendered
	}
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	//retrieve the template set from the cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	//write template to buffer instead of to the ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//write out the HTTP status code
	// if template is written to the buffer without errors,
	//it is ok to go ahead and write the http status to ResponseWriter
	w.WriteHeader(status)

	//wirte the content of the buffer to the ResponseWriter
	buf.WriteTo(w)
}

// writes error message and stack trace to errorLog and return a http 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends a specfic status code and description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) info(w http.ResponseWriter, content string) {
	app.infoLog.Output(2, content)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextkey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
