package main

import (
	"github.com/rickferrdev/gongo-simple-auth/config"
	"github.com/rickferrdev/gongo-simple-auth/internal/auth"
	"github.com/rickferrdev/gongo-simple-auth/internal/me"
	"github.com/rickferrdev/gongo-simple-auth/internal/middlewares"
	"github.com/rickferrdev/gongo-simple-auth/internal/platform"
	"github.com/rickferrdev/gongo-simple-auth/internal/server"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	env, err := config.NewEnvr()
	if err != nil {
		panic(err)
	}

	engine := server.NewServer()
	group := server.NewGroup(engine)
	authGroup := group.Group("/auth")
	usersGroup := group.Group("/users")
	client, err := config.NewMongo(env.GongoMongoURI)
	if err != nil {
		panic(err)
	}

	storage, err := auth.NewAuthStorage(client)
	if err != nil {
		panic(err)
	}
	tokenizer := platform.NewTokenizer(env.GongoJWTSecret)
	hasher := platform.NewHasher(bcrypt.DefaultCost)
	service := auth.NewAuthService(storage, tokenizer, hasher)

	auth.NewAuthHandler(authGroup, service)
	middlewares.NewGuardMiddleware(usersGroup, tokenizer)
	me.NewMeHandler(usersGroup)

	if err := engine.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
