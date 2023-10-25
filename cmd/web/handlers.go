package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"poc.eduardo-luz.eu/internal/models"
	"poc.eduardo-luz.eu/internal/validator"
)

type snippetCreateForm struct {
	Title               string     `form:"title"` //tells the decoder to store the value from the input title in the title field
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //tells the decode to ignore a field during the decoding
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Call newTemplateData with the default data and then append the snippets slice to it
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.go.tmpl", data)
}

func (app *application) sniView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header()["Date"] = nil

	// values of any named parameters will be stored in the request context
	// returns a slice containing names and values params
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//initialize s alice with paths plus base layout and partial
	files := []string{
		"./ui/html/base.go.tmpl",
		"./ui/html/partials/nav.go.tmpl",
		"./ui/html/pages/view.go.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) sniCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.go.tmpl", data)
}

func (app *application) sniCreatePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096) //limits request body to 4kb

	var form snippetCreateForm
	// fills with values the struct with the request and the pointer to the struct. instead of initializing
	// and fulfill the struct manually
	err := app.decodePostForm(r, &form) // requires a non-nil pointer or returns a form.InvalidDecoderError
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field is required")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This fields cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field is required")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// if threre is any error, dump in a plain text HTTP response
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.go.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//changes to format handled by httprouter
	http.Redirect(w, r, fmt.Sprintf("/sni/view/%d", id), http.StatusSeeOther)
}
