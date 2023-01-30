package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

type Omat struct {
	Credentials *CredentialCache

	OrganizationPrefix string `yaml:"organizationPrefix"`
	Region             string `yaml:"region"`
	Environment        string `yaml:"environment"`
	DeployService      string `yaml:"deployService"`
	BuildAccountSlug   string `yaml:"buildAccountSlug"`
	DeployAccountSlug  string `yaml:"deployAccountSlug"`
}

func NewOmat() *Omat {
	return &Omat{
		OrganizationPrefix: "teak",
		Region:             "us-east-1",
		Environment:        "development",
		DeployService:      "deployomat",
		BuildAccountSlug:   "ci-cd",
		DeployAccountSlug:  "workload",
	}
}

func findOmatConfig(dir string) (string, error) {
	result, err := filepath.Abs(dir)
	if err != nil {
		return result, errors.WithStack(err)
	}

	expectedPath := path.Join(result, ".omat.yml")
	if _, err = os.Stat(expectedPath); errors.Is(err, os.ErrNotExist) {
		if result != "/" {
			return findOmatConfig(filepath.Dir(result))
		}

		return result, errors.New("Couldn't find .omat.yml anywhere")
	}

	return expectedPath, nil
}

func (omat *Omat) loadConfigFromFile(path string) error {
	configData, err := os.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = yaml.Unmarshal(configData, omat); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (omat *Omat) loadConfigFromEnv() {
	if organizationPrefix, wasSet := os.LookupEnv("OMAT_ORGANIZATION_PREFIX"); wasSet {
		omat.OrganizationPrefix = organizationPrefix
	}

	if region, wasSet := os.LookupEnv("OMAT_REGION"); wasSet {
		omat.Region = region
	}

	if environment, wasSet := os.LookupEnv("OMAT_ENVIRONMENT"); wasSet {
		omat.Environment = environment
	}

	if deployService, wasSet := os.LookupEnv("OMAT_DEPLOY_SERVICE"); wasSet {
		omat.DeployService = deployService
	}

	if buildAccountSlug, wasSet := os.LookupEnv("OMAT_BUILD_ACCOUNT_SLUG"); wasSet {
		omat.BuildAccountSlug = buildAccountSlug
	}

	if deployAccountSlug, wasSet := os.LookupEnv("OMAT_DEPLOY_ACCOUNT_SLUG"); wasSet {
		omat.DeployAccountSlug = deployAccountSlug
	}
}

func (omat *Omat) LoadConfig() error {
	path, err := findOmatConfig(".")
	if err != nil {
		return errors.Wrap(err, "couldn't find config file")
	}

	if err = omat.loadConfigFromFile(path); err != nil {
		return errors.Wrap(err, "couldn't load config from file")
	}

	omat.loadConfigFromEnv()

	return nil
}

func (omat *Omat) InitCredentials() {
	omat.Credentials = newCredentialCache(omat)
}

func (omat *Omat) Prefix() string {
	return fmt.Sprintf("/%s/%s", omat.OrganizationPrefix, omat.Environment)
}
