package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/ui"
)

type templateData struct {
	Form            any
	Flash           string
	IsAuthenticated bool
	Superchats      []*models.Superchat
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
