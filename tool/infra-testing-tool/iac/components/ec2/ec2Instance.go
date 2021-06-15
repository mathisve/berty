package ec2

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"infratesting/iac"
	"infratesting/iac/components/networking"
)

type Instance struct {
	Name         string
	InstanceType string
	KeyName      string

	AvailabilityZone string

	RootBlockDevice RootBlockDevice

	NodeType string
	UserData string

	// NETWORKING
	NetworkInterfaces          []*networking.NetworkInterface
	NetworkInterfaceAttachment []NetworkInterfaceAttachment

	Tags []Tag
}

type NetworkInterfaceAttachment struct {
	DeviceIndex        int
	NetworkInterfaceId string
}

type Tag struct {
	Key string
	Value string
}

func NewInstance() Instance {
	// with defaults in place
	return Instance{
		Name:         fmt.Sprintf("%s-%s", Ec2NamePrefix, uuid.NewString()),
		InstanceType: Ec2InstanceTypeDefault,
		KeyName:      Ec2InstanceKeyNameDefault,

		RootBlockDevice: NewRootBlockDevice(),
	}
}

func NewInstanceWithAttributes(ni *networking.NetworkInterface) (c Instance) {
	c = NewInstance()
	c.NetworkInterfaces = []*networking.NetworkInterface{
		ni,
	}

	return c
}

// GetTemplate returns Ec2 template
func (c Instance) GetTemplate() string {
	return Ec2HCLTemplate
}

// GetType returns Ec2 type
func (c Instance) GetType() string {
	return Ec2Type
}

// GetPrivateIp returns the terraform formatting of this instances' ip
func (c Instance) GetPrivateIp() string {
	return fmt.Sprintf("aws_instance.%s.private_ip", c.Name)
}

// SetNodeType sets the node type
func (c *Instance) SetNodeType(s string) {
	c.NodeType = s
}

// Validate validates the component
func (c Instance) Validate() (iac.Component, error) {
	// checks if NetworkInterface is attached/configured
	if len(c.NetworkInterfaces) > 0 {
		// generates NetworkInterfaceAttachment form c.NetworkInterfaces
		for i, ni := range c.NetworkInterfaces {
			var nia = NetworkInterfaceAttachment{
				DeviceIndex:        i,
				NetworkInterfaceId: ni.GetId(),
			}
			c.NetworkInterfaceAttachment = append(c.NetworkInterfaceAttachment, nia)
		}
	} else if len(c.NetworkInterfaces) > 2 && c.InstanceType == Ec2InstanceTypeDefault {
		// more than 2 network interfaces on a t2.micro is not allowed
		return c, errors.New(Ec2ErrTooManyNetworkInterfaces)
	} else {
		// can't have no network interfaces, there is just no point
		return c, errors.New(Ec2ErrNoNetworkInterfaceConfigured)
	}

	// checks if all network interfaces are on the same Availability Zone
	for _, networkInterface := range c.NetworkInterfaces {
		if c.AvailabilityZone == "" {
			c.AvailabilityZone = networkInterface.GetAvailabilityZone()
		} else if c.AvailabilityZone != networkInterface.GetAvailabilityZone() {
			return c, errors.New(Ec2ErrNetworkInterfaceAZMismatch)
		}

	}

	// checks RootBlockDevice size
	// we don't check the VolumeType because it doesn't really matter for now
	if c.RootBlockDevice.VolumeSize < 8 {
		return c, errors.New(Ec2ErrRootBlockDeviceTooSmall)
	}

	c.Tags = append(c.Tags, Tag{
		Key:   Ec2TagName,
		Value: c.Name,
	})

	c.Tags = append(c.Tags, Tag{
		Key:   Ec2TagType,
		Value: c.NodeType,
	})

	return c, nil
}
