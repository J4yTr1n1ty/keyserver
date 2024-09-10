package htmx

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/J4yTr1n1ty/keyserver/pkg/models"
)

var (
	ErrParsingTemplate   = fmt.Errorf("error parsing template")
	ErrExecutingTemplate = fmt.Errorf("error executing template")
)

var ErrorTemplate = `
<div class="alert alert-danger">
  <p>{{ .Message }}</p>
</div>
`

var SuccessTemplate = `
<div class="alert alert-success">
  <p>{{ .Message }}</p>
</div>
`

var InfoTemplate = `
<div class="alert alert-info">
  <p>{{ .Message }}</p>
</div>
`

var ListTemplate = `
{{ range . }}
<option value="{{ . }}">{{ . }}</option>
{{ end }}
`

var EmptyListTemplate = `
<option value="" selected disabled>No identities found</option>
`

func RenderError(w http.ResponseWriter, httpStatus int, message string) error {
	w.WriteHeader(httpStatus)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(ErrorTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderSuccess(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(SuccessTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderInfo(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(InfoTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderListIdentities(w http.ResponseWriter, identityList []models.Identity) error {
	w.WriteHeader(http.StatusOK)

	if len(identityList) == 0 {
		value, err := renderTemplate(EmptyListTemplate, []string{})
		if err != nil {
			return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
		}

		w.Write([]byte(value))
		return nil
	}

	names := []string{}
	for _, identity := range identityList {
		names = append(names, fmt.Sprintf("%s (%s)", identity.Name, identity.KeyFingerprint))
	}

	value, err := renderTemplate(ListTemplate, names)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func renderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParsingTemplate, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	return buf.String(), nil
}
