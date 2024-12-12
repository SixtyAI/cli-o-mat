package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"

	"github.com/SixtyAI/cli-o-mat/util"
)

const (
	CantFindOrgPrefix       = 14
	ErrorLookingUpRoleParam = 15
	OrgPrefixEmpty          = 16
)

type Omat struct {
	Credentials *CredentialCache `yaml:"-"`

	AccountName        string `yaml:"-"`
	OrganizationPrefix string `yaml:"organizationPrefix"`
	Region             string `yaml:"region"`
	Environment        string `yaml:"environment"`
	DeployService      string `yaml:"deployService"`
	BuildAccountSlug   string `yaml:"buildAccountSlug"`
	DeployAccountSlug  string `yaml:"deployAccountSlug"`
}

func NewOmat() *Omat {
	return &Omat{
		AccountName:        "",
		OrganizationPrefix: "",
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
	omat.InitCredentials()
	omat.FetchOrgPrefix()

	return nil
}

func (omat *Omat) FetchOrgPrefix() error {
	ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
	roleParamName := "/omat/organization_prefix"

	roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(roleParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			util.Fatalf(CantFindOrgPrefix, "Couldn't find org prefix parameter: %s\n", roleParamName)
		}

		util.Fatalf(ErrorLookingUpRoleParam,
			"Error looking up org prefix parameter %s, got: %s\n", roleParamName, err.Error())
	}

	orgPrefix := aws.StringValue(roleParam.Parameter.Value)
	if orgPrefix == "" {
		util.Fatalf(OrgPrefixEmpty, "Paramater '%s' was empty.\n", roleParamName)
	}

	omat.OrganizationPrefix = orgPrefix

	return nil
}

func (omat *Omat) InitCredentials() {
	omat.Credentials = newCredentialCache(omat)
}

func (omat *Omat) Prefix() string {
	// TODO: This is wrong.  All wrong!
	return fmt.Sprintf("/%s/%s", omat.OrganizationPrefix, omat.Environment)
}
