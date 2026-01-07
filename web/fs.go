package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

//go:embed templates static
var EFS embed.FS

// FuncMap com funções customizadas para templates
var templateFuncs = template.FuncMap{
	"formatMonthYear": func(t time.Time) string {
		return t.Format("Jan 2006")
	},
	"formatMonthYearPtr": func(t *time.Time) string {
		if t == nil {
			return "Presente"
		}
		return t.Format("Jan 2006")
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"formatDatePtr": func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02")
	},
	"join": func(sep string, arr []string) string {
		return strings.Join(arr, sep)
	},
	"seniorityLabel": func(s string) string {
		labels := map[string]string{
			"JUNIOR":    "Júnior",
			"MID_LEVEL": "Pleno",
			"SENIOR":    "Sênior",
			"LEAD":      "Tech Lead",
			"PRINCIPAL": "Principal",
		}
		if label, ok := labels[s]; ok {
			return label
		}
		return s
	},
	"locationLabel": func(l string) string {
		labels := map[string]string{
			"ON_SITE": "Presencial",
			"REMOTE":  "Remoto",
			"HYBRID":  "Híbrido",
			"ANY":     "Qualquer",
		}
		if label, ok := labels[l]; ok {
			return label
		}
		return l
	},
}

// GetStaticAssets retorna um FileSystem pronto para ser servido via HTTP
func GetStaticAssets() http.FileSystem {
	f, _ := fs.Sub(EFS, "static")
	return http.FS(f)
}

// ParseTemplate ajuda a parsear templates com layout base e componentes
func ParseTemplate(page string, components ...string) (*template.Template, error) {
	patterns := []string{
		"templates/layouts/base.html",
		"templates/" + page,
	}
	// Adiciona componentes ao parsing
	for _, comp := range components {
		patterns = append(patterns, "templates/components/"+comp)
	}
	return template.New("base.html").Funcs(templateFuncs).ParseFS(EFS, patterns...)
}

// ParseTemplateFragment parseia templates de fragmentos (sem layout base)
func ParseTemplateFragment(templates ...string) (*template.Template, error) {
	paths := make([]string, len(templates))
	for i, t := range templates {
		paths[i] = "templates/" + t
	}
	return template.New("").Funcs(templateFuncs).ParseFS(EFS, paths...)
}
