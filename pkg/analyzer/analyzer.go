package analyzer

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// New returns new constant analyzer.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "constantcheck",
		Doc:      "A linter that detect the possibility to use constants.",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	var (
		insp             *inspector.Inspector
		packages         []*ast.File
		packageConstList map[string]string
		literalsFromFile []*ast.BasicLit
		importsFromFile  []*ast.BasicLit
	)
	insp = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	packages = getFilesFromImport(insp)
	packageConstList = getAllConstantFromImportingFiles(packages)
	literalsFromFile, importsFromFile = getAllLiteralsAndImportsFromFile(insp)

	for _, v := range literalsFromFile {
		if containsValueInImports(importsFromFile, v.Value) {
			continue
		}
		if key := containsValueInMap(packageConstList, v.Value); key != "" {
			pass.Reportf(v.Pos(), "%v literal contains in constant with name %v",
				getBasicLitValue(v), key)
		}
	}

	return nil, nil
}

// getAllLiteralsFromFile parses all literals and imports files from the parsed file
func getAllLiteralsAndImportsFromFile(insp *inspector.Inspector) (
	literalsList []*ast.BasicLit, importsList []*ast.BasicLit) {
	types := []ast.Node{
		(*ast.BasicLit)(nil),
		(*ast.ImportSpec)(nil),
	}
	insp.Preorder(types, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.BasicLit:
			if n.Kind == token.STRING && n.Value != "\"\"" {
				literalsList = append(literalsList, n)
			}
		case *ast.ImportSpec:
			importsList = append(importsList, n.Path)
		}
	})
	return
}

// getAllConstantFromImportingFiles takes all constants from import files
func getAllConstantFromImportingFiles(packages []*ast.File) (constList map[string]string) {
	constList = make(map[string]string)
	for _, file := range packages {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok != token.CONST {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)

				for _, vn := range valueSpec.Names {
					if len(vn.Obj.Decl.(*ast.ValueSpec).Values) > 0 {
						if !isStartedWithUppercase(vn.Name) {
							continue
						}
						switch v := vn.Obj.Decl.(*ast.ValueSpec).Values[0].(type) {
						case *ast.BasicLit:
							constList[vn.Name] = v.Value
						case *ast.Ident:
							constList[vn.Name] = v.Name
						case *ast.UnaryExpr:
							switch uv := v.X.(type) {
							case *ast.BasicLit:
								constList[vn.Name] = uv.Value
							}
						}
					}
				}
			}
		}
	}
	return
}

// getFilesFromImport get paths to all import files
func getFilesFromImport(insp *inspector.Inspector) []*ast.File {
	var importsSlice []string

	types := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// pick up all imports filename
	insp.Preorder(types, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				switch spec := spec.(type) {
				case *ast.ImportSpec:
					filename := spec.Path.Value[1 : len(spec.Path.Value)-1]
					if strings.Contains(filename, "/") {
						importsSlice = append(importsSlice, filename)
					}
				}
			}
		}
	})

	// pick up all the files next to the import
	var packages []*ast.Package
	for _, file := range importsSlice {
		if strings.Contains(file, "_test") {
			continue
		}
		bi, err := build.Import(file, "", 0)
		if err != nil {
			log.Fatal(err)
		}

		fileCheck := func(fi os.FileInfo) bool {
			nm := fi.Name()
			for _, f := range bi.GoFiles {
				if f == nm {
					return true
				}
			}
			return false
		}

		pkgs, err := parser.ParseDir(token.NewFileSet(), bi.Dir, fileCheck, 0)
		if err != nil {
			log.Fatal(err)
		}

		packages = append(packages, pkgs[bi.Name])
	}

	var files []*ast.File
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
	}
	return files
}

// containsValueInMap checks if there is a value in the map
func containsValueInMap(fConst map[string]string, value string) string {
	for k, v := range fConst {
		if v == value {
			return k
		}
	}
	return ""
}

// containsValueInImports checks if there is a value in the importsList
func containsValueInImports(impotsList []*ast.BasicLit, value string) bool {
	for _, iv := range impotsList {
		if iv.Value == value {
			return true
		}
	}
	return false
}

// getBasicLitValue returns BasicLit value as string without quotes.
func getBasicLitValue(basicLit *ast.BasicLit) string {
	var val strings.Builder
	for _, r := range basicLit.Value {
		if r == '"' {
			continue
		} else {
			val.WriteRune(r)
		}
	}
	return val.String()
}

func isStartedWithUppercase(str string) bool {
	for i, r := range str {
		if unicode.IsUpper(r) {
			return true
		}
		if i != 0 {
			break
		}
	}
	return false
}
