package slow

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
	delayOptions *DelayOptions
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if err := cs.validateCreateVolumeRequest(req); err != nil {
		glog.Errorf("CreateVolumeRequest validation failed: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	runDelay("CreateVolume", cs.delayOptions.CreateVolume)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            string(newVolumeID()),
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			Attributes:    req.GetParameters(),
		},
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if err := cs.validateDeleteVolumeRequest(req); err != nil {
		glog.Errorf("DeleteVolumeRequest validation failed: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	runDelay("DeleteVolume", cs.delayOptions.DeleteVolume)

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, cap := range req.VolumeCapabilities {
		if cap.GetBlock() != nil {
			return &csi.ValidateVolumeCapabilitiesResponse{Supported: false, Message: ""}, nil
		}
	}
	return &csi.ValidateVolumeCapabilitiesResponse{Supported: true}, nil
}
