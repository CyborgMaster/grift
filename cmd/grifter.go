package cmd

import (
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/markbates/going/randx"
)

type grifter struct {
	CurrentDir        string
	BuildPath         string
	GriftsPackagePath string
	ExePath           string
	Verbose           bool
}

func newGrifter() (*grifter, error) {
	g := &grifter{}

	pwd, err := os.Getwd()
	if err != nil {
		return g, err
	}
	g.CurrentDir = pwd
	base := randx.String(10)
	g.BuildPath = path.Join(os.Getenv("GOPATH"), "src", "grift.build", base)
	g.GriftsPackagePath = path.Join("grift.build", base, "grifts")
	return g, nil
}

func (g *grifter) Setup() error {
	err := os.MkdirAll(g.BuildPath, 0777)
	if err != nil {
		return err
	}

	return g.Build()
}

func (g *grifter) Build() error {
	err := g.copyGrifts()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(g.BuildPath, "grifts", "grift_loader.go"), []byte(loaderTmpl), 0644)
	if err != nil {
		return err
	}

	t, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(g.BuildPath, "main.go"))
	if err != nil {
		return err
	}

	err = t.Execute(f, g)
	if err != nil {
		return err
	}

	g.ExePath = path.Join(g.BuildPath, "main.go")
	return nil
}

func (g *grifter) TearDown() error {
	return os.RemoveAll(g.BuildPath)
}

func (g *grifter) copyGrifts() error {
	cp := exec.Command("cp", "-rv", path.Join(g.CurrentDir, "grifts"), g.BuildPath)
	return cp.Run()
}
