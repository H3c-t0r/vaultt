package aws

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	awsauth "github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/hashicorp/vault/command/agent/auth"
)

const (
	typeEC2          = "ec2"
	typeIAM          = "iam"
	identityEndpoint = "http://169.254.169.254/latest/dynamic/instance-identity"
)

type awsMethod struct {
	logger       hclog.Logger
	authType     string
	nonce        string
	mountPath    string
	role         string
	headerValue  string
	accessKey    string
	secretKey    string
	sessionToken string
	watchCh      chan string
	stopCh       chan struct{}
	doneCh       chan struct{}
}

func NewAWSAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &awsMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
		watchCh:   make(chan string),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing 'type' value")
	}
	a.authType, ok = typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	a.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	switch {
	case a.role == "":
		return nil, errors.New("'role' value is empty")
	case a.authType == "":
		return nil, errors.New("'type' value is empty")
	case a.authType != typeEC2 && a.authType != typeIAM:
		return nil, errors.New("'type' value is invalid")
	}

	accessKeyRaw, ok := conf.Config["access_key"]
	if ok {
		a.accessKey, ok = accessKeyRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'access_key' value into string")
		}
	}

	secretKeyRaw, ok := conf.Config["secret_key"]
	if ok {
		a.secretKey, ok = secretKeyRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'secret_key' value into string")
		}
	}

	sessionTokenRaw, ok := conf.Config["session_token"]
	if ok {
		a.sessionToken, ok = sessionTokenRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'session_token' value into string")
		}
	}

	headerValueRaw, ok := conf.Config["header_value"]
	if ok {
		a.headerValue, ok = headerValueRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'header_value' value into string")
		}
	}

	return a, nil
}

func (a *awsMethod) Authenticate(client *api.Client) (*api.Secret, error) {
	a.logger.Trace("beginning authentication")

	data := make(map[string]interface{})

	switch a.authType {
	case typeEC2:
		client := cleanhttp.DefaultClient()

		// Fetch document
		{
			resp, err := client.Get(fmt.Sprintf("%s/document", identityEndpoint))
			if err != nil {
				return nil, errwrap.Wrapf("error fetching instance document: {{err}}", err)
			}
			if resp == nil {
				return nil, errors.New("empty response fetching instance document")
			}
			defer resp.Body.Close()
			doc, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errwrap.Wrapf("error reading instance document response body: {{err}}", err)
			}
			data["identity"] = base64.StdEncoding.EncodeToString(doc)
		}

		// Fetch signature
		{
			resp, err := client.Get(fmt.Sprintf("%s/signature", identityEndpoint))
			if err != nil {
				return nil, errwrap.Wrapf("error fetching instance document signature: {{err}}", err)
			}
			if resp == nil {
				return nil, errors.New("empty response fetching instance document signature")
			}
			defer resp.Body.Close()
			sig, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errwrap.Wrapf("error reading instance document signature response body: {{err}}", err)
			}
			data["signature"] = string(sig)
		}

		// Add the reauthentication value, if we have one
		if a.nonce == "" {
			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return nil, errwrap.Wrapf("error generating uuid for reauthentication value: {{err}}", err)
			}
			a.nonce = uuid
		}
		data["nonce"] = a.nonce

	default:
		var err error
		data, err = awsauth.GenerateLoginData(a.accessKey, a.secretKey, a.sessionToken, a.headerValue)
		if err != nil {
			return nil, errwrap.Wrapf("error creating login value: {{err}}", err)
		}
	}

	data["role"] = a.role

	secret, err := client.Logical().Write(fmt.Sprintf("%s/login", a.mountPath), data)
	if err != nil {
		return nil, errwrap.Wrapf("error logging in: {{err}}", err)
	}

	return secret, nil
}

func (a *awsMethod) NewCreds() chan struct{} {
	return nil
}

func (a *awsMethod) Shutdown() {
}
