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
}

func NewStructParser(serviceName string, typeSpec *ast.TypeSpec, file *ast.File) *StructParser {
	return &StructParser{serviceName: serviceName, typeSpec: typeSpec, file: file}
}

func (s *StructParser) Parse() (sDesc *StructDesc, err error) {
	sDesc = &StructDesc{
		Name:    s.serviceName,
		Methods: make([]*MethodDesc, 0),
		Imports: make(map[string]*ImportDesc),
	}
	comments, err := getAtComments(s.typeSpec.Doc)
	sDesc.Comments = comments
	if err != nil {
		return nil, err
	}

	funcList, err := s.getFuncList()
	if err != nil {
		return nil, err
	}
	if funcList == nil || len(funcList) == 0 {
		return nil, nil
	}
	for _, f := range funcList {
		methodDesc, err := s.parserMethod(f)
		if err != nil {
			return nil, err
		}
		sDesc.Methods = append(sDesc.Methods, methodDesc)
	}
	sDesc.Annotations = parseAnnotation(comments, CurrentAnnotationMode)
	return sDesc, nil
}

func (s *StructParser) getFuncList() ([]*ast.FuncDecl, error) {
	list := make([]*ast.FuncDecl, 0)
	for _, decl := range s.file.Decls {
		// 如果声明是一个函数
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
	if method.Type.Params != nil && method.Type.Params.List != nil && len(method.Type.Params.List) > 0 {
		return nil, nil
	}
	params := make([]*Field, 0)
	for _, param := range method.Type.Params.List {
		field, err := parseField(param)
		if err != nil {
			return nil, err
		}
		params = append(params, field)
	}
	methodDesc.Params = params

	if method.Type.Results == nil || method.Type.Results.List == nil || len(method.Type.Results.List) == 0 {
		return nil, nil
	}
	results := make([]*Field, 0)
	for _, result := range method.Type.Results.List {
		field, err := parseField(result)
		if err != nil {
			return nil, err
		}
		results = append(results, field)
	}
	methodDesc.Results = results

	comments := make([]string, 0)
	if method.Doc != nil {
		for _, com := range method.Doc.List {
			atIndex := strings.Index(com.Text, "@")
			if atIndex != -1 {
				commentText := com.Text[atIndex:]
				comments = append(comments, commentText)
			}
		}
	}
	methodDesc.Comments = comments
	methodDesc.Annotations = parseAnnotation(comments, CurrentAnnotationMode)
	return methodDesc, err
}
