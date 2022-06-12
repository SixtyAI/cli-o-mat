package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

type Omat struct {
	OrganizationPrefix string `yaml:"organizationPrefix"`
	Region             string `yaml:"region"`
	Environment        string `yaml:"environment"`
}

func NewOmat() *Omat {
	return &Omat{
		OrganizationPrefix: "teak",
		Region:             "us-east-1",
		Environment:        "development",
	}
}

func (omat *Omat) LoadConfigFromFile(path string) error {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = yaml.Unmarshal(configData, omat); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (omat *Omat) LoadConfigFromEnv() {
	if organizationPrefix, wasSet := os.LookupEnv("OMAT_ORGANIZATION_PREFIX"); wasSet {
		omat.OrganizationPrefix = organizationPrefix
	}

	if region, wasSet := os.LookupEnv("OMAT_REGION"); wasSet {
		omat.Region = region
	}

	if environment, wasSet := os.LookupEnv("OMAT_ENVIRONMENT"); wasSet {
		omat.Environment = environment
	}
}

func (omat *Omat) Prefix() string {
	return fmt.Sprintf("/%s/%s", omat.OrganizationPrefix, omat.Environment)
}

func FindOmatConfig(dir string) (string, error) {
	result, err := filepath.Abs(dir)
	if err != nil {
		return result, errors.WithStack(err)
	}

	expectedPath := path.Join(result, ".omat.yml")
	if _, err = os.Stat(expectedPath); errors.Is(err, os.ErrNotExist) {
		if result != "/" {
			return FindOmatConfig(filepath.Dir(result))
		}

		return result, errors.New("Couldn't find .omat.yml anywhere")
	}

	return expectedPath, nil
}
