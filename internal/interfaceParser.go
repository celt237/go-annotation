package internal

import (
	"fmt"
	"go/ast"
)

type InterfaceParser struct {
	typeSpec      *ast.TypeSpec
	interfaceSpec *ast.InterfaceType
	serviceName   string
}

func NewInterfaceParser(serviceName string, typeSpec *ast.TypeSpec) *InterfaceParser {
	return &InterfaceParser{serviceName: serviceName, typeSpec: typeSpec, interfaceSpec: typeSpec.Type.(*ast.InterfaceType)}
}

func (s *InterfaceParser) Parse() (*InterfaceDesc, error) {
	comments, err := getAtComments(s.typeSpec.Doc)
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
		Methods:     methods,
		Imports:     make(map[string]*ImportDesc),
		Comments:    comments,
		Annotations: parseAnnotation(comments, CurrentAnnotationMode),
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
	// 方法名
	if method.Names == nil || len(method.Names) == 0 {
		err = fmt.Errorf("method name is empty")
		return
	}
	if len(method.Names) > 1 {
		err = fmt.Errorf("method name is not unique")
		return
	}
	methodDesc.Name = method.Names[0].Name
	// 方法类型
	if funcType, ok := method.Type.(*ast.FuncType); ok {
		// 参数
		if funcType.Params != nil {
			for _, param := range funcType.Params.List {
				field, err := parseField(param)
				if err != nil {
					return nil, err
				}
				methodDesc.Params = append(methodDesc.Params, field)
			}
		}
		// 返回值
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				field, err := parseField(result)
				if err != nil {
					return nil, err
				}
				methodDesc.Results = append(methodDesc.Results, field)
			}
		}
		// 注释
		methodDesc.Comments, err = getAtComments(method.Doc)
		if err != nil {
			return nil, err
		}
		methodDesc.Annotations = parseAnnotation(methodDesc.Comments, CurrentAnnotationMode)
		return methodDesc, err
	} else {
		err = fmt.Errorf("method type is not funcType")
		return
	}
}
