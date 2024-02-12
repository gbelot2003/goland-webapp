package render

import (
	"bytes"
	"firstwebapp/cmd/pkg/config"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// set the config for the new templates packages
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string) {

	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get the template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, nil)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil

}
