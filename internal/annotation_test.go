package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func getInstanceFromJsonFile(fileName string) *FileDesc {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var obj FileDesc
	json.Unmarshal(byteValue, &obj)
	return &obj
}

func TestAnnotation(t *testing.T) {

	tests := []struct {
		name       string
		fileName   string
		mode       AnnotationMode
		wantResult *FileDesc
		wantErr    bool
	}{
		{
			name:       "数组模式单接口测试",
			fileName:   "test/data/arraymode/arraymode_single_interface.go",
			mode:       AnnotationModeArray,
			wantResult: getInstanceFromJsonFile("test/data/arraymode/arraymode_single_interface.json"),
			wantErr:    false,
		},
		{
			name:       "数组模式单结构体测试",
			fileName:   "test/data/arraymode/arraymode_single_struct.go",
			mode:       AnnotationModeArray,
			wantResult: getInstanceFromJsonFile("test/data/arraymode/arraymode_single_struct.json"),
			wantErr:    false,
		},
		{
			name:       "数组模式混合测试",
			fileName:   "test/data/arraymode/arraymode_mult.go",
			mode:       AnnotationModeArray,
			wantResult: getInstanceFromJsonFile("test/data/arraymode/arraymode_mult.json"),
			wantErr:    false,
		},
		{
			name:       "map模式单接口测试",
			fileName:   "test/data/mapmode/mapmode_single_interface.go",
			mode:       AnnotationModeMap,
			wantResult: getInstanceFromJsonFile("test/data/mapmode/mapmode_single_interface.json"),
			wantErr:    false,
		},
		{
			name:       "map模式单结构体测试",
			fileName:   "test/data/mapmode/mapmode_single_struct.go",
			mode:       AnnotationModeMap,
			wantResult: getInstanceFromJsonFile("test/data/mapmode/mapmode_single_struct.json"),
			wantErr:    false,
		},
		{
			name:       "map模式混合测试",
			fileName:   "test/data/mapmode/mapmode_mult.go",
			mode:       AnnotationModeMap,
			wantResult: getInstanceFromJsonFile("test/data/mapmode/mapmode_mult.json"),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CurrentAnnotationMode = tt.mode
			fileParser := GetFileParser(tt.fileName)
			fileDesc, err := fileParser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !deepCompare(fileDesc, tt.wantResult, "fileDesc") {
				t.Errorf("Parse() gotResult = %v, want %v", fileDesc, tt.wantResult)
			}
		})
	}

}

func deepCompare(a, b interface{}, fieldName string) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	equal := true
	tempEqual := true
	switch aVal.Kind() {
	case reflect.Ptr:
		return deepCompare(aVal.Elem().Interface(), bVal.Elem().Interface(), fieldName)
	case reflect.Struct:
		for i := 0; i < aVal.NumField(); i++ {
			tempEqual = deepCompare(aVal.Field(i).Interface(), bVal.Field(i).Interface(), fieldName+"."+aVal.Type().Field(i).Name)
			if !tempEqual {
				equal = false
			}
		}
	case reflect.Slice, reflect.Array:
		if aVal.Len() != bVal.Len() {
			fmt.Printf("Field %s slice Length is different: a = %v, b = %v\n", fieldName, aVal.Len(), bVal.Len())
			return false
		}
		for i := 0; i < aVal.Len(); i++ {
			tempEqual = deepCompare(aVal.Index(i).Interface(), bVal.Index(i).Interface(), fmt.Sprintf("%s[%d]", fieldName, i))
			if !tempEqual {
				equal = false
			}
		}
	case reflect.Map:
		if aVal.Len() != bVal.Len() {
			fmt.Printf("Field %s map Length is different: a = %v, b = %v\n", fieldName, aVal.Len(), bVal.Len())
			return false
		}
		for _, key := range aVal.MapKeys() {
			tempEqual = deepCompare(aVal.MapIndex(key).Interface(), bVal.MapIndex(key).Interface(), fmt.Sprintf("%s[%v]", fieldName, key))
			if !tempEqual {
				equal = false
			}
		}
	default:
		tempEqual = reflect.DeepEqual(aVal.Interface(), bVal.Interface())
		if !tempEqual {
			equal = false
			fmt.Printf("Field %s is different: a = %v, b = %v\n", fieldName, aVal.Interface(), bVal.Interface())
		}
	}
	return equal
}
