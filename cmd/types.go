package cmd

import (
	"github.com/xetys/hetzner-kube/pkg/clustermanager"
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type WgCtlContext struct {
	Config *Config
	ssh *clustermanager.SSHCommunicator
}

type Config struct {
	Name string `json:"name"`
	Nodes []clustermanager.Node `json:"nodes"`
	SSHKeys []clustermanager.SSHKey
}

type Client struct {
	Address string                   `json:"address"`
	KeyPair clustermanager.WgKeyPair `json:"key_pair"`
}

// LoadFromLocalDir tries to load the config from local dir.
// returns true if found and loaded, otherwise false
// false does not mean an error exists. If no config was found, it returns false and no error
func (ctx *WgCtlContext) LoadFromLocalDir() (bool, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	dir, err := filepath.Abs(pwd)
	if err != nil {
		return false, err
	}

	configPath := dir + "/config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false, nil
	}


	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return false, err
	}

	ctx.Config = &Config{}
	err = json.Unmarshal(configBytes, &ctx.Config)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ctx *WgCtlContext) SaveConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir, err := filepath.Abs(pwd)
	if err != nil {
		return err
	}

	configPath := dir + "/config.json"
	configBytes, err := json.Marshal(ctx.Config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, configBytes, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *WgCtlContext) IsConfigLoaded() bool {
	return ctx.Config != nil
}

func (ctx *WgCtlContext) SSH() *clustermanager.SSHCommunicator {
	if !ctx.IsConfigLoaded() {
		panic(errors.New("you can only run ssh with a loaded config"))
	}

	if ctx.ssh == nil {
		ctx.ssh = clustermanager.NewSSHCommunicator(ctx.Config.SSHKeys).(*clustermanager.SSHCommunicator)
	}

	return ctx.ssh
}
