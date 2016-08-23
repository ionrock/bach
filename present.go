package bach

import (
	"errors"
	"fmt"
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

func (sm ServiceMap) Get(name string) (*Service, error) {
	s, ok := sm.Services[name]

	if ok {
		return s, nil
	}

	return nil, errors.New("Missing service")
}

func LocalConfig(name string) *memberlist.Config {
	c := memberlist.DefaultLocalConfig()
	c.Name = fmt.Sprintf("%s-%s", name, c.Name)
	return c
}

func InitializeMembership(config *memberlist.Config, memberIp string) ServiceMap {
	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	if memberIp != "" {
		_, err := list.Join([]string{memberIp})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}

	sm := make(map[string]*Service)

	// Ask for members of the cluster
	for _, member := range list.Members() {
		parts := strings.Split(member.Name, "-")

		_, ok := sm[parts[0]]
		if !ok {
			sm[parts[0]] = &Service{Name: parts[0]}

		}
		sm[parts[0]].Add(member)
	}

	return ServiceMap{
		Services:    sm,
		ServiceList: list,
	}
}
