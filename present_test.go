package bach

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestClusterSetup(t *testing.T) {
	fooConfig := LocalConfig("foo")
	barConfig := LocalConfig("bar")
	barConfig.BindPort++

	// no ip just yet
	sm := InitializeMembership(fooConfig, "")

	s, err := sm.Get("foo")
	if err != nil {
		log.Info(len(sm.ServiceList.Members()))
		t.Error("Foo not connected")
	}

	assert.Equal(t, len(s.Hosts), 1)

	ip := fmt.Sprintf("%s:%d", sm.Ip(), fooConfig.BindPort)
	sm = InitializeMembership(barConfig, ip)

	assert.Len(t, sm.ServiceList.Members(), 2)
}
