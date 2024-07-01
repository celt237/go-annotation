package go_annotation

import "github.com/celt237/go-annotation/internal"

// GetFileDesc 获取文件描述
// fileName: 文件名
func GetFileDesc(fileName string) (*internal.FileDesc, error) {
	return internal.GetFileParser(fileName).Parse()
}

// GetFilesDescList 获取文件描述列表
// directory: 目录
func GetFilesDescList(directory string) ([]*internal.FileDesc, error) {
	var filesDesc []*internal.FileDesc
	// 读取目录下的所有文件
	fileNames, err := internal.GetFileNames(directory)
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
