package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gman0/slow-csi/pkg/slow"
	"github.com/golang/glog"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName = flag.String("drivername", "csi-slowplugin", "name of the driver")
	nodeId     = flag.String("nodeid", "", "node id")

	defaultDelay = flag.String("defaultdelay", "", "default delay applied to all CSI calls")
	delay        = flag.String("delay", "", "delay settings for individual CSI calls, e.g. CreateVolume=20..50,DeleteVolume=5,NodeUnstageVolume=inf")
)

func main() {
	flag.Parse()

	var delayOptions slow.DelayOptions

	if *defaultDelay != "" {
		if dr, err := parseDelay(*defaultDelay); err != nil {
			glog.Fatalf("invalid format in default delay: %v", err)
		} else {
			setDefaultDelay(&delayOptions, dr)
		}
	}

	if *delay != "" {
		if err := parseDelayArgs(&delayOptions, *delay); err != nil {
			glog.Fatalf("invalid delay: %v", err)
		}
	}

	driver := slow.NewSlowDriver(delayOptions)
	driver.Run(*driverName, *nodeId, *endpoint)

	os.Exit(0)
}

func setDefaultDelay(o *slow.DelayOptions, dr slow.DelayRange) {
	o.CreateVolume = dr
	o.DeleteVolume = dr
	o.NodeStageVolume = dr
	o.NodeUnstageVolume = dr
	o.NodePublishVolume = dr
	o.NodeUnpublishVolume = dr
}

func parseDelayArgs(o *slow.DelayOptions, args string) error {
	parts := strings.Split(args, ",")
	for _, p := range parts {
		eq := strings.Index(p, "=")
		if eq <= 0 {
			return fmt.Errorf("invalid format: '%s' is not in format of 'RPC=delay'", p)
		}

		rpc := p[:eq]
		val := p[eq+1:]

		var dr *slow.DelayRange

		switch rpc {
		case "CreateVolume":
			dr = &o.CreateVolume
		case "DeleteVolume":
			dr = &o.DeleteVolume
		case "NodeStageVolume":
			dr = &o.NodeStageVolume
		case "NodeUnstageVolume":
			dr = &o.NodeUnstageVolume
		case "NodePublishVolume":
			dr = &o.NodePublishVolume
		case "NodeUnpublishVolume":
			dr = &o.NodeUnpublishVolume
		default:
			return fmt.Errorf("no delay option for RPC '%s'", rpc)
		}

		ret, err := parseDelay(val)
		if err != nil {
			return fmt.Errorf("failed to parse delay value for '%s': %v", p, err)
		}

		*dr = ret
	}

	return nil
}

func parseDelay(delay string) (slow.DelayRange, error) {
	if delay == "inf" {
		return slow.DelayRange{-1, -1}, nil
	}

	dots := strings.Index(delay, "..")
	if dots == -1 {
		if d, err := strconv.Atoi(delay); err != nil {
			return slow.DelayRange{}, err
		} else {
			return slow.DelayRange{d, d}, nil
		}
	} else {
		min, err := strconv.Atoi(delay[:dots])
		if err != nil {
			return slow.DelayRange{}, err
		}

		max, err := strconv.Atoi(delay[dots+2:])
		if err != nil {
			return slow.DelayRange{}, err
		}

		if min >= max {
			return slow.DelayRange{}, fmt.Errorf("%d ≮ %d", min, max)
		}

		return slow.DelayRange{Min: min, Max: max}, nil
	}
}
