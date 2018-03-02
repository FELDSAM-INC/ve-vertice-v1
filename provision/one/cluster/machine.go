package cluster

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/virtengine/libgo/cmd"
	"github.com/virtengine/libgo/safe"
	"github.com/virtengine/opennebula-go/api"
	"github.com/virtengine/opennebula-go/compute"
	"github.com/virtengine/opennebula-go/disk"
	"github.com/virtengine/opennebula-go/images"
	"github.com/virtengine/opennebula-go/virtualmachine"
	"net"
	"net/url"
	"strconv"
	"time"
)

// CreateVM creates a vm in the specified node.
// It returns the vm, or an error, in case of failures.
const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

var ErrConnRefused = errors.New("connection refused")

func (c *Cluster) CreateVM(opts compute.VirtualMachine, throttle, storage string) (string, string, string, error) {

	var (
		addr    string
		machine string
		vmid    string
		err     error
	)
	maxTries := 5
	for ; maxTries > 0; maxTries-- {

		nodlist, err := c.Nodes()

		for _, v := range nodlist {
			if v.Metadata[api.ONEZONE] == opts.Region {
				addr = v.Address
				opts.Vnets, opts.ClusterId = c.getVnets(v, opts.Vnets, storage)
				if v.Metadata[api.VCPU_PERCENTAGE] != "" {
					opts.Cpu = cpuThrottle(v.Metadata[api.VCPU_PERCENTAGE], opts.Cpu)
				} else {
					opts.Cpu = cpuThrottle(throttle, opts.Cpu)
				}
			}
		}

		if addr == "" {
			return addr, machine, vmid, fmt.Errorf("%s", cmd.Colorfy("Unavailable region ( "+opts.Region+" ) nodes (hint: start or beat it).\n", "red", "", ""))
		}
		if err == nil {
			machine, vmid, err = c.createVMInNode(opts, addr)
			if err == nil {
				c.handleNodeSuccess(addr)
				break
			}
			log.Errorf("  > Trying... %s", addr)
		}
		shouldIncrementFailures := false
		isCreateMachineErr := false
		baseErr := err
		if nodeErr, ok := baseErr.(OneNodeError); ok {
			isCreateMachineErr = nodeErr.cmd == "createVM"
			baseErr = nodeErr.BaseError()
		}
		if urlErr, ok := baseErr.(*url.Error); ok {
			baseErr = urlErr.Err
		}
		_, isNetErr := baseErr.(*net.OpError)
		if isNetErr || isCreateMachineErr || baseErr == ErrConnRefused {
			shouldIncrementFailures = true
		}
		c.handleNodeError(addr, err, shouldIncrementFailures)
		return addr, machine, vmid, err
	}
	if err != nil {
		return addr, machine, vmid, fmt.Errorf("CreateVM: maximum number of tries exceeded, last error: %s", err.Error())
	}
	return addr, machine, vmid, err
}

//create a vm in a node.
func (c *Cluster) createVMInNode(opts compute.VirtualMachine, nodeAddress string) (string, string, error) {
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return "", "", err
	}

	if opts.ClusterId != "" {
		opts.TemplateName = node.template
	} else {
		opts.TemplateName = opts.Image
	}

	opts.T = node.Client

	res, err := opts.Create()
	if err != nil {
		return "", "", err
	}
	vmid := res.(string)
	return opts.Name, vmid, nil
}

func (c *Cluster) GetVM(opts virtualmachine.Vnc, region string) (*virtualmachine.VM, error) {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return nil, err
	}

	opts.T = node.Client
	res, err := opts.GetVm()
	if err != nil {
		return nil, wrapErrorWithCmd(node, err, "GetVM")
	}

	return res, err
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) DestroyVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Delete()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DestroyVM")
	}

	return nil
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) ForceDestoryVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.RecoverDelete()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DestroyVM")
	}

	return nil
}

func (c *Cluster) VM(opts compute.VirtualMachine, action string) error {
	switch action {
	case START:
		return c.StartVM(opts)
	case STOP:
		return c.StopVM(opts)
	case RESTART:
		return c.RestartVM(opts)
	default:
		return nil
	}
}
func (c *Cluster) StartVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Resume()
	if err != nil {
		return wrapErrorWithCmd(node, err, "StartVM")
	}

	return nil
}

func (c *Cluster) RestartVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Reboot()
	if err != nil {
		return wrapErrorWithCmd(node, err, "RebootVM")
	}

	return nil
}

func (c *Cluster) StopVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Poweroff()
	if err != nil {
		return wrapErrorWithCmd(node, err, "StopVM")
	}

	return nil
}

func (c *Cluster) getNodeRegion(region string) (node, error) {
	return c.getNode(func(s Storage) (Node, error) {
		return s.RetrieveNode(region)
	})
}

func (c *Cluster) SnapVMDisk(opts compute.Image) (string, error) {
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return "", err
	}
	opts.T = node.Client

	res, err := opts.DiskSnap()
	if err != nil {
		return "", wrapErrorWithCmd(node, err, "CreateSnap")
	}
	imageId := res.(string)
	return imageId, nil
}

func (c *Cluster) RemoveSnap(opts compute.Image) error {
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.RemoveImage()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DeleteSnap")
	}

	return nil
}

func (c *Cluster) IsSnapReady(v *images.Image, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	v.T = node.Client
	err = safe.WaitCondition(10*time.Minute, 10*time.Second, func() (bool, error) {
		res, err := v.ImageShow()
		if err != nil || res.State_string() == "failure" {
			return false, fmt.Errorf("fails to create snapshot")
		}
		return (res.State_string() == "ready"), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) GetDiskId(vd *disk.VmDisk, region string) ([]int, error) {
	var a []int
	node, err := c.getNodeRegion(region)
	if err != nil {
		return a, err
	}
	vd.T = node.Client

	dsk, err := vd.ListDisk()
	if err != nil {
		return a, err
	}

	a = dsk.GetDiskIds()
	return a, nil
}

func (c *Cluster) AttachDisk(v *disk.VmDisk, region string) error {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	v.T = node.Client

	_, err = v.AttachDisk()
	if err != nil {
		return wrapErrorWithCmd(node, err, "AttachDisk")
	}

	return nil
}

func (c *Cluster) DetachDisk(v *disk.VmDisk, region string) error {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	v.T = node.Client

	_, err = v.DetachDisk()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DetachDisk")
	}

	return nil
}

func cpuThrottle(vcpu, cpu string) string {
	ThrottleFactor, _ := strconv.Atoi(vcpu)
	cpuThrottleFactor := float64(ThrottleFactor)
	ICpu, _ := strconv.Atoi(cpu)
	throttle := float64(ICpu)
	realCPU := throttle / cpuThrottleFactor
	cpu = strconv.FormatFloat(realCPU, 'f', 6, 64) //ugly, compute has the info.
	return cpu
}
