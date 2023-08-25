package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionDetails struct {
	Session     *session.Session
	Credentials *credentials.Credentials
	Config      *aws.Config
}

type CredentialCache struct {
	RootAWSConfig *aws.Config      `yaml:"-"`
	RootSession   *session.Session `yaml:"-"`
	// metaCreds     *SessionDetails
	byARN map[string]*SessionDetails `yaml:"-"`
}

func newCredentialCache(omat *Omat) *CredentialCache {
	return &CredentialCache{
		RootAWSConfig: aws.NewConfig().WithRegion(omat.Region),
		RootSession:   session.Must(session.NewSession()),

		byARN: make(map[string]*SessionDetails),
	}
}

func (cache *CredentialCache) ForARN(arn string) *SessionDetails {
	cred := cache.byARN[arn]
	if cred == nil {
		cache.byARN[arn] = &SessionDetails{}
		cred = cache.byARN[arn]
		cred.Session = session.Must(session.NewSession())
		cred.Credentials = stscreds.NewCredentials(cred.Session, arn)
		cred.Config = &aws.Config{Credentials: cred.Credentials}
	}

	return cred
}
