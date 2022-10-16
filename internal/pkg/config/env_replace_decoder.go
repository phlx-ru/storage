package config

import (
	"bytes"
	"os"
	"regexp"

	"github.com/go-kratos/kratos/v2/config"
	"gopkg.in/yaml.v3"
)

const (
	envKeywordRegex              = `\${((.+?)(:.+?)*?)}`
	replacementFromRegexPosition = 0
	keyAndValueFromRegexPosition = 1
	keyFromRegexPosition         = 2
)

func EnvReplaceDecoder(kv *config.KeyValue, v map[string]any) error {
	configData := replaceEnv(kv.Value)
	return yaml.Unmarshal(configData, v)
}

func replaceEnv(configData []byte) []byte {
	for _, match := range regexp.MustCompile(envKeywordRegex).FindAllSubmatch(configData, -1) {
		key := string(match[keyFromRegexPosition])
		value := []byte(os.Getenv(key))
		if len(value) == 0 {
			value = bytes.TrimLeft(match[keyAndValueFromRegexPosition], key+":")
		}
		if len(value) > 0 {
			configData = bytes.Replace(configData, match[replacementFromRegexPosition], value, 1)
		}
	}

	return configData
}
