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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	//changes to format handled by httprouter
	http.Redirect(w, r, fmt.Sprintf("/sni/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.go.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field is required")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 chars long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.go.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEmail) {
			form.AddFieldError("email", "Email already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.go.tmpl", data)
		} else {
			app.serverError(w, err)
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Signup was successful!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.go.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.go.tmpl", data)
		return
	}

	//checks if the credentials are valid
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.go.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//generates new session id during the authenticate state change
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/sni/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout user")
}
