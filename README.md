# wgctl

This tool simplifies management of wireguard nodes.


## Quick Tour

```

# this will create a new folder in the current directory and place a config.json with a name for usage. You must enter the folder for all further actions
wgctl new my-servers
cd my-servers

# this will import your ssh keys for connecting the nodes
wgctl import-ssh-keys

# this will add three nodes you are going to manage
wgctl add 42.42.42.1 10.0.0.1 id_rsa
wgctl add 42.42.42.2 10.0.0.2 id_rsa
wgctl add 42.42.42.3 10.0.0.3 id_rsa
# this step ensures your nodes will have wireguard install
wgctl install

# take a look into the configurations if you want
wgctl do -d
# setup the network
wgctl do
```

With these steps your servers should be connected. Now you could add an client config:

```
# adding a client 'me' to your list
wgctl client add me 10.0.1.1 10.0.1.0/24
# You must run this, to generate a certificate!
wgctl do
# generates a config for a client. You can share this one!
wgctl client gen-config me > wgme.conf
# connect to the VPN
wg-quick up ./wgme.conf
```


## Limitations

The install sub-command only works on Ubuntu nodes, yet.