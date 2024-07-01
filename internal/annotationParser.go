package internal

import (
	"strconv"
	"strings"
)

type AnnotationParser interface {
	Parse(comments []string) map[string]*Annotation
}

func getAnnotationParser(mode AnnotationMode) AnnotationParser {
	switch mode {
	case AnnotationModeMap:
		return &MapAnnotationParser{}
	default:
		return &ArrayAnnotationParser{}
	}
}

type ArrayAnnotationParser struct{}

func (a *ArrayAnnotationParser) Parse(comments []string) map[string]*Annotation {
	annotations := make(map[string]*Annotation)
	for _, comment := range comments {
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
	}
	return annotations
}

type MapAnnotationParser struct{}

func (a *MapAnnotationParser) Parse(comments []string) map[string]*Annotation {
	annotations := make(map[string]*Annotation)
	for _, comment := range comments {
		if strings.Contains(comment, "(") {
			// get annotation name
			name := comment[:strings.Index(comment, "(")]
			// get attribute
			attributeStr := comment[strings.Index(comment, "(")+1 : strings.LastIndex(comment, ")")]
			attributeSlice := strings.Split(attributeStr, ",")
			attribute := make(map[string]string)
			for _, item := range attributeSlice {
				// get attribute name and value
				itemSlice := strings.Split(item, "=")
				if len(itemSlice) != 2 {
					continue
				}
				attributeName := strings.TrimSpace(itemSlice[0])
				attributeValue := strings.TrimSpace(itemSlice[1])
				if strings.HasPrefix(attributeValue, "\"") && strings.HasSuffix(attributeValue, "\"") {
					attributeValue = attributeValue[1 : len(attributeValue)-1]
				}
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
			annotation := &Annotation{Name: comment, Attributes: []map[string]string{}}
			annotations[comment] = annotation
		}
	}
	return annotations
}
