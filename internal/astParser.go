package internal

import (
	"fmt"
	"go/ast"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// getModuleName 获取模块名
func getModuleName() string {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error finding module root:", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getFullPackageName 获取当前文件所在完整包名
func getFullPackageName(moduleName string, filePath string) string {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error finding module root:", err)
		return ""
	}
	moduleRoot := strings.TrimSpace(string(output))

	// 计算文件路径相对于模块根目录的路径
	relativePath, err := filepath.Rel(moduleRoot, absolutePath)
	if err != nil {
		fmt.Println("Error calculating relative path:", err)
		return ""
	}
	// 将相对路径转换为包的导入路径
	dirPath := filepath.Dir(relativePath)
	importPath := filepath.ToSlash(dirPath)

	fullPackageName := moduleName + "/" + importPath // 模块名 + 相对路径
	return fullPackageName
}

func parseAtComments(commentGroup *ast.CommentGroup) (comments []string) {
	comments = make([]string, 0)
	if commentGroup != nil {
		prefix := "// " + AnnotationPrefix
		for _, com := range commentGroup.List {
			if strings.HasPrefix(com.Text, prefix) {
				commentText := strings.TrimPrefix(com.Text, prefix)
				comments = append(comments, commentText)
			}
		}
	}
	return comments
}

func parseDescription(name string, commentGroup *ast.CommentGroup) (description string) {
	if commentGroup == nil {
		return ""
	}
	description = ""
	prefix := "// " + name
	for _, com := range commentGroup.List {
		if strings.HasPrefix(com.Text, prefix) {
			// 取前缀之后的内容 去掉前后空格
			description = strings.TrimSpace(strings.TrimPrefix(com.Text, prefix))
			break
		}
	}
	return description

}

func parseAnnotation(comments []string, mode AnnotationMode) (annotations map[string]*Annotation) {
	annotations = make(map[string]*Annotation)
	for _, comment := range comments {
		if mode == AnnotationModeArray {
			commentSlice := splitComment(comment)
			if len(commentSlice) == 0 {
				continue
			}
			name := commentSlice[0]
			attribute := make(map[string]string)
			for i := 1; i < len(commentSlice); i++ {
				attribute[strconv.Itoa(i-1)] = commentSlice[i]
			}
			if _, ok := annotations[name]; ok {
				if len(attribute) > 0 {
					annotations[name].Attributes = append(annotations[name].Attributes, attribute)
				}
			} else {
				annotation := &Annotation{Name: name, Attributes: []map[string]string{}}
				if len(attribute) > 0 {
					annotation.Attributes = append(annotation.Attributes, attribute)
				}
				annotations[name] = annotation
			}
		} else {
			// map mode
			if strings.Contains(comment, "(") {
				// 获取注解名
				name := comment[:strings.Index(comment, "(")]
				// 获取注解属性
				attributeStr := comment[strings.Index(comment, "(")+1 : strings.LastIndex(comment, ")")]
				attributeSlice := strings.Split(attributeStr, ",")
				attribute := make(map[string]string)
				for _, item := range attributeSlice {
					// 获取属性名和属性值
					itemSlice := strings.Split(item, "=")
					if len(itemSlice) != 2 {
						continue
					}
					attributeName := strings.TrimSpace(itemSlice[0])
					attributeValue := strings.TrimSpace(itemSlice[1])
					attribute[attributeName] = attributeValue
				}
				if _, ok := annotations[name]; ok {
					if len(attribute) > 0 {
						annotations[name].Attributes = append(annotations[name].Attributes, attribute)
					}
				} else {
					annotation := &Annotation{Name: name, Attributes: []map[string]string{}}
					if len(attribute) > 0 {
						annotation.Attributes = append(annotation.Attributes, attribute)
					}
					annotations[name] = annotation
				}
			} else {
				// 不包含"("的注解
				annotation := &Annotation{Name: comment, Attributes: []map[string]string{}}
				annotations[comment] = annotation
			}

		}
	}
	return annotations
}

func splitComment(comment string) []string {
	re := regexp.MustCompile("[\\s　]") // 使用半角空格和全角空格作为分隔符
	commentSlice := re.Split(comment, -1)
	var filteredSlice []string
	for _, item := range commentSlice {
		trimmedItem := strings.TrimSpace(item)
		if trimmedItem != "" {
			filteredSlice = append(filteredSlice, trimmedItem)
		}
	}
	return filteredSlice
}

func parseField(field *ast.Field) (fieldDesc *Field, err error) {
	fieldDesc = &Field{}
	if field.Names != nil || len(field.Names) > 0 {
		fieldDesc.Name = field.Names[0].Name
	}
	fieldDesc.DataType = exprToString(field.Type)
	if strings.HasPrefix(fieldDesc.DataType, "*") {
		fieldDesc.RealDataType = fieldDesc.DataType[1:]
		fieldDesc.IsPtr = true
	} else {
		fieldDesc.RealDataType = fieldDesc.DataType
		fieldDesc.IsPtr = false
	}

	packageName := ""
	if ident, ok := field.Type.(*ast.Ident); ok {
		packageName = ident.Name
	} else if se, ok := field.Type.(*ast.SelectorExpr); ok {
		if id, ok := se.X.(*ast.Ident); ok {
			packageName = id.Name
		}
	} else if se2, ok := field.Type.(*ast.StarExpr); ok {
		if selExpr, ok := se2.X.(*ast.SelectorExpr); ok {
			if id, ok := selExpr.X.(*ast.Ident); ok {
				packageName = id.Name
			}
		}
	}
	fieldDesc.PackageName = packageName
	return fieldDesc, err
}

func exprToString(expr ast.Expr) string {

	switch t := expr.(type) {
	case *ast.Ident:
		// 标识符
		return t.Name
	case *ast.SelectorExpr:
		// 选择器
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		//  指针类型
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		// 数组或切片类型
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		// map类型
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.StructType:
		// 结构体类型
		var fields []string
		for _, field := range t.Fields.List {
			var names []string
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
			fields = append(fields, strings.Join(names, ", ")+" "+exprToString(field.Type))
		}
		return "struct{" + strings.Join(fields, "; ") + "}"
	case *ast.InterfaceType:
		// 接口类型
		if t.Methods == nil || len(t.Methods.List) == 0 {
			return "interface{}"
		}
		var methods []string
		for _, method := range t.Methods.List {
			var names []string
			for _, name := range method.Names {
				names = append(names, name.Name)
			}
			methods = append(methods, strings.Join(names, ", ")+" "+exprToString(method.Type))
		}
		return "interface{" + strings.Join(methods, "; ") + "}"
	default:
		return fmt.Sprintf("%T", t)
	}
}
