package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FileParser struct {
	filePath string
}

func GetFileParser(filePath string) *FileParser {
	return &FileParser{filePath: filePath}
}

func GetFileNames(directory string) ([]string, error) {
	fileNames := make([]string, 0)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	return fileNames, err
}

func (f *FileParser) Parse() (*FileDesc, error) {
	// parse file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, f.filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %s", err)
	}
	importsDic, err := f.parseImport(node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse import: %s", err)
	}
	structs := make([]*StructDesc, 0)
	interfaces := make([]*InterfaceDesc, 0)
	genDecls, err := getGenDecls(node)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %s", err)
	}
	if len(genDecls) == 0 {
		return nil, nil
	}
	for _, genDecl := range genDecls {
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					structParser := NewStructParser(typeSpec.Name.Name, typeSpec, genDecl, node, importsDic)
					structDesc, err := structParser.Parse()
					if err != nil {
						return nil, fmt.Errorf("failed to parse struct: %s", err)
					}
					structs = append(structs, structDesc)
				} else if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					interfaceParser := NewInterfaceParser(typeSpec.Name.Name, typeSpec, genDecl, importsDic)
					interfaceDesc, err := interfaceParser.Parse()
					if err != nil {
						return nil, fmt.Errorf("failed to parse interface: %s", err)
					}
					interfaces = append(interfaces, interfaceDesc)
				}
			}
		}

	}
	fileInfo, err := os.Stat(f.filePath)
	moduleName := getModuleName()
	fileDesc := &FileDesc{
		FileName:        fileInfo.Name(),
		PackageName:     node.Name.Name,
		FullPackageName: getFullPackageName(moduleName, f.filePath),
		//RelativePath: "", // todo unimplemented
		Imports:    importsDic,
		Structs:    structs,
		Interfaces: interfaces,
	}
	return fileDesc, nil
}

func getGenDecls(file *ast.File) (list []*ast.GenDecl, err error) {
	list = make([]*ast.GenDecl, 0)
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			list = append(list, genDecl)
		}
	}
	return list, err
}

func (f *FileParser) parseImport(file *ast.File) (result map[string]*ImportDesc, err error) {
	result = make(map[string]*ImportDesc)
	ast.Inspect(file, func(n ast.Node) bool {
		if importSpec, ok := n.(*ast.ImportSpec); ok {
			path := strings.Trim(importSpec.Path.Value, "\"")
			name := ""
			hasAlias := importSpec.Name != nil
			if importSpec.Name != nil {
				name = importSpec.Name.Name
			} else {
				parts := strings.Split(path, "/")
				name = parts[len(parts)-1]
			}
			if importSpec.Name != nil {
				name = importSpec.Name.Name
			}
			result[name] = &ImportDesc{
				Path:     path,
				HasAlias: hasAlias,
				Name:     name,
			}
		}
		return true
	})
	return result, err
}
