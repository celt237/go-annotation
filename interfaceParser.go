package go_annotation

import (
	"fmt"
	"go/ast"
)

type InterfaceParser struct {
	typeSpec      *ast.TypeSpec
	genDecl       *ast.GenDecl
	interfaceSpec *ast.InterfaceType
	serviceName   string
	fileImports   map[string]*ImportDesc
}

func NewInterfaceParser(
	serviceName string,
	typeSpec *ast.TypeSpec,
	genDecl *ast.GenDecl,
	fileImports map[string]*ImportDesc) *InterfaceParser {
	return &InterfaceParser{
		serviceName:   serviceName,
		typeSpec:      typeSpec,
		genDecl:       genDecl,
		interfaceSpec: typeSpec.Type.(*ast.InterfaceType),
		fileImports:   fileImports,
	}
}

func (s *InterfaceParser) Parse() (*InterfaceDesc, error) {
	comments := parseAtComments(s.genDecl.Doc)
	description := parseDescription(s.serviceName, s.genDecl.Doc)
	funcList, err := s.getFuncList()
	if err != nil {
		return nil, err
	}
	if funcList == nil || len(funcList) == 0 {
		return nil, nil
	}
	methods := make([]*MethodDesc, 0)
	for _, method := range funcList {
		methodDesc, err := s.parserMethod(method)
		if err != nil {
			return nil, err
		}
		methods = append(methods, methodDesc)
	}
	sDesc := &InterfaceDesc{
		Name:        s.serviceName,
		Description: description,
		Methods:     methods,
		Imports:     s.parserImports(methods),
		Comments:    comments,
		Annotations: getAnnotationParser(CurrentAnnotationMode).Parse(comments),
	}
	return sDesc, nil
}

func (s *InterfaceParser) getFuncList() ([]*ast.Field, error) {
	list := make([]*ast.Field, 0)
	for _, method := range s.interfaceSpec.Methods.List {
		list = append(list, method)
	}
	return list, nil
}

func (s *InterfaceParser) parserMethod(method *ast.Field) (methodDesc *MethodDesc, err error) {
	methodDesc = &MethodDesc{
		Comments: make([]string, 0),
		Params:   make([]*Field, 0),
		Results:  make([]*Field, 0),
	}
	// method name
	if method.Names == nil || len(method.Names) == 0 {
		err = fmt.Errorf("method name is empty")
		return
	}
	if len(method.Names) > 1 {
		err = fmt.Errorf("method name is not unique")
		return
	}
	methodDesc.Name = method.Names[0].Name
	// funcType
	if funcType, ok := method.Type.(*ast.FuncType); ok {
		// params
		if funcType.Params != nil {
			for _, param := range funcType.Params.List {
				field, err := parseField(param)
				if err != nil {
					return nil, err
				}
				methodDesc.Params = append(methodDesc.Params, field)
			}
		}
		// results
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				field, err := parseField(result)
				if err != nil {
					return nil, err
				}
				methodDesc.Results = append(methodDesc.Results, field)
			}
		}
		// comment
		methodDesc.Comments = parseAtComments(method.Doc)
		methodDesc.Description = parseDescription(methodDesc.Name, method.Doc)
		methodDesc.Annotations = getAnnotationParser(CurrentAnnotationMode).Parse(methodDesc.Comments)
		return methodDesc, err
	} else {
		err = fmt.Errorf("method type is not funcType")
		return
	}
}

func (s *InterfaceParser) parserImports(methods []*MethodDesc) (imports map[string]*ImportDesc) {
	imports = make(map[string]*ImportDesc)
	fields := make([]*Field, 0)
	for _, method := range methods {
		for _, param := range method.Params {
			fields = append(fields, param)
		}
		for _, result := range method.Results {
			fields = append(fields, result)
		}
	}
	for _, field := range fields {
		if imp, ok := s.fileImports[field.PackageName]; ok {
			imports[field.PackageName] = imp
		}
	}
	return imports
}
