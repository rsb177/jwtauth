package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-cz/jwtauth"
	jwt "github.com/golang-jwt/jwt/v4"
)

var tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

func main() {
	fmt.Printf("Starting server on http://localhost:3333\n\n")

	// Generate JWT token for debugging purposes.
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})

	fmt.Printf("Try the following commands in new terminal window:\n")
	fmt.Printf("curl http://localhost:3333/\n")
	fmt.Printf("curl http://localhost:3333/admin\n")
	fmt.Printf("curl -H \"Authorization: BEARER %s\" http://localhost:3333/admin\n", tokenString)

	http.ListenAndServe(":3333", router())
}

func router() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("welcome to protected area (user_id=%v)\n", claims["user_id"])))
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous\n"))
		})
	})

	return r
}
