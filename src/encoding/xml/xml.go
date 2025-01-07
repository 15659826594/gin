package xml

import (
	"encoding/xml"
	"regexp"
	"strconv"
)

func Marshal(v any) ([]byte, error) {
	var buf []byte
	var err error
	if m, ok := v.(map[string]any); ok {
		doc := NewDocument()
		doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
		root := doc.CreateElement("root")
		xml2map(m, root)
		doc.Indent(4)
		//fmt.Println(doc.WriteTo(os.Stdout))
		buf, err = doc.WriteToBytes()
		if err != nil {
			return nil, err
		}
	} else {
		buf, err = xml.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func xml2map(m map[string]any, ele *Element) {
	for k, v := range m {
		child := ele.CreateElement(k)
		if k2, ok := v.(map[string]any); ok {
			xml2map(k2, child)
		} else {
			child.CreateText(v.(string))
		}
	}
}

func Unmarshal(data []byte, v any) error {
	var err error
	if m, ok := v.(*map[string]any); ok {
		doc := NewDocument()
		if err = doc.ReadFromBytes(data); err != nil {
			return err
		}
		root := doc.Root()
		if *m == nil {
			*m = make(map[string]any)
		}
		for _, ele := range root.ChildElements() {
			(*m)[ele.name()] = map2xml(ele)
		}
	} else {
		err = xml.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func map2xml(ele *Element) any {
	m := make(map[string]any)
	childs := ele.ChildElements()
	childNum := len(childs)
	if childNum == 0 {
		return strconvBase(ele.Text())
	} else {
		for _, child := range childs {
			m[child.name()] = map2xml(child)
		}
	}
	return m
}

// 字符串转基础数据类型
func strconvBase(str string) any {
	if str == "true" {
		return true
	}
	if str == "false" {
		return false
	}

	if isNumeric(str) {
		integer, err := strconv.Atoi(str)
		if err == nil {
			return integer
		}
		float1, err := strconv.ParseFloat(str, 32)
		if err == nil {
			return float1
		}
		float2, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return float2
		}
	}
	return str
}

func isNumeric(str string) bool {
	reg := regexp.MustCompile(`^\s*[+-]?((\d+(\.\d*)?)|(\.\d+))([eE][+-]?\d+)?\s*$`)
	return reg.MatchString(str)
}
