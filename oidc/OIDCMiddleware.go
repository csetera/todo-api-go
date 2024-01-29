package oidc

import (
	"github.com/gin-gonic/gin"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
)

type OIDCMiddleware struct {
	authorizer *authorization.Authorizer[*oauth.IntrospectionContext]
}

func New(authorizer *authorization.Authorizer[*oauth.IntrospectionContext]) *OIDCMiddleware {
	return &OIDCMiddleware{
		authorizer: authorizer,
	}
}

func (mw *OIDCMiddleware) RequiresRole(role string) gin.HandlerFunc {
	check := authorization.WithRole(role)

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		inspectCtx, err := mw.authorizer.CheckAuthorization(c, token, check)

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
