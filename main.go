package main

import (
	"bytes"
	"flag"
	"fmt"
	htmpl "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	tmpl = ".tmpl"
)

var (
	config Config
	// https://choosealicense.com/licenses
	supportedLicenses = []string{"GPL3", "MIT", "APACHE", "UNLICENSE"}
)

// Config ...
type Config struct {
	Project     string
	ProjectRoot string
	Template    string
	License     htmpl.HTML
}

func main() {
	fmt.Println(`
   ____   ____   ____
  / ___\ /  _ \ / ___\
 / /_/  >  <_> ) /_/  >
 \___  / \____/\___  /
/_____/       /_____/`)

	parseConfig()
	generateProject()

	fmt.Println("\n\nmake setup and happy hacking!")
}

func parseConfig() {
	var err error
	var template string
	var project string
	var license string

	flag.StringVar(&project, "project", "", "project name")
	flag.StringVar(&template, "template", "", "project template")
	flag.StringVar(&license, "license", "GPL3", "license for project")
	flag.Parse()

	if template == "" {
		template = findGogTemplate()
	}
	config.Template = template

	if project == "" {
		log.Fatalln("must specify project name")
	}
	config.Project = project

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	root := pwd + "/" + project
	config.ProjectRoot = root

	if isLicenseSupported(supportedLicenses, license) {
		txt, err := readLicense(license)
		if err == nil {
			config.License = htmpl.HTML(string(txt))
		} else {
			log.Println(err)
		}
	} else {
		log.Println("currently not support this license " + license)
	}
}

func generateProject() {
	err := filepath.Walk(config.Template, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return parseTemplateAndOutput(path, f)
	})
	if err != nil {
		log.Fatalln(err)
	}
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
		tpl := htmpl.Must(htmpl.New("gog").Parse(buf.String()))
		// discard previous content
		buf.Reset()
		if err := tpl.Execute(buf, config); err != nil {
			log.Fatal(err)
		}
	}
	// TODO: write executable file
	return ioutil.WriteFile(p, buf.Bytes(), 0644)
}

func isLicenseSupported(licenses []string, license string) bool {
	l := strings.ToUpper(license)
	for i := 0; i < len(licenses); i++ {
		if l == licenses[i] {
			return true
		}
	}
	return false
}

func readLicense(name string) ([]byte, error) {
	file := findGogRoot() + "/licenses" + "/" + strings.ToUpper(name)
	return ioutil.ReadFile(file)
}

func findGogRoot() string {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	for _, path := range gopaths {
		root := strings.TrimSpace(path) + "/src/github.com/cosmtrek/gog"
		if _, err := os.Stat(root); err == nil {
			return root
		}
	}
	return ""
}

func findGogTemplate() string {
	root := findGogRoot()
	return root + "/template"
}
