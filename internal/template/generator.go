package template

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (tc *TemplateCompiler) Walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	tc.sourcePaths = append(tc.sourcePaths, path)
	return nil
}

type TemplateCompiler struct {
	sourcePaths []string
	root        string
	templateDir string
	output      string
	packageRoot string
	fs          afero.Fs
}

func NewTemplateCompiler(fs afero.Fs, root, templateDir, output, packageRoot string) *TemplateCompiler {
	return &TemplateCompiler{
		fs:          fs,
		root:        root,
		templateDir: templateDir,
		output:      output,
		packageRoot: packageRoot,
		sourcePaths: make([]string, 0),
	}
}

type DispatchEntry struct {
	RelativePath   string
	RelativeImport string
	FuncName       string
	PackageRoot    string
}

func (de DispatchEntry) ImportAlias() string {
	parts := strings.Split(de.RelativeImport, string(os.PathSeparator))
	return strcase.ToLowerCamel(strings.Join(parts, " "))
}

func (de DispatchEntry) dispatchMapEntry() (jen.Code, jen.Code) {
	key := jen.Lit(de.RelativePath)
	val := jen.Qual(de.importPath(), de.FuncName)
	return key, val
}

func (tc *TemplateCompiler) CompileGeneratedTemplates() error {
	err := afero.Walk(tc.fs, filepath.Join(tc.root, tc.templateDir), tc.Walk)
	if err != nil {
		return err
	}
	dispatchers := make([]DispatchEntry, 0)
	for _, source := range tc.sourcePaths {
		t := NewTemplateRenderer(source, tc.root, tc.templateDir, tc.output, tc.packageRoot, tc.fs)
		err := t.WriteCompiled()
		if err != nil {
			return err
		}
		dispatchers = append(dispatchers, t.DispatchEntry())
	}
	return tc.renderDispatchMap(dispatchers)
}

func (de DispatchEntry) importPath() string {
	return path.Join(de.PackageRoot, de.RelativeImport)
}

func (tc *TemplateCompiler) renderDispatchMap(dispatchers []DispatchEntry) error {
	f := jen.NewFile("dispatch")
	mapEntries := make(jen.Dict)
	for _, d := range dispatchers {
		importPath := d.importPath()
		alias := d.ImportAlias()
		f.ImportAlias(importPath, alias)
		k, v := d.dispatchMapEntry()
		mapEntries[k] = v
	}
	mapType := jen.Map(jen.String()).Func().Params().String()
	mapAssignment := jen.Op("=").Map(jen.String()).Func().Params().String().Values(mapEntries)
	f.Var().Id("TemplateDispatchMap").Add(mapType).Add(mapAssignment)
	fh, err := tc.fs.Create(filepath.Join(tc.output, "dispatch.go"))
	if err != nil {
		return err
	}
	defer fh.Close()
	return f.Render(fh)
}
