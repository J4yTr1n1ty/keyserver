package htmx

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
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

func RenderError(w http.ResponseWriter, httpStatus int, message string) error {
	w.WriteHeader(httpStatus)

	value, err := renderResponse(ErrorTemplate, message)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderSuccess(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	value, err := renderResponse(SuccessTemplate, message)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderInfo(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	value, err := renderResponse(InfoTemplate, message)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func renderResponse(templateStr string, message string) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParsingTemplate, err)
	}

	data := map[string]string{"Message": message}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	return buf.String(), nil
}
