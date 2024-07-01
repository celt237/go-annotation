package generate

type Config struct {
	// service文件所在目录 必传
	SourcePath string `yaml:"servicePath"`

	// 要生成的代码文件所在目录 必传
	GenFilePath string `yaml:"genFilePath"`

	// 模版文件地址
	TemplateFile string `yaml:"templateFile"`
}
