package bach

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Transport interface {
	Cp(string, string)
	Run(string)
}

type SSHClient struct {
	User   string
	Host   string
	Noop   bool
	Config *ssh.ClientConfig
	conn   *ssh.Client
	agent  agent.Agent
}

func (client *SSHClient) Url() string {
	return fmt.Sprintf("%s@%s", client.User, client.Host)
}

func loadKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		// log.Fatalf("unable to read private key: %v", err)
		panic(err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		// log.Fatalf("unable to parse private key: %v", err)
		panic(err)
	}

	if signer == nil {
		log.Fatalf("The signer from %s is nil", path)
	}

	return ssh.PublicKeys(signer)
}

func findKeys() []ssh.AuthMethod {
	keys := []ssh.AuthMethod{
		loadKey("/Users/eric/.ssh/id_rsa"),
	}

	sshconfigdir := filepath.Join(os.Getenv("HOME"), ".ssh")
	log.Printf("Reading keys from %s", sshconfigdir)

	fns, err := ioutil.ReadDir(sshconfigdir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fn := range fns {
		name := fn.Name()
		if fn.IsDir() {
			continue
		}

		if strings.HasSuffix(name, ".pub") {
			continue
		}

		if strings.Contains(name, "id_rsa") {
			log.Printf("Found key: %s", name)
			keyfile := filepath.Join(sshconfigdir, name)
			keys = append(keys, loadKey(keyfile))
		}
	}

	return keys
}

func SSHAgent() agent.Agent {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return agent.NewClient(sshAgent)
	}
	return nil
}

func (client *SSHClient) LoadConfig() {
	var keys []ssh.AuthMethod

	ag := SSHAgent()
	if ag != nil {
		client.agent = ag
		keys = []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers)}
	} else {
		keys = findKeys()
	}

	client.Config = &ssh.ClientConfig{
		User: client.User,
		Auth: keys,
	}
}

func (client *SSHClient) connect() {
	if client.conn != nil {
		return
	}

	log.Print(client.Host)
	log.Printf("%#v", client.Config)

	conn, err := ssh.Dial("tcp", client.Host+":22", client.Config)
	if err != nil {
		log.Fatal(err)
	}

	if client.agent != nil {
		agent.ForwardToAgent(conn, client.agent)
	}

	client.conn = conn
}

func (client *SSHClient) Run(cmd string) error {
	client.connect()

	sess, err := client.conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	out, err := sess.Output(cmd)
	if err != nil {
		log.Printf("Error starting cmd: ", err)
		return err
	}

	if len(out) > 0 {
		return nil
	}

	for _, line := range bytes.Split(out, []byte("\n")) {
		if len(line) > 0 {
			log.Printf("%s", line)
		}
	}

	return nil
}

func (client *SSHClient) Cp(local string, remote string) {
	client.connect()
	sess, err := client.conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	if client.agent != nil {
		if err := agent.RequestAgentForwarding(sess); err != nil {
			log.Fatal(err)
		}
	}

	targetFile := filepath.Base(remote)

	src, err := os.Open(local)

	if err != nil {
		panic(err)
	}
	defer src.Close()

	stat, err := src.Stat()

	if err != nil {
		panic(err)
	}

	go func() {
		w, _ := sess.StdinPipe()
		defer w.Close()

		fmt.Fprintln(w, "C0644", stat.Size(), targetFile)
		io.Copy(w, src)
		fmt.Fprint(w, "\x00")
	}()

	cmd := fmt.Sprintf("/usr/bin/scp -t %s", targetFile)

	err = sess.Run(cmd)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	mv := fmt.Sprintf("sudo mv %s %s", targetFile, remote)
	// log.Printf("Doing: %s", mv)
	client.Run(mv)
}
