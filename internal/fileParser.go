package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type FileParser struct {
	fileName string
}

func GetFileParser(fileName string) *FileParser {
	return &FileParser{fileName: fileName}
}

func (f *FileParser) Parse() (*FileDesc, error) {
	// 解析文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, f.fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %s", err)
	}
	importsDic, err := f.parseImport(node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse import: %s", err)
	}
	structs := make([]*StructDesc, 0)
	interfaces := make([]*InterfaceDesc, 0)
	serviceSpecs, err := getServiceSpecs(node)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %s", err)
	}
	if len(serviceSpecs) == 0 {
		return nil, nil
	}
	for _, serviceSpec := range serviceSpecs {
		if _, ok := serviceSpec.Type.(*ast.StructType); ok {
			structParser := NewStructParser(serviceSpec.Name.Name, serviceSpec, node)
			structDesc, err := structParser.Parse()
			if err != nil {
				return nil, fmt.Errorf("failed to parse struct: %s", err)
			}
			structs = append(structs, structDesc)
		} else if _, ok := serviceSpec.Type.(*ast.InterfaceType); ok {
			interfaceParser := NewInterfaceParser(serviceSpec.Name.Name, serviceSpec)
			interfaceDesc, err := interfaceParser.Parse()
			if err != nil {
				return nil, fmt.Errorf("failed to parse interface: %s", err)
			}
			interfaces = append(interfaces, interfaceDesc)
		}
	}
	fileDesc := &FileDesc{
		FileName:     f.fileName,
		PackageName:  node.Name.Name,
		RelativePath: "", // todo 未实现
		Imports:      importsDic,
		Structs:      structs,
		Interfaces:   interfaces,
	}
	return fileDesc, nil
}

func getServiceSpecs(file *ast.File) (specList []*ast.TypeSpec, err error) {
	specList = make([]*ast.TypeSpec, 0)
	for _, decl := range file.Decls {
		// 如果声明是一个类型声明
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// 遍历类型声明中的所有规格
			for _, spec := range genDecl.Specs {
				// 如果规格是一个类型规格
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					specList = append(specList, typeSpec)
				}
			}
		}
	}
	return specList, err
}

func (f *FileParser) parseImport(file *ast.File) (result map[string]*ImportDesc, err error) {
	result = make(map[string]*ImportDesc)
	ast.Inspect(file, func(n ast.Node) bool {
		if importSpec, ok := n.(*ast.ImportSpec); ok {
			path := strings.Trim(importSpec.Path.Value, "\"")
			// 如果有别名，使用别名作为包名
			name := ""
			hasAlias := importSpec.Name != nil
			if importSpec.Name != nil {
				name = importSpec.Name.Name
			} else {
				// 否则，使用路径的最后一个元素作为包名
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
