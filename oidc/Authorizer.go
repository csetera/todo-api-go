package oidc

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type ZitadelParameters struct {
	// ZITADEL instance domain (in the form: <instance>.zitadel.cloud or <yourdomain>)
	Domain string `required:"true"`

	// Path to the key.json
	Key string `required:"true"`

	// Port the Zitadel server is listening on
	Port string `required:"true"`

	// Whether the Zitadel port is not using secure transport
	Insecure bool `default:"false"`
}

type Authorizer struct {
	zitadelAuthorizer *authorization.Authorizer[*oauth.IntrospectionContext]
}

func New() (*Authorizer, error) {
	ctx := context.Background()

	params := ZitadelParameters{}
	err := envconfig.Process("zitadel", &params)
	if err != nil {
		return nil, err
	}

	// Initiate the authorization by providing a zitadel configuration and a verifier.
	var z *zitadel.Zitadel
	if params.Insecure {
		z = zitadel.New(params.Domain, zitadel.WithInsecure(params.Port))
	} else {
		z = zitadel.New(params.Domain)
	}

	authZ, err := authorization.New(ctx, z, oauth.DefaultAuthorization(params.Key))
	if err != nil {
		return nil, err
	}

	return &Authorizer{
		zitadelAuthorizer: authZ,
	}, nil
}

func (authz *Authorizer) RequiresRole(role string) gin.HandlerFunc {
	check := authorization.WithRole(role)

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		inspectCtx, err := authz.zitadelAuthorizer.CheckAuthorization(c, token, check)

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("introspection", inspectCtx)

		c.Next()
	}
}
