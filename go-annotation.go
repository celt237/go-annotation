package go_annotation

// GetFileDesc 获取文件描述
// fileName: 文件名
func GetFileDesc(fileName string) (*FileDesc, error) {
	return GetFileParser(fileName).Parse()
}

// GetFilesDescList 获取文件描述列表
// directory: 目录
func GetFilesDescList(directory string) ([]*FileDesc, error) {
	var filesDesc []*FileDesc
	// 读取目录下的所有文件
	fileNames, err := GetFileNames(directory)
	if err != nil {
		return nil, err
	}
	for _, fileName := range fileNames {
		fileDesc, err := GetFileDesc(fileName)
		if err != nil {
			return nil, err
		}
		filesDesc = append(filesDesc, fileDesc)
	}
	return filesDesc, nil
}
