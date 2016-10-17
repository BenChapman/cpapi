package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BenChapman/platform/middleware"
	"github.com/BenChapman/platform/user"
	"github.com/Sirupsen/logrus"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type SessionClaims struct {
	SessionToken string `json:"session_token"`
	jwt.StandardClaims
}

func main() {
	rand.Seed(time.Now().UnixNano())

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		Debug:          true,
	})

	server := negroni.New()
	server.Use(negroni.NewRecovery())
	server.Use(middleware.NewLogger(logrus.New()))
	server.Use(c)
	// server.Use(negroni.NewStatic(http.Dir("public/app/")))

	db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/platform_n?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&user.User{})
	defer db.Close()

	userServer := user.NewUserServer(db)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("eey9iek9luoy6tah{thae7Eca1gayoo7quieNu^quaix,iequ9IekooG5eiP<ei9"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.HandleFunc("/api/v1/new_session", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, _ := generateRandomString(32)
		currentTime := time.Now()
		expiryTime := currentTime.Add(time.Hour * 3)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, SessionClaims{
			sessionToken,
			jwt.StandardClaims{
				ExpiresAt: expiryTime.Unix(),
				NotBefore: currentTime.Unix(),
				Issuer:    "org.coolestprojects.platform.session",
			},
		})
		tokenString, err := token.SignedString([]byte("eey9iek9luoy6tah{thae7Eca1gayoo7quieNu^quaix,iequ9IekooG5eiP<ei9"))
		fmt.Printf("%#v", err)
		w.Write([]byte(tokenString))
	})
	r.Handle("/api/v1/ninja/{id}", securedHandler(jwtMiddleware, userServer.HandleGet)).Methods("GET")
	r.Handle("/api/v1/ninja", securedHandler(jwtMiddleware, userServer.HandleGetAll)).Methods("GET")
	r.Handle("/api/v1/ninja", securedHandler(jwtMiddleware, userServer.HandlePost)).Methods("POST")
	r.Handle("/api/v1/ninja/{id}", securedHandler(jwtMiddleware, userServer.HandlePut)).Methods("PUT")
	r.Handle("/api/v1/ninja/{id}", securedHandler(jwtMiddleware, userServer.HandleDelete)).Methods("DELETE")
	server.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}

func securedHandler(jwtMiddleware *jwtmiddleware.JWTMiddleware, handlerToWrap http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handlerToWrap)),
	)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
