package main

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"time"
	"strings"
	"log"
	"context"
	"errors"
	"os"
)

type KeyNode struct {
	NodeNetwork string
	ContainerNetwork string
	AlreadyUsedIp string

}

type EtcdHelper struct {
	HeaderTimeoutPerRequest   time.Duration
	Client           client.Client
}

func (IpamS *IpamConfig) etcdConn() (EtcdConn EtcdHelper){

	var etcdServerList []string = strings.Split(IpamS.Ipam.Etcdcluster, ",")
	cli, err := client.New(client.Config{
		Endpoints:   etcdServerList,
		HeaderTimeoutPerRequest: 1* time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		os.Exit(-1)

	}

	return EtcdHelper{
		HeaderTimeoutPerRequest: 1 * time.Second,
		Client: cli,
		}


}

func IsKeyExist(Rang *map[string]string,Key string) error{

	if _, ok := (*Rang)[Key +"rangeStart"]; !ok {
		return errors.New("ETCD Lack rangeStart")
	}


	if _, ok := (*Rang)[Key +"rangeEnd"]; !ok {
		return errors.New("ETCD Lack rangeEnd")

	}
	return nil
}





func (Cli EtcdHelper) setKey (key string,ip string,containerID string) error {
	kapi := client.NewKeysAPI(Cli.Client)
	_, err := kapi.Set(context.Background(), key+ip, containerID, nil)
	if err != nil {
		return nil
	}
	return nil
}


func (Cli EtcdHelper) getKey(key string) (NodesInfo *map[string]string){
	kapi := client.NewKeysAPI(Cli.Client)
	//get host
	resp, err := kapi.Get(context.Background(), key, &client.GetOptions{Recursive: true})
	if err != nil {
		log.Fatal(err)
		return
	}
	skydnsNodesInfo := make(map[string]string)
	getAllNode(resp.Node, skydnsNodesInfo)
	return &skydnsNodesInfo
}


func getAllNode(rootNode *client.Node, nodesInfo map[string]string) {
	if !rootNode.Dir {
		nodesInfo[rootNode.Key] = rootNode.Value
		return
	}
	for node := range rootNode.Nodes {
		getAllNode(rootNode.Nodes[node], nodesInfo)
	}
}