package encoding

import (
	"encoding/json"
	"errors"
	"gin/src/encoding/xml"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"slices"
)

var allowConfigurationFile = []string{".json", ".yml", ".yaml", ".xml"}

/*Parse
 * 解析配置文件或内容(".json", ".yml", ".yaml", ".xml")
 * @param  string $config 配置文件路径或内容
 */
func Parse(file string) (map[string]any, error) {
	var data map[string]any
	ext := filepath.Ext(file)
	if !slices.Contains(allowConfigurationFile, ext) {
		return nil, errors.New("not supporting file format: " + ext)
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	switch ext {
	case ".json":
		err = json.Unmarshal(bytes, &data)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(bytes, &data)
	case ".xml":
		err = xml.Unmarshal(bytes, &data)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
