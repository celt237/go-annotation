package internal

type AnnotationMode string // 注解模式

const (
	AnnotationModeArray AnnotationMode = "array" // 数组注解模式
	AnnotationModeMap   AnnotationMode = "map"   // map注解模式
)

// Annotation 注解
type Annotation struct {
	Name       string              // 注解名称
	Attributes []map[string]string // 注解属性
}

// FileDesc  文件信息
type FileDesc struct {
	PackageName     string // 包名
	FullPackageName string // 完整包名
	FileName        string // 文件名
	//RelativePath string // todo 相对路径
	Imports    map[string]*ImportDesc
	Structs    []*StructDesc
	Interfaces []*InterfaceDesc
}

// ImportDesc  import信息
type ImportDesc struct {
	Name     string // 包名
	HasAlias bool   // 是否有别名
	Path     string // 路径
}

// StructDesc  结构体信息
type StructDesc struct {
	Name        string                 // 结构体名
	Imports     map[string]*ImportDesc // 导入信息
	Comments    []string               // 注释
	Annotations map[string]*Annotation // 注解
	//Fields      []*Field               // 字段 暂不支持
	Methods     []*MethodDesc // 方法
	Description string        // 描述
}

// InterfaceDesc  接口信息
type InterfaceDesc struct {
	Name        string                 // 接口名
	Imports     map[string]*ImportDesc // 导入信息
	Comments    []string               // 注释
	Annotations map[string]*Annotation // 注解
	Methods     []*MethodDesc          // 方法
	Description string                 // 描述
}

// MethodDesc  方法信息
type MethodDesc struct {
	Name        string                 // 方法名
	Description string                 // 描述
	Comments    []string               // 注释
	Annotations map[string]*Annotation // 注解
	Params      []*Field               // 参数
	Results     []*Field               // 返回值
}

// Field  字段信息（入参、出参）
type Field struct {
	Name         string //  字段名
	DataType     string // 字段类型
	PackageName  string // 包名
	RealDataType string // 真实类型 不含指针
	IsPtr        bool   // 是否是指针
}
