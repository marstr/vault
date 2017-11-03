package nomad

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Request token",
			},
		},

		Renew:  b.secretTokenRenew,
		Revoke: b.secretTokenRevoke,
	}
}

func (b *backend) secretTokenRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

func (b *backend) secretTokenRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, intErr := b.Client(req.Storage)
	if intErr != nil {
		return nil, intErr
	}

	tokenRaw, _ := req.Secret.InternalData["accessor_id"]

	_, err := c.ACLTokens().Delete(tokenRaw.(string), nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}