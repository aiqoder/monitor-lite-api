package bind

import "github.com/gin-gonic/gin"

// Parse binds URI, query, form and JSON fields from the request.
func Parse(c *gin.Context, obj any) error {
	_ = c.ShouldBindUri(obj)
	return c.ShouldBind(obj)
}
