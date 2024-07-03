package go_annotation

import (
	"fmt"
	"go/ast"
	"os/exec"
	"path/filepath"
	"regexp"
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

	fullPackageName := moduleName + "/" + importPath
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
			// Remove the prefix and trim the spaces
			description = strings.TrimSpace(strings.TrimPrefix(com.Text, prefix))
			break
		}
	}
	return description

}

func splitComment(comment string) []string {
	re := regexp.MustCompile("[\\s　]") // Use half-width space and full-width space as separators
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
	if strings.Contains(fieldDesc.DataType, "*") {
		fieldDesc.RealDataType = strings.Replace(fieldDesc.DataType, "*", "", -1)
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
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		// pointer
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.StructType:
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
