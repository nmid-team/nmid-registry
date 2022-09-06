package registry

import (
	"encoding/json"
	"fmt"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"nmid-registry/pkg/cluster"
	"sync"
)

type Registry struct {
	lock sync.Mutex

	servicem map[string]*Service // serviceid-env -> service
	cluster  cluster.Cluster
}

func NewRegistry(cls cluster.Cluster) *Registry {
	return &Registry{
		servicem: make(map[string]*Service),
		cluster:  cls,
	}
}

//Register a new service.
func (r *Registry) Register(c *bm.Context, arg *ArgRegister, ins *Instance) (err error) {
	var sc *Service

	key := smapKey(arg.ServiceId, arg.Env)
	r.lock.Lock()
	if sc, ok := r.servicem[key]; !ok {
		sc = NewService(arg)
		r.servicem[key] = sc
	}
	r.lock.Unlock()

	sc.Instances = append(sc.Instances, ins)

	serviceVal, err := json.Marshal(sc)
	if nil != err {
		return err
	}

	//put to etcd cluster
	r.cluster.Put(key, string(serviceVal))

	return
}

func (r *Registry) Renew(c *bm.Context, arg *ArgRenew) (ins *Instance, err error) {

	return
}

func (r *Registry) LogOff(c *bm.Context, arg *ArgLogOff) (err error) {

	return
}

func smapKey(serviceId, env string) string {
	return fmt.Sprintf("%s-%s", serviceId, env)
}
