package util

import (
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func GetIPRange(ip_start, ip_end string) []string {
	var ips []string
	ip_s, ipnet_s, start_err := net.ParseCIDR(ip_start)
	if start_err != nil {
		log.Fatal(start_err)
	}
	ip_e, ipnet_e, end_err := net.ParseCIDR(ip_end)
	if end_err != nil {
		log.Fatal(end_err)
	}

	if ipnet_s.Mask.String() != ipnet_e.Mask.String() {
		log.Fatalf("%s and %s are not in the same subnet", ip_start, ip_end)
	}
	for ip := ip_s; ipnet_s.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		if ip.Equal(ip_e) {
			break
		}
	}
	return ips
}

func Get4BytesMask(mask string) string {
	mask_int, _ := strconv.Atoi(mask)
	mask_bytes := []byte(net.CIDRMask(mask_int, 32))
	var mask_strings []string
	for _, mask_byte := range mask_bytes {
		mask_strings = append(mask_strings, strconv.Itoa(int(mask_byte)))
	}
	return strings.Join(mask_strings, ".")
}

func GetIPNetAndMask(ip_cidr string) (string, string) {
	ip_obj, ipnet_obj, err := net.ParseCIDR(ip_cidr)
	if err != nil {
		log.Fatal(err)
	}
	return ip_obj.Mask(ipnet_obj.Mask).String(), strings.Split(ipnet_obj.String(), "/")[1]
}

func GetIPAndCIDR(ip_cidr string) (string, string) {
	ip, cidr, err := net.ParseCIDR(ip_cidr)
	if err != nil {
		log.Fatal(err)
	}
	return ip.String(), cidr.String()
}

func GetMask(ip_cidr string) int {
	_, cidr, _ := net.ParseCIDR(ip_cidr)
	mask, _ := cidr.Mask.Size()
	return mask
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
