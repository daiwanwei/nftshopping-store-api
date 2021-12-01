package casbins

import (
	"github.com/casbin/casbin"
	"nftshopping-store-api/pkg/config"
)

var (
	enforcerInstance *casbin.Enforcer
)

func GetEnforcer() (instance *casbin.Enforcer, err error) {
	if enforcerInstance == nil {
		instance, err = newEnforcer()
		if err != nil {
			return nil, err
		}
		enforcerInstance = instance
	}
	return enforcerInstance, nil
}

func newEnforcer() (instance *casbin.Enforcer, err error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	casbinConfig := c.Casbin
	enforcer, err := casbin.NewEnforcerSafe(casbinConfig.Model, casbinConfig.Policy)
	if err != nil {
		return
	}
	return enforcer, nil
}
