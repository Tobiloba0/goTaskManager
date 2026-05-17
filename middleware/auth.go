package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "gotask/db"
    "gotask/models"
)

// 1. Change this line to accept the string and return a gin.HandlerFunc
func RequireAuth(jwtSecret string) gin.HandlerFunc {
    
    // 2. Return the actual execution function that Gin expects
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }
 
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
 
        claims, _ := token.Claims.(jwt.MapClaims)
        userID := uint(claims["sub"].(float64))
 
        var user models.User
        if err := db.DB.First(&user, userID).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            return
        }
 
        c.Set("currentUser", user)
        c.Next()
    }
}