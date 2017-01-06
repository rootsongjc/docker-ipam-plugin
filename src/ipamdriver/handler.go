package ipamdriver

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/ipam"
	netlabel "github.com/docker/libnetwork/netlabel"
	"util"
)

type MyIPAMHandler struct {
}

func (iph *MyIPAMHandler) GetCapabilities() (response *ipam.CapabilitiesResponse, err error) {
	log.Infof("GetCapabilities")

	return &ipam.CapabilitiesResponse{RequiresMACAddress: true}, nil
}

func (iph *MyIPAMHandler) GetDefaultAddressSpaces() (response *ipam.AddressSpacesResponse, err error) {
	log.Infof("GetDefaultAddressSpaces")
	return &ipam.AddressSpacesResponse{}, nil
}

func (iph *MyIPAMHandler) RequestPool(request *ipam.RequestPoolRequest) (response *ipam.RequestPoolResponse, err error) {
	var request_json []byte = nil
	request_json, err = json.Marshal(request)
	if err != nil {
		return nil, err
	}
	log.Infof("RequestPool: %s", request_json)
	ip_net, _ := util.GetIPNetAndMask(request.Pool)
	_, ip_cidr := util.GetIPAndCIDR(request.Pool)
	options := request.Options
	return &ipam.RequestPoolResponse{ip_net, ip_cidr, options}, nil
}

func (iph *MyIPAMHandler) ReleasePool(request *ipam.ReleasePoolRequest) (err error) {
	var request_json []byte = nil
	request_json, err = json.Marshal(request)
	if err != nil {
		return err
	}
	log.Infof("ReleasePool %s is danger, you should do this by manual.", request_json)
	return nil
}

func (iph *MyIPAMHandler) RequestAddress(request *ipam.RequestAddressRequest) (response *ipam.RequestAddressResponse, err error) {
	var request_json []byte = nil
	request_json, err = json.Marshal(request)
	if err != nil {
		return nil, err
	}
	log.Infof("RequestAddress %s", request_json)
	ip_net := request.PoolID
	ip := request.Address
	config, _ := GetConfig(ip_net)

	if value, ok := request.Options["RequestAddressType"]; ok && value == netlabel.Gateway || len(request.Options) == 0 {
		log.Infof("Skip allocate gateway ip %s", ip)
		return &ipam.RequestAddressResponse{fmt.Sprintf("%s/%s", ip, config.Mask), nil}, nil
	}
	ip, err = AllocateIP(ip_net, ip)
	return &ipam.RequestAddressResponse{fmt.Sprintf("%s/%s", ip, config.Mask), nil}, err
}

func (iph *MyIPAMHandler) ReleaseAddress(request *ipam.ReleaseAddressRequest) (err error) {
	var request_json []byte = nil
	request_json, err = json.Marshal(request)
	if err != nil {
		return err
	}
	log.Infof("ReleaseAddress %s", request_json)
	err = ReleaseIP(request.PoolID, request.Address)
	return err
}
