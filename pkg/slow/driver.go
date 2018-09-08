package slow

import (
	"github.com/golang/glog"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	Version = "0.3.0"
)

type slowDriver struct {
	driver *csicommon.CSIDriver

	is *identityServer
	ns *nodeServer
	cs *controllerServer

	caps   []*csi.VolumeCapability_AccessMode
	cscaps []*csi.ControllerServiceCapability

	delayOptions DelayOptions
}

func NewSlowDriver(delayOptions DelayOptions) *slowDriver {
	return &slowDriver{delayOptions: delayOptions}
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver, delayOptions *DelayOptions) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		delayOptions:            delayOptions,
	}
}

func NewNodeServer(d *csicommon.CSIDriver, delayOptions *DelayOptions) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		delayOptions:      delayOptions,
	}
}

func (fs *slowDriver) Run(driverName, nodeId, endpoint string) {
	glog.Infof("Slowly running driver: %v version: %v", driverName, Version)
	glog.Infof("Delays: %s", fs.delayOptions.ToString())

	// Initialize default library driver

	fs.driver = csicommon.NewCSIDriver(driverName, Version, nodeId)
	if fs.driver == nil {
		glog.Fatalln("Failed to initialize CSI driver")
	}

	fs.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})

	fs.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})

	// Create gRPC servers

	fs.is = NewIdentityServer(fs.driver)
	fs.ns = NewNodeServer(fs.driver, &fs.delayOptions)
	fs.cs = NewControllerServer(fs.driver, &fs.delayOptions)

	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(endpoint, fs.is, fs.cs, fs.ns)
	server.Wait()
}
