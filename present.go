package bach

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/memberlist"
)

func RunScript(script string) error {
	log.Debug("Running Script: ", script)

	if script != "" {
		cmd := NewCommand(script)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

type NodeInfo struct {
	Name string
}

type Service struct {
	Name  string
	Hosts []*memberlist.Node
}

func (s *Service) Add(n *memberlist.Node) {
	s.Hosts = append(s.Hosts, n)
}

func (s *Service) Ip(index int) string {
	return s.Hosts[index].Addr.String()
}

type ServiceMap struct {
	Config      *memberlist.Config
	ClusterAddr string
	Services    map[string]*Service
	ServiceList *memberlist.Memberlist
}

// Find the first IP in order to join the cluster
func (sm ServiceMap) Ip() string {
	for _, v := range sm.Services {
		if len(v.Hosts) > 0 {
			return v.Ip(0)
		}
	}
	return ""
}

func (sm *ServiceMap) Get(name string) (*Service, error) {
	s, ok := sm.Services[name]

	if ok {
		return s, nil
	}

	return nil, errors.New("Missing service")
}

func (sm *ServiceMap) Load() *ServiceMap {
	list, err := memberlist.Create(sm.Config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	sm.ServiceList = list
	return sm
}

func (sm *ServiceMap) Join() *ServiceMap {
	if sm.ServiceList == nil {
		sm.Load()
	}

	// Join an existing cluster by specifying at least one known member.
	if sm.ClusterAddr != "" {
		_, err := sm.ServiceList.Join([]string{sm.ClusterAddr})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}

		sm.Sync()
	}

	return sm
}

func (sm *ServiceMap) Leave() *ServiceMap {
	if sm.ServiceList != nil {
		sm.Leave()
	}
	return sm
}

func (sm *ServiceMap) Sync() *ServiceMap {
	if sm.ServiceList == nil {
		sm.Load()
	}

	m := make(map[string]*Service)

	// Ask for members of the cluster
	for _, member := range sm.ServiceList.Members() {
		parts := strings.Split(member.Name, "-")

		_, ok := m[parts[0]]
		if !ok {
			m[parts[0]] = &Service{Name: parts[0]}

		}
		m[parts[0]].Add(member)
	}

	sm.Services = m

	return sm
}

func (sm *ServiceMap) AsJson() {

}

func (sm *ServiceMap) CopyJsonTo(fh io.Writer) {
	services := make(map[string]string)
	localNode := sm.ServiceList.LocalNode()

	for _, n := range sm.Services {
		hosts := make([]string, len(n.Hosts))
		for i, h := range n.Hosts {
			if h != nil && h != localNode {
				hosts[i] = fmt.Sprintf("%s:%d", h.Addr.String(), h.Port)
			}
		}
		services[strings.ToUpper(n.Name)] = strings.Join(hosts, ", ")
	}

	doc := ServicesDocument{Services: services}
	err := json.NewEncoder(fh).Encode(doc)
	if err != nil {
		panic(err)
	}
}

func LocalConfig(name string) *memberlist.Config {
	c := memberlist.DefaultLocalConfig()
	c.Name = fmt.Sprintf("%s-%s", name, c.Name)
	return c
}

type ServiceDocument struct {
	Name      string `json:"NAME"`
	Addresses string `json:"ADDRESSES"`
}

type ServicesDocument struct {
	Services map[string]string `json:"BACH_SERVICES"`
}
