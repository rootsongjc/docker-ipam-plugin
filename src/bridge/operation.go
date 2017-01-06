package bridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	log "github.com/Sirupsen/logrus"
	"util"
)

func createBridge(ip, subnet, gateway string) error {
	var err error = nil

	if err := create_network(ip, subnet, gateway); err != nil {
		return err
	}
	if err = configure_network(ip, subnet, gateway); err != nil {
		return err
	}
	return err
}

func create_network(ip, subnet, gateway string) error {
	command := "docker"
	args := fmt.Sprint("network create ",
		"--opt=com.docker.network.bridge.enable_icc=true ",
		"--opt=com.docker.network.bridge.enable_ip_masquerade=false ",
		"--opt=com.docker.network.bridge.host_binding_ipv4=0.0.0.0 ",
		"--opt=com.docker.network.bridge.name=br0 ",
		"--opt=com.docker.network.driver.mtu=1500 ",
		"--ipam-driver=talkingdata ",
		"--subnet=%s ",
		"--gateway=%s ",
		"--aux-address=DefaultGatewayIPv4=%s ",
		"mynet")

	args = fmt.Sprintf(args, subnet, ip, gateway)

	var out []byte
	out, err := exec.Command(command, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		log.Fatal(err, string(out))
		return err
	}

	return nil
}

func configure_network(ip, subnet, gateway string) error {
	bridge_command := "BRIDGE=br0"
	_, interface_name, _ := get_network_information()
	old_interface := "/etc/sysconfig/network-scripts/ifcfg-" + interface_name
	new_interface := "/etc/sysconfig/network-scripts/ifcfg-br0"
	netmask := util.Get4BytesMask(strings.Split(subnet, "/")[1])

	old_fd, old_err := os.OpenFile(old_interface, os.O_RDWR|os.O_APPEND, 0666)
	if old_err != nil {
		log.Fatal(old_err)
		return old_err
	}
	defer old_fd.Close()
	br, _ := ioutil.ReadAll(old_fd)

	if strings.Contains(string(br), bridge_command) != true {
		buf := []byte(bridge_command)
		old_fd.Write(buf)
	}

	new_fd, new_err := os.OpenFile(new_interface, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if new_err != nil {
		log.Fatal(new_err)
		return new_err
	}
	defer new_fd.Close()

	br0_content := fmt.Sprint("DEVICE=br0\n",
		"TYPE=Bridge\n",
		"BOOTPROTO=static\n",
		"IPADDR=%s\n",
		"GATEWAY=%s\n",
		"NETMASK=%s\n",
		"ONBOOT=yes\n",
		"NOZEROCONF=yes\n",
		"IPV6INIT=no\n",
		"NM_CONTROLLED=no\n",
		"DELAY=0")
	br0_content = fmt.Sprintf(br0_content, ip, gateway, netmask)
	new_fd.WriteString(br0_content)

	return nil
}

func restart_network() error {
	command := "systemctl"

	disable_args := fmt.Sprint("disable ", "NetworkManager")
	_, err := exec.Command(command, strings.Split(disable_args, " ")...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	stop_args := fmt.Sprint("stop ", "NetworkManager")
	_, err = exec.Command(command, strings.Split(stop_args, " ")...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Restarting network")
	restart_args := fmt.Sprint("restart ", "network")
	_, err = exec.Command(command, strings.Split(restart_args, " ")...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Info("Restart network done")

	return nil
}

func get_network_information() (gateway, interface_name, local_ip string) {
	command := "ip"
	args := fmt.Sprint("route get ", "8.8.8.8")
	out, err := exec.Command(command, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		return
	}
	network_information := strings.Split(string(out), " ")
	gateway = network_information[2]
	interface_name = network_information[4]
	local_ip = network_information[6]

	return
}
