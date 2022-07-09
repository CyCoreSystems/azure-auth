package smtp

import (
	"encoding/json"
	"fmt"

	"github.com/emersion/go-sasl"
)

// The XAUTH2 mechanism name.
const OAuthBearer = "XOAUTH2"

type OAuthBearerError struct {
	Status  string `json:"status"`
	Schemes string `json:"schemes"`
	Scope   string `json:"scope"`
}

type OAuthBearerOptions struct {
	Username string
	AccessToken    string
}

// Implements error
func (err *OAuthBearerError) Error() string {
	return fmt.Sprintf("OAUTHBEARER authentication error (%v)", err.Status)
}

type oauthBearerClient struct {
	OAuthBearerOptions
}

func (a *oauthBearerClient) Start() (mech string, ir []byte, err error) {
	mech = OAuthBearer

	str := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.Username, a.AccessToken)

	ir = []byte(str)
	return
}

func (a *oauthBearerClient) Next(challenge []byte) ([]byte, error) {
	authBearerErr := &OAuthBearerError{}
	if err := json.Unmarshal(challenge, authBearerErr); err != nil {
		return nil, err
	} else {
		return nil, authBearerErr
	}
}

// An implementation of the OAUTHBEARER authentication mechanism, as
// described in RFC 7628.
func NewOAuthBearerClient(opt *OAuthBearerOptions) sasl.Client {
	return &oauthBearerClient{*opt}
}

