package internal

import (
	"go/ast"
	"strings"
	"unicode"
)

type StructParser struct {
	serviceName string
	file        *ast.File
	typeSpec    *ast.TypeSpec
	genDecl     *ast.GenDecl
	fileImports map[string]*ImportDesc
}

func NewStructParser(serviceName string,
	typeSpec *ast.TypeSpec,
	genDecl *ast.GenDecl,
	file *ast.File,
	fileImports map[string]*ImportDesc) *StructParser {
	return &StructParser{serviceName: serviceName,
		typeSpec:    typeSpec,
		genDecl:     genDecl,
		file:        file,
		fileImports: fileImports}
}

func (s *StructParser) Parse() (*StructDesc, error) {
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
	for _, f := range funcList {
		methodDesc, err := s.parserMethod(f)
		if err != nil {
			return nil, err
		}
		methods = append(methods, methodDesc)
	}
	sDesc := &StructDesc{
		Name:        s.serviceName,
		Description: description,
		Methods:     methods,
		Imports:     s.parserImports(methods),
		Comments:    comments,
		Annotations: getAnnotationParser(CurrentAnnotationMode).Parse(comments),
	}
	return sDesc, nil
}

func (s *StructParser) getFuncList() ([]*ast.FuncDecl, error) {
	list := make([]*ast.FuncDecl, 0)
	for _, decl := range s.file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if unicode.IsUpper(rune(funcDecl.Name.Name[0])) && funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						if ident.Name == s.serviceName && funcDecl.Doc != nil && strings.Contains(funcDecl.Doc.Text(), AnnotationPrefix) {
							list = append(list, funcDecl)
						}
					}
				}
			}
		}
	}
	return list, nil
}

func (s *StructParser) parserMethod(method *ast.FuncDecl) (methodDesc *MethodDesc, err error) {
	methodDesc = &MethodDesc{}
	methodDesc.Name = method.Name.Name
	// params
	params := make([]*Field, 0)
	if method.Type.Params != nil {
		for _, param := range method.Type.Params.List {
			field, err := parseField(param)
			if err != nil {
				return nil, err
			}
			params = append(params, field)
		}
	}
	methodDesc.Params = params
	// results
	results := make([]*Field, 0)
	if method.Type.Results != nil {
		for _, result := range method.Type.Results.List {
			field, err := parseField(result)
			if err != nil {
				return nil, err
			}
			results = append(results, field)
		}
	}
	methodDesc.Results = results
	// comment
	methodDesc.Comments = parseAtComments(method.Doc)
	methodDesc.Description = parseDescription(methodDesc.Name, method.Doc)
	methodDesc.Annotations = getAnnotationParser(CurrentAnnotationMode).Parse(methodDesc.Comments)
	return methodDesc, err
}

func (s *StructParser) parserImports(methods []*MethodDesc) (imports map[string]*ImportDesc) {
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
