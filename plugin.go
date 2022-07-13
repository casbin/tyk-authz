package main

import (
	"errors"
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/casbin/casbin/v2"
	"net/http"
)

var logger = log.Get()


func CasbinAuthz(rw http.ResponseWriter, r *http.Request) {
	logger.Info("Casbin Authz Plugin start processing")

	auth, err := parseConfigData(ctx.GetDefinition(r).ConfigData)
	if err != nil {
		logger.Error("Error parsing config data: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	username := r.Header.Get("username")
	path := r.URL.Path
	method := r.Method

	logger.Info("username: ", username)
	logger.Info("path: ", path)
	logger.Info("method: ", method)

	ok, err := auth.e.Enforce(username, path, method)
	if err != nil {
		logger.Error("Error checking permission: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	if !ok {
		logger.Error("Authorization failed")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
}

type Auth struct {
	e 		   *casbin.Enforcer
	ModelPath  string
	PolicyPath string
}

func init() {
	logger.Info("Initialising Casbin Authz Plugin")
}

func parseConfigData(confMap map[string]interface{}) (*Auth, error) {
	auth := &Auth{}

	authConfig := confMap["casbin_authz_plugin"].(map[string]interface{})

	key := "model_path"
	value := authConfig[key].(string)
	logger.Debug( key + value)
	if value == "" {
		return nil, errors.New(key + " not found")
	}
	auth.ModelPath = value

	key = "policy_path"
	value = authConfig[key].(string)
	logger.Debug( key + value)
	if value == "" {
		return nil, errors.New(key + " not found")
	}
	auth.PolicyPath = value

	e, err := casbin.NewEnforcer(auth.ModelPath, auth.PolicyPath)
	if err != nil {
		return nil, err
	}
	auth.e = e

	return auth, nil
}

func main() {}
