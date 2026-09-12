package controllers

import (
	"html/template"
	"io/fs"
	"reshier/utils"
)

// ViewsFS adalah filesystem untuk template views.
// Di-set oleh api/index.go (Vercel) menggunakan embed.FS,
// atau tetap nil untuk local dev (fallback ke os filesystem).
var ViewsFS fs.FS

// parseTemplate mem-parse template dari ViewsFS (embed) jika tersedia,
// atau dari filesystem OS (local dev) jika tidak.
func parseTemplate(name string, path string) (*template.Template, error) {
	if ViewsFS != nil {
		return template.New(name).Funcs(utils.TemplateFuncs).ParseFS(ViewsFS, path)
	}
	return template.New(name).Funcs(utils.TemplateFuncs).ParseFiles(path)
}
