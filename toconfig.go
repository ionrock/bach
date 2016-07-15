package bach

import (
	"os"
	"text/template"
)

type TmplData struct {
	Filename string
}

func ApplyConfig(t string, c string) error {

	funcMap := template.FuncMap{
		"Get": os.Getenv,
	}

	tmpl, err := template.New(t).Funcs(funcMap).ParseFiles(t)
	if err != nil {
		panic(err)
	}

	d := TmplData{Filename: t}

	fh := os.Stdout
	if c != "" {
		fh, err = os.Create(c)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
	}

	err = tmpl.Execute(fh, d)
	if err != nil {
		panic(err)
	}
	return nil
}
