package cmd

import (
	"fmt"
	"github.com/xetys/hetzner-kube/pkg/clustermanager"
	"encoding/json"
)

func GenerateKeyPairs(node clustermanager.Node, count int) []clustermanager.WgKeyPair {
	genKeyPairs := fmt.Sprintf(`echo "[" ;for i in {1..%d}; do pk=$(wg genkey); pubk=$(echo $pk | wg pubkey);echo "{\"private\":\"$pk\",\"public\":\"$pubk\"},"; done; echo "]";`, count)
	// gives an invalid JSON back
	o, err := wgContext.SSH().RunCmd(node, genKeyPairs)
	FatalOnError(err)
	o = o[0:len(o)-4] + "]"
	// now it's a valid json

	var keyPairs []clustermanager.WgKeyPair
	err = json.Unmarshal([]byte(o), &keyPairs)
	FatalOnError(err)

	return keyPairs
}

func GenerateWireguardConf(node clustermanager.Node, nodes []clustermanager.Node, clients []Client) string {
	headerTpl := `[Interface]
Address = %s
PrivateKey = %s
ListenPort = 51820
`
	// print header block
	output := fmt.Sprintf(headerTpl, node.PrivateIPAddress, node.WireGuardKeyPair.Private)

	// peer positions
	output = fmt.Sprintf("%s\n%s", output, renderPeerConfigs(nodes, &node, clients))

	return output
}

func GenerateClientConf(client Client, nodes []clustermanager.Node) string {
	headerTpl := `[Interface]
Address = %s
PrivateKey = %s
`
	// print header block
	output := fmt.Sprintf(headerTpl, client.Address, client.KeyPair.Private)

	// peer positions
	output = fmt.Sprintf("%s\n%s", output, renderPeerConfigs(nodes, nil, nil))

	return output
}
func renderPeerConfigs(nodes []clustermanager.Node, excludeNode *clustermanager.Node, clients []Client) string {
	peerTpl := `# %s
[Peer]
PublicKey = %s
AllowedIps = %s/32
Endpoint = %s:51820
`
	clientTpl := `[Peer]
PublicKey = %s
AllowedIps = %s
`
	var output string
	for _, peer := range nodes {
		if excludeNode != nil && peer.Name == excludeNode.Name {
			continue
		}

		output = fmt.Sprintf("%s\n%s",
			output,
			fmt.Sprintf(peerTpl, peer.Name, peer.WireGuardKeyPair.Public, peer.PrivateIPAddress, peer.IPAddress),
		)
	}

	for _, client := range clients {
		output = fmt.Sprintf("%s\n%s",
			output,
			fmt.Sprintf(clientTpl, client.KeyPair.Public, client.CIDR),
		)
	}
	return output
}

func (config *Config) SetupEncryptedNetwork(configNumber int, dryRun bool) error {
	// render a public/private key pair
	clients := config.Clients
	lenNodes := len(config.Nodes)
	keyPairs := GenerateKeyPairs(config.Nodes[0], lenNodes+len(clients))

	for i, keyPair := range keyPairs {
		if i < lenNodes {
			config.Nodes[i].WireGuardKeyPair = keyPair
		} else {
			config.Clients[i-lenNodes].KeyPair = keyPair
		}
	}

	// for each node, get specific IP and install it on node
	errChan := make(chan error)
	trueChan := make(chan bool)
	numProc := 0
	for _, node := range config.Nodes {
		numProc++
		go func(node clustermanager.Node) {
			fmt.Println(node.Name, "configure wireguard")
			wireGuardConf := GenerateWireguardConf(node, config.Nodes, config.Clients)
			configPath := fmt.Sprintf("/etc/wireguard/wg%d.conf", configNumber)
			if !dryRun {
				err := wgContext.SSH().WriteFile(node, configPath, wireGuardConf, false)
				if err != nil {
					errChan <- err
				}

				systemctlCommand := fmt.Sprintf("systemctl enable wg-quick@wg%d && systemctl restart wg-quick@wg%d", configNumber, configNumber)
				_, err = wgContext.SSH().RunCmd(node, systemctlCommand)

				if err != nil {
					errChan <- err
				}
			} else {
				fmt.Println(node.Name)
				fmt.Println(wireGuardConf)
			}

			fmt.Println(node.Name, "wireguard configured")
			trueChan <- true
		}(node)
	}

	waitOrError(trueChan, errChan, &numProc)
	return nil
}
