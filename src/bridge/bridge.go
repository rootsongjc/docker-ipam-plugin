package bridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"db"
	"util"
	"path/filepath"
	"strings"
	log "github.com/Sirupsen/logrus"
)

const (
	network_key_prefix = "/talkingdata/hosts"
)

type Config struct {
	Subnet  string
	Gateway string
}

func AllocateHostRange(ip_start, ip_end, gateway string) []string {
	ips := util.GetIPRange(ip_start, ip_end)
	ip_net, mask := util.GetIPNetAndMask(ip_start)
	for _, ip := range ips {
		if checkIPAssigned(ip) {
			log.Warnf("IP %s has been allocated", ip)
			continue
		}
		db.SetKey(filepath.Join(network_key_prefix, "pool", ip), "")
	}
	initializeConfig(ip_net, fmt.Sprint(ip_net, "/", mask), gateway)
	fmt.Println("Allocate Hosts Done! Total:", len(ips))
	return ips
}

func initializeConfig(ip_net, subnet, gateway string) error {
	config := &Config{Subnet: subnet, Gateway: gateway}
	config_bytes, _ := json.Marshal(config)
	err := db.SetKey(filepath.Join(network_key_prefix, "config"), string(config_bytes))
	if err == nil {
		log.Infof("Initialized Config %s for network %s", string(config_bytes), ip_net)
	}
	return err
}

func getConfig() (*Config, error) {
	config, err := db.GetKey(filepath.Join(network_key_prefix, "config"))
	if err == nil {
		log.Debugf("getConfig %s", config)
	}
	conf := &Config{}
	json.Unmarshal([]byte(config), conf)
	return conf, err
}

func allocateHost(ip string) error {
	if ip == "" {
		return errors.New("arg ip is lack")
	}
	err := db.DeleteKey(filepath.Join(network_key_prefix, "pool", ip))
	if err != nil {
		return err
	}
	if err = db.SetKey(filepath.Join(network_key_prefix, "assigned", ip), ""); err != nil {
		return err
	}
	log.Infof("Allocated host %s", ip)
	return nil
}

func getHost(ip string) (string, error) {
	ip_pool, err := db.GetKeys(filepath.Join(network_key_prefix, "pool"))
	if err != nil {
		return "", err
	}
	if len(ip_pool) == 0 {
		return "", errors.New("Pool is empty")
	}
	if ip == "" {
		find_ip := strings.Split(ip_pool[0].Key, "/")
		ip = find_ip[len(find_ip)-1]
	} else if exist := db.IsKeyExist(filepath.Join(network_key_prefix, "pool", ip)); exist != true {
		return "", errors.New(fmt.Sprintf("Host %s not in pool", ip))
	}
	if assigned := checkIPAssigned(ip); assigned == true {
		return "", errors.New(fmt.Sprintf("Host %s has been allocated", ip))
	}
	return ip, nil
}

func checkIPAssigned(ip string) bool {
	if exist := db.IsKeyExist(filepath.Join(network_key_prefix, "assigned", ip)); exist {
		return true
	}
	return false
}

/*
func checkNetwork() error {
	num := 0
	interface_name := ""
	for {
		_, interface_name, _ = get_network_information()
		if num >= 10 {
			return errors.New("network is not ok, after 10 times try.")
		} else if interface_name == "br0" {
			return nil
		}
		num += 1
		time.Sleep(2 * time.Second)
		log.Infof("try test network %s time", num)
	}
	return nil
}
*/

func ReleaseHost(ip string) error {
	err := db.DeleteKey(filepath.Join(network_key_prefix, "assigned", ip))
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetKey(filepath.Join(network_key_prefix, "pool", ip), "")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Release host %s", ip)
	return err
}

func CreateNetwork(ip string) {
	var assigned_ip string
	var config *Config
	var err error

	if config, err = getConfig(); err != nil {
		log.Fatal(err)
	}
	if assigned_ip, err = getHost(ip); err != nil {
		log.Fatal(err)
	}
	if err = allocateHost(assigned_ip); err != nil {
		log.Fatal(err)
	}
	if err = createBridge(assigned_ip, config.Subnet, config.Gateway); err != nil {
		log.Fatal(err)
	}
	if err = restart_network(); err != nil {
		log.Fatal(err)
	}
	log.Infof("Create network %s done", assigned_ip)
}
