package middlewares

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/anggadarkprince/crud-employee-go/models"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/golang-jwt/jwt/v5"
)

// Context key for storing user data
type contextKey string

const userContextKey contextKey = "user"

// Auth holds dependencies for middleware
type Auth struct {
	UserRepository *repositories.UserRepository
	SecretKey      string
}

func (c *Auth) GetAuthToken(r *http.Request) string {
	var tokenString string
	var cookieName = os.Getenv("COOKIE_NAME")
	if cookieName == "" {
		cookieName = "auth_token"
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		// Cookie not found, try Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return ""
		}

		// Extract token from "Bearer <token>"
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) {
			return ""
		}
		tokenString = authHeader[len(bearerPrefix):]
	} else {
		tokenString = cookie.Value
	}
	return tokenString
}

// AuthMiddleware protects routes - redirects to login if not authenticated
func (c *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT token from cookie or Header
		authToken := c.GetAuthToken(r)

		if authToken == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate JWT token
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(c.SecretKey), nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get user ID from "sub" claim
		var userID int
		switch v := claims["sub"].(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		case int64:
			userID = int(v)
		case string:
			// Sometimes sub is stored as string
			parsed, err := strconv.Atoi(v)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			userID = parsed
		default:
			// fmt.Printf("Unexpected sub type: %T, value: %v\n", claims["sub"], claims["sub"])
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Query user from database
		user, err := c.UserRepository.GetById(r.Context(), userID)
		if err != nil {
			// User not found or database error
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Store user in context for later retrieval
		ctx := context.WithValue(r.Context(), userContextKey, user)

		// Pass the new context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GuestMiddleware for login/register pages - redirects to dashboard if already authenticated
func (c *Auth) GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := c.GetAuthToken(r)
		if authToken == "" {
			// No token, user is guest
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(c.SecretKey), nil
		})

		if err != nil || !token.Valid {
			// Invalid token, user is guest
			next.ServeHTTP(w, r)
			return
		}

		// User is authenticated, redirect to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})
}

// Call this anywhere in your handlers after AuthMiddleware
func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
