package telegram

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LocalizationFile is a locale file for the bot
type LocalizationFile struct {
	Version int `yaml:"version"`
	Strings struct {
		GENERALABORT      string `yaml:"GENERAL_ABORT"`
		GETMEDIAGETNAME   string `yaml:"GETMEDIA_GET_NAME"`
		DUPLICATESHEADER  string `yaml:"DUPLICATES_HEADER"`
		DUPLICATESFOOTER  string `yaml:"DUPLICATES_FOOTER"`
		DUPLICATESCANCEL  string `yaml:"DUPLICATES_CANCEL"`
		ISMOVIEGETTYPE    string `yaml:"ISMOVIE_GET_TYPE"`
		CONFIRMMEDIAASK   string `yaml:"CONFIRM_MEDIA_ASK"`
		CONFIRMMEDIALEAVE string `yaml:"CONFIRM_MEDIA_LEAVE"`
	} `yaml:"strings"`
	Proceed struct {
		No  []string `yaml:"no"`
		Yes []string `yaml:"yes"`
	} `yaml:"proceed"`
}

// LoadLocale loads a locale file
func LoadLocale(locale string) (*LocalizationFile, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get workdir: %v", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(wd, fmt.Sprintf("localization/%s.yaml", locale)))
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("locale '%s' not found: %v", locale, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to read locale file: %v", err)
	}

	var l *LocalizationFile
	err = yaml.Unmarshal(b, &l)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal localization file: %v", err)
	}

	return l, nil
}
