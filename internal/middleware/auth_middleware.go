package middleware

// AuthMiddleware godoc
// @Security BearerAuth
// @Description JWT authentication middleware
// @Param Authorization header string true "Bearer {token}"
// func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "authorization header required",
// 			})
// 			return
// 		}

// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "invalid authorization header format",
// 			})
// 			return
// 		}

// 		doctorID, err := authService.ValidateToken(tokenParts[1])
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "invalid token: " + err.Error(),
// 			})
// 			return
// 		}

// 		c.Set("doctorID", doctorID)
// 		c.Next()
// 	}
// }
