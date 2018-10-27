package analysis

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	prose "gopkg.in/jdkato/prose.v2"
)

// Model is a type: string: true
type Model map[string]map[string]bool

var model Model

func tokenize(msg string) ([]prose.Token, error) {
	doc, err := prose.NewDocument(msg)
	if err != nil {
		return nil, err
	}

	return doc.Tokens(), nil
}

// Train the various models
func Train() error {
	log.Printf("[analysis/train] training models")

	d, err := os.Getwd()
	if err != nil {
		log.Printf("[analysis/train] failed to get wd")
		return err
	}

	model = make(map[string]map[string]bool)

	datasetPath := filepath.Join(d, "dataset")
	err = filepath.Walk(datasetPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(path, datasetPath+"/")
		s := strings.Split(relPath, "/")
		modelName := s[0]
		file := strings.TrimSuffix(s[1], filepath.Ext(s[1]))

		if _, ok := model[modelName]; !ok {
			model[modelName] = make(map[string]bool)
		}

		log.Printf("[analysis/train] process: %s (model: %s, file: %s)", relPath, modelName, file)

		// only works for the proceed anlysis currently
		v, err := strconv.ParseBool(file)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		r := csv.NewReader(f)
		l, err := r.ReadAll()
		if err != nil {
			return err
		}

		for _, line := range l {
			s := line[0]
			// log.Printf("[analysis/train:parse] add line: %s to class %s", line[0], strconv.FormatBool(v))
			model[modelName][s] = v
		}

		if err != nil {
			log.Printf("[analysis/train] failed to parse file: %s", info.Name())
			os.Exit(1)
		}

		return nil
	})
	return err
}
