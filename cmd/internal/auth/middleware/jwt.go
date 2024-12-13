package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// JWT Secret Key (sebaiknya diambil dari environment variable)
var jwtSecretKey = []byte(os.Getenv("JWT_KEY"))

// JWTClaims struktur untuk claims JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userID, email string) (string, error) {
    claims := JWTClaims{
        UserID: userID, 
        Email:  email, 
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
            IssuedAt:  time.Now().Unix(), 
            Issuer:    "wafiuddin", 
        }, 
    }

    // Use the correct signing method with full details
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Ensure error is checked
    signedToken, err := token.SignedString(jwtSecretKey)
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }
    
    return signedToken, nil
}
func ValidateToken(tokenString string) (*JWTClaims, error) {
    // Use a key function to validate the signing method
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Check the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    // Type assert and validate
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header missing", 
            })
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, "Bearer ")
        if len(parts) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token format", 
            })
            c.Abort()
            return
        }

        tokenString := strings.TrimSpace(parts[1])
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Empty token", 
            })
            c.Abort()
            return
        }

        claims, err := ValidateToken(tokenString)
        if err != nil {
            log.Printf("Token validation error: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token", 
            })
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Next()
    }
}

func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}