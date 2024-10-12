package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Southclaws/enumerator/generate"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	if err := run(root); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(path string) error {
	fs := token.NewFileSet()

	pkgs, err := parser.ParseDir(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return nil
	}

	for _, pkg := range pkgs {
		if strings.Contains(pkg.Name, "_test") {
			continue
		}
		if len(pkg.Files) == 0 {
			continue
		}

		enums := make(map[string][]generate.Value)

		output := filepath.Join(path, fmt.Sprintf("%s_enum_gen.go", pkg.Name))

		for _, file := range pkg.Files {
			m := processFile(file)
			for k, v := range m {
				enums[k] = v
			}
		}

		if len(enums) == 0 {
			continue
		}

		var enumNames []string
		for k := range enums {
			enumNames = append(enumNames, k)
		}
		sort.Strings(enumNames)

		es := []generate.Enum{}
		for _, sourceName := range enumNames {
			values := enums[sourceName]
			name := generate.Title(strings.TrimSuffix(sourceName, "Enum"))

			es = append(es, generate.Enum{
				Name:       name,
				Values:     values,
				Sourcename: sourceName,
			})
		}

		of, err := os.Create(output)
		if err != nil {
			return err
		}

		err = generate.Generate(output, pkg.Name, es, of)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(file *ast.File) map[string][]generate.Value {
	enumTypes := gatherEnumTypes(file)

	types := make(map[string][]generate.Value)

	for _, d := range file.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if decl.Tok != token.CONST {
			continue
		}

		if len(decl.Specs) < 2 {
			continue
		}

		for _, spec := range decl.Specs {
			v, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			ident, ok := v.Type.(*ast.Ident)
			if !ok {
				continue
			}

			name := ident.Name

			_, ok = enumTypes[name]
			if !ok {
				continue
			}

			l := types[name]

			token := v.Names[0].Name

			var pretty string
			if v.Comment != nil {
				pretty = strings.Trim(v.Comment.Text(), "\n\t")
			}

			l = append(l, generate.Value{
				Token:  token,
				Pretty: pretty,
			})

			types[name] = l
		}
	}

	return types
}

func gatherEnumTypes(file *ast.File) map[string]*ast.TypeSpec {
	types := make(map[string]*ast.TypeSpec)

	for _, d := range file.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if decl.Tok != token.TYPE {
			continue
		}

		spec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		_, ok = spec.Type.(*ast.Ident)
		if !ok {
			continue
		}

		if !strings.HasSuffix(spec.Name.Name, "Enum") {
			continue
		}

		types[spec.Name.Name] = spec
	}

	return types
}
