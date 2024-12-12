package config

import (
	"encoding/json"
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
	CantFindSSMParam       = 14
	ErrorLookingUpSSMParam = 15
	SSMParamEmpty          = 16
	CantParseSSMParam      = 17
)

type Omat struct {
	Credentials *CredentialCache `yaml:"-"`

	AccountName        string `yaml:"-"`
	OrganizationPrefix string `yaml:"organizationPrefix"`
	Region             string `yaml:"region"`
	Environment        string `yaml:"environment"`
	DeployService      string `yaml:"deployService"`
	ParamPrefix        string `yaml:"-"`
}

type accountInfoConfig struct {
	AccountID   string `json:"account_id"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	Purpose     string `json:"purpose"`
	Slug        string `json:"slug"`
}

func NewOmat(accountName string) *Omat {
	return &Omat{
		AccountName:        accountName,
		OrganizationPrefix: "",
		Region:             "us-east-1",
		Environment:        "development",
		DeployService:      "deployomat",
		ParamPrefix:        "",
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
	omat.FetchAccountInfo()

	return nil
}

func (omat *Omat) FetchOrgPrefix() {
	ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
	roleParamName := "/omat/organization_prefix"

	roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(roleParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			util.Fatalf(CantFindSSMParam, "Couldn't find org prefix parameter: %s\n", roleParamName)
		}

		util.Fatalf(ErrorLookingUpSSMParam,
			"Error looking up org prefix parameter %s, got: %s\n", roleParamName, err.Error())
	}

	orgPrefix := aws.StringValue(roleParam.Parameter.Value)
	if orgPrefix == "" {
		util.Fatalf(SSMParamEmpty, "Paramater '%s' was empty.\n", roleParamName)
	}

	omat.OrganizationPrefix = orgPrefix
}

func (omat *Omat) FetchAccountInfo() {
	ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
	infoParamName := "/omat/account_registry/" + omat.AccountName

	infoParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(infoParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			util.Fatalf(CantFindSSMParam, "Couldn't find account info parameter: %s\n", infoParamName)
		}

		util.Fatalf(ErrorLookingUpSSMParam,
			"Error looking up account info parameter %s, got: %s\n", infoParamName, err.Error())
	}

	accountInfo := aws.StringValue(infoParam.Parameter.Value)
	if accountInfo == "" {
		util.Fatalf(SSMParamEmpty, "Paramater '%s' was empty.\n", infoParamName)
	}

	var data accountInfoConfig
	if err = json.Unmarshal([]byte(accountInfo), &data); err != nil {
		util.Fatalf(CantParseSSMParam, "Couldn't parse account info parameter: %s\nGot: %s\n", infoParamName, accountInfo)
	}

	omat.ParamPrefix = data.Prefix
}

func (omat *Omat) InitCredentials() {
	omat.Credentials = newCredentialCache(omat)
}
