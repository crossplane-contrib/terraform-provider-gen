package template

import (
	"bytes"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TemplateRenderer interface {
	WriteCompiled() error
	DispatchEntry() DispatchEntry
}

type templateRenderer struct {
	path        string
	root        string
	templateDir string
	fs          afero.Fs
	outputBase  string
	packageRoot string
}

func (tr *templateRenderer) fileName() string {
	return strcase.ToSnake(tr.funcName()) + ".go"
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func (tr *templateRenderer) funcName() string {
	_, fname := filepath.Split(tr.path)
	parts := strings.Split(fname, ".")
	return strcase.ToCamel(normalize(parts[0]))
}

func (tr *templateRenderer) packageName() string {
	dir, _ := filepath.Split(tr.path)
	dirs := strings.Split(dir, string(os.PathSeparator))
	// the last element is always a quote because of the dangling slash
	if dirs[len(dirs)-1] == "" {
		return dirs[len(dirs)-2]
	}
	return dirs[len(dirs)-1]
}

func (tr *templateRenderer) relativePath() string {
	templatePath := filepath.Join(tr.root, tr.templateDir)
	rel, _ := filepath.Rel(templatePath, tr.path)
	return rel
}

func (tr *templateRenderer) relativeDir() string {
	rel := tr.relativePath()
	dir, _ := filepath.Split(rel)
	return dir
}

func (tr *templateRenderer) readSource() (string, error) {
	b := bytes.NewBuffer(nil)
	fh, err := tr.fs.Open(tr.path)
	if err != nil {
		return b.String(), err
	}
	defer fh.Close()
	_, err = io.Copy(b, fh)
	return b.String(), err
}

func (tr *templateRenderer) outputWriter() (io.WriteCloser, error) {
	rel := tr.relativeDir()
	outputDir := filepath.Join(tr.outputBase, rel)
	err := tr.fs.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	fh, err := tr.fs.Create(filepath.Join(outputDir, tr.fileName()))
	if err != nil {
		return nil, err
	}
	return fh, nil
}

func (tr *templateRenderer) WriteCompiled() error {
	contents, err := tr.readSource()
	if err != nil {
		return err
	}
	f := jen.NewFile(tr.packageName())
	f.Func().Id(tr.funcName()).Params().String().Block(
		jen.Return(jen.Lit(contents)))
	fh, err := tr.outputWriter()
	if err != nil {
		return err
	}
	defer fh.Close()
	return f.Render(fh)
}

func (tr *templateRenderer) DispatchEntry() DispatchEntry {
	return DispatchEntry{
		FuncName:       tr.funcName(),
		RelativeImport: tr.relativeDir(),
		RelativePath:   tr.relativePath(),
		PackageRoot:    tr.packageRoot,
	}
}

func NewTemplateRenderer(path, root, templateDir, outputBase, packageRoot string, fs afero.Fs) TemplateRenderer {
	return &templateRenderer{
		path:        path,
		root:        root,
		templateDir: templateDir,
		packageRoot: packageRoot,
		fs:          fs,
		outputBase:  outputBase,
	}
}
