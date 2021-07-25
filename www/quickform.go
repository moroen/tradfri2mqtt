package www

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"reflect"
)

type FieldDefinition struct {
	Id        string
	FieldType reflect.Kind
	Label     string
	Text      string
}

type FormDefinition struct {
	Action string
	Fields template.HTML
}

func (f FieldDefinition) Render() template.HTML {
	switch f.FieldType {
	case reflect.String:
		return template.HTML(fmt.Sprintf("<div>%s:</div> <div><input></input></div>", f.Label))
	}
	return template.HTML("")
}

func QuickForm(action string, fields interface{}) template.HTML {
	/* tpl_field, err := template.New("field").Parse("<div>{{.Label}}</div>")
	if err != nil {
		log.Fatal(err.Error())
	}
	*/
	val := reflect.ValueOf(fields)

	if val.Kind() != reflect.Struct {
		log.Fatal("Wrong type")
	}

	structType := val.Type()

	// formName := structType.Name()

	templates, err := template.ParseGlob("www/widgets/*.tmpl")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var tpl bytes.Buffer
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag
		templates.ExecuteTemplate(&tpl, "input.tmpl", FieldDefinition{field.Name, reflect.String, tag.Get("label"), "Text"})
	}

	form_def := FormDefinition{Action: action, Fields: template.HTML(tpl.Bytes())}
	tpl.Reset()
	templates.ExecuteTemplate(&tpl, "form.tmpl", form_def)
	return template.HTML(tpl.Bytes())
}
