package CSP

import "sync"

type Singleton struct {
	name string
}

var instance *Singleton
var onice sync.Once

func GetSingleton(name string) *Singleton {
	onice.Do(func() {
		instance = &Singleton{name: name}
	})
	return instance
}
