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
<option value="{{ .Fingerprint }}">{{ .Name }} ({{ .Fingerprint }})</option>
{{ end }}
`

type ListData struct {
	Fingerprint string
	Name        string
}

var EmptyListTemplate = `
<option value="" selected disabled>No identities found</option>
`

var KeyTableTemplate = `
<table class="table table-hover">
  <thead>
    <tr>
      <th scope="col">Name</th>
      <th scope="col">Fingerprint</th>
      <th scope="col">Comment</th>
      <th scope="col">Email</th>
      <th scope="col">Download</th>
    </tr>
  </thead>
  <tbody>
    {{ range . }}
    <tr>
      <td>{{ .Name }}</td>
      <td>{{ .KeyFingerprint }}</td>
      <td>{{ .Comment }}</td>
      <td>{{ .Email }}</td>
      {{if .Email}}
      <td><a href="/key/{{ .Email }}">Download</td>
      {{else}}
      <td></td>
      {{end}}
    </tr>
    {{ end }}
  </tbody>
</table>
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

	data := make([]ListData, len(identityList))
	for _, identity := range identityList {
		data = append(data, ListData{
			Name:        identity.Name,
			Fingerprint: identity.KeyFingerprint,
		})
	}

	value, err := renderTemplate(ListTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderKeyTable(w http.ResponseWriter, identityList []models.Identity) error {
	w.WriteHeader(http.StatusOK)

	if len(identityList) == 0 {
		value, err := renderTemplate(EmptyListTemplate, []string{})
		if err != nil {
			return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
		}

		w.Write([]byte(value))
		return nil
	}

	data := append([]models.Identity{}, identityList...)

	value, err := renderTemplate(KeyTableTemplate, data)
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
