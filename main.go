package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config holds all info
type Config struct {
	Project     string
	ProjectRoot string
	Template    string
}

var (
	config Config
)

const (
	tmpl = ".tmpl"
)

func main() {
	fmt.Println(`
   ____   ____   ____
  / ___\ /  _ \ / ___\
 / /_/  >  <_> ) /_/  >
 \___  / \____/\___  /
/_____/       /_____/

`)

	var err error
	var template string
	var project string
	flag.StringVar(&template, "template", "", "project template")
	flag.StringVar(&project, "project", "", "project name")
	flag.Parse()

	if template == "" {
		log.Fatal("must specify template dir")
	}
	if project == "" {
		log.Fatal("must specify project name")
	}
	config.Template = template
	config.Project = project

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	root := pwd + "/" + project
	config.ProjectRoot = root

	err = filepath.Walk(template, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return parseTemplateAndOutput(path, f)
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(`
Happy hacking!
`)
}

func parseTemplateAndOutput(path string, file os.FileInfo) error {
	ps := strings.Split(path, "template")
	fn := ps[len(ps)-1]
	if strings.Contains(fn, tmpl) {
		fn = strings.Split(fn, tmpl)[0]
	}

	p := config.ProjectRoot + fn
	if file.IsDir() {
		return os.MkdirAll(p, 0755)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	if strings.Contains(file.Name(), tmpl) {
		tpl := template.Must(template.New("gog").Parse(buf.String()))
		// discard previous content
		buf.Reset()
		if err := tpl.Execute(buf, config); err != nil {
			log.Fatal(err)
		}
	}
	return ioutil.WriteFile(p, buf.Bytes(), 0644)
}
