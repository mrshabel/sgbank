package handlers

// import (
// 	"github.com/gin-gonic/gin"
// )

// // SuccessResponse sends a successful response to the client
// type SuccessResponse struct {
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

// func (r *SuccessResponse) Send(c *gin.Context, status int) {
// 	if r.Message == "" && r.Data == "" {
// 		panic("response message or data should be present")
// 	}
// 	c.JSON(status, r)
// }

// func SendErrorResponse(c *gin.Context, status int, message string) {
// 	body := map[string]string{
// 		"message": message,
// 	}
// 	c.JSON(status, body)
// }
