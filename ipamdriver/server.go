package ipamdriver

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/ipam"

	"oam-docker-ipam/db"
	"oam-docker-ipam/util"
)

const (
	network_key_prefix = "/talkingdata/containers"
)

type Config struct {
	Ipnet string
	Mask  string
}

func StartServer() {
	d := &MyIPAMHandler{}
	h := ipam.NewHandler(d)
	h.ServeUnix("root", "talkingdata")
}

func AllocateIPRange(ip_start, ip_end string) []string {
	ips := util.GetIPRange(ip_start, ip_end)
	ip_net, mask := util.GetIPNetAndMask(ip_start)
	for _, ip := range ips {
		if checkIPAssigned(ip_net, ip) {
			log.Warnf("IP %s has been allocated", ip)
			continue
		}
		db.SetKey(filepath.Join(network_key_prefix, ip_net, "pool", ip), "")
	}
	initializeConfig(ip_net, mask)
	fmt.Println("Allocate Containers IP Done! Total:", len(ips))
	return ips
}

func ReleaseIP(ip_net, ip string) error {
	err := db.DeleteKey(filepath.Join(network_key_prefix, ip_net, "assigned", ip))
	if err != nil {
		log.Infof("Skip Release IP %s", ip)
		return nil
	}
	err = db.SetKey(filepath.Join(network_key_prefix, ip_net, "pool", ip), "")
	if err == nil {
		log.Infof("Release IP %s", ip)
	}
	return nil
}

func AllocateIP(ip_net, ip string) (string, error) {
	ip_pool, err := db.GetKeys(filepath.Join(network_key_prefix, ip_net, "pool"))
	if err != nil {
		return ip, err
	}
	if len(ip_pool) == 0 {
		return ip, errors.New("Pool is empty")
	}
	if ip == "" {
		find_ip := strings.Split(ip_pool[0].Key, "/")
		ip = find_ip[len(find_ip)-1]
	}
	exist := checkIPAssigned(ip_net, ip)
	if exist == true {
		return ip, errors.New(fmt.Sprintf("IP %s has been allocated", ip))
	}
	err = db.DeleteKey(filepath.Join(network_key_prefix, ip_net, "pool", ip))
	if err != nil {
		return ip, err
	}
	db.SetKey(filepath.Join(network_key_prefix, ip_net, "assigned", ip), "")
	log.Infof("Allocated IP %s", ip)
	return ip, err
}

func checkIPAssigned(ip_net, ip string) bool {
	if exist := db.IsKeyExist(filepath.Join(network_key_prefix, ip_net, "assigned", ip)); exist {
		return true
	}
	return false
}

func initializeConfig(ip_net, mask string) error {
	config := &Config{Ipnet: ip_net, Mask: mask}
	config_bytes, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetKey(filepath.Join(network_key_prefix, ip_net, "config"), string(config_bytes))
	if err == nil {
		log.Infof("Initialized Config %s for network %s", string(config_bytes), ip_net)
	}
	return err
}

func DeleteNetWork(ip_net string) error {
	err := db.DeleteKey(filepath.Join(network_key_prefix, ip_net))
	if err == nil {
		log.Infof("DeleteNetwork %s", ip_net)
	}
	return err
}

func GetConfig(ip_net string) (*Config, error) {
	config, err := db.GetKey(filepath.Join(network_key_prefix, ip_net, "config"))
	if err == nil {
		log.Debugf("GetConfig %s from network %s", config, ip_net)
	}
	conf := &Config{}
	json.Unmarshal([]byte(config), conf)
	return conf, err
}
