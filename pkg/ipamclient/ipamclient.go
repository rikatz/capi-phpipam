package ipamclient

import (
	"fmt"
	"strings"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
	"github.com/rikatz/capi-phpipam/api/v1alpha1"
)

type IPAMClient struct {
	subnetid int
	ctrl     *addresses.Controller
	subctrl  *subnets.Controller
}

type Subnet struct {
	Mask    string `json:"mask,omitempty"`
	Gateway struct {
		IPAddress string `json:"ip_addr,omitempty"`
	} `json:"gateway,omitempty"`
}

type addrId struct {
	ID        int    `json:"id,omitempty"`
	SubnetID  int    `json:"subnetId,omitempty"`
	IPAddress string `json:"ip,omitempty"`
}

// We create our new client with all the controllers already encapsulated
// TODO: Should support HTTPs and skip insecure :)
func NewIPAMClient(cfg phpipam.Config, subnetid int) *IPAMClient {
	sess := session.NewSession(cfg)
	return &IPAMClient{
		subnetid: subnetid,
		ctrl:     addresses.NewController(sess),
		subctrl:  subnets.NewController(sess),
	}
}

func (i *IPAMClient) GetAddress(hostname string) (string, error) {
	myaddr, err := i.searchForAddress(hostname)
	if err == nil {
		return myaddr.IPAddress, nil
	}

	if err != nil && !isNotFoundHostname(err) {
		return "", fmt.Errorf("failed to find if previous address is in use %w", err)
	}

	addr, err := i.ctrl.CreateFirstFreeAddress(i.subnetid, addresses.Address{Description: hostname, Hostname: hostname})
	if err != nil {
		return "", err
	}
	return addr, nil
}

func (i *IPAMClient) ReleaseAddress(hostname string) error {
	// The library is broken on addrStruct so we need a simple one just to get the allocated ID
	// TODO: Improve error handling, being able to check if the error is something like "not found"
	myaddr, err := i.searchForAddress(hostname)
	if err != nil {
		if isNotFoundHostname(err) {
			return nil
		}
		return fmt.Errorf("failed to find the address, maybe it doesn't exist anymore? %w", err)
	}

	_, err = i.ctrl.DeleteAddress(myaddr.ID, false)
	if err != nil {
		return err
	}
	return nil
}

func (i *IPAMClient) GetSubnetConfig() (*Subnet, error) {
	var subnet Subnet
	err := i.subctrl.SendRequest("GET", fmt.Sprintf("/subnets/%d/", i.subnetid), &struct{}{}, &subnet)
	if err != nil {
		return nil, err
	}
	return &subnet, nil
}

func (i *IPAMClient) searchForAddress(hostname string) (*addrId, error) {
	myaddr := make([]addrId, 0)

	err := i.ctrl.SendRequest("GET", fmt.Sprintf("/addresses/search_hostname/%s", hostname), &struct{}{}, &myaddr)
	if err == nil && len(myaddr) > 0 && myaddr[0].SubnetID == i.subnetid {
		return &myaddr[0], nil
	}
	return nil, err
}

func SpecToClient(spec *v1alpha1.PHPIPAMPoolSpec) (*IPAMClient, error) {
	if spec == nil {
		return nil, fmt.Errorf("spec cannot be null")
	}

	if spec.SubnetID < 0 || spec.Credentials == nil {
		return nil, fmt.Errorf("subnet id and credentials are required")
	}

	return NewIPAMClient(phpipam.Config{
		AppID:    spec.Credentials.AppID,
		Username: spec.Credentials.Username,
		Password: spec.Credentials.Password,
		Endpoint: spec.Credentials.Endpoint,
	}, spec.SubnetID), nil
}

func isNotFoundHostname(err error) bool {
	return strings.Contains(err.Error(), "Hostname not found")
}
