/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package zram

import (
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"

	csicommon "github.com/boris257/csi-driver-zram/pkg/csi-common"
	"github.com/boris257/csi-driver-zram/pkg/mounter"
	"k8s.io/klog/v2"
	mount "k8s.io/mount-utils"
)

const (
	DefaultDriverName = "zram.csi.k8s.io"
	TopologyKeyNode   = "topology.hostpath.csi/node"
	mountOptionsField = "mountoptions"
	capacityField     = "capacity"
)

// DriverOptions defines driver parameters specified in driver deployment
type DriverOptions struct {
	NodeID               string
	DriverName           string
	EnableGetVolumeStats bool
	EnableTopology       bool
	WorkingMountDir      string
}

// Driver implements all interfaces of CSI drivers
type Driver struct {
	csicommon.CSIDriver
	mounter *mount.SafeFormatAndMount
	// A map storing all volumes with ongoing operations so that additional operations
	// for that same volume (as defined by VolumeID) return an Aborted error
	volumeLocks          *volumeLocks
	workingMountDir      string
	enableGetVolumeStats bool
	enableTopology       bool
}

// NewDriver Creates a NewCSIDriver object. Assumes vendor version is equal to driver version &
// does not support optional driver plugin info manifest field. Refer to CSI spec for more details.
func NewDriver(options *DriverOptions) *Driver {
	driver := Driver{}
	driver.Name = options.DriverName
	driver.Version = driverVersion
	driver.NodeID = options.NodeID
	driver.enableGetVolumeStats = options.EnableGetVolumeStats
	driver.enableTopology = options.EnableTopology
	driver.workingMountDir = options.WorkingMountDir
	driver.volumeLocks = newVolumeLocks()
	return &driver
}

// Run driver initialization
func (d *Driver) Run(endpoint string, testMode bool) {
	versionMeta, err := GetVersionYAML(d.Name)
	if err != nil {
		klog.Fatalf("%v", err)
	}
	klog.V(2).Infof("\nDRIVER INFORMATION:\n-------------------\n%s\n\nStreaming logs below:", versionMeta)

	d.mounter, err = mounter.NewSafeMounter()
	if err != nil {
		klog.Fatalf("Failed to get safe mounter. Error: %v", err)
	}

	// Initialize default library driver
	d.AddControllerServiceCapabilities(
		[]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
		})

	d.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	})

	nodeCap := []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
		csi.NodeServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
	}
	if d.enableGetVolumeStats {
		nodeCap = append(nodeCap, csi.NodeServiceCapability_RPC_GET_VOLUME_STATS)
	}
	d.AddNodeServiceCapabilities(nodeCap)

	s := csicommon.NewNonBlockingGRPCServer()
	// Driver d act as IdentityServer, ControllerServer and NodeServer
	s.Start(endpoint, d, d, d, testMode)
	s.Wait()
}

func IsCorruptedDir(dir string) bool {
	_, pathErr := mount.PathExists(dir)
	return pathErr != nil && mount.IsCorruptedMnt(pathErr)
}

// setKeyValueInMap set key/value pair in map
// key in the map is case insensitive, if key already exists, overwrite existing value
func setKeyValueInMap(m map[string]string, key, value string) {
	if m == nil {
		return
	}
	for k := range m {
		if strings.EqualFold(k, key) {
			m[k] = value
			return
		}
	}
	m[key] = value
}

// replaceWithMap replace key with value for str
func replaceWithMap(str string, m map[string]string) string {
	for k, v := range m {
		if k != "" {
			str = strings.ReplaceAll(str, k, v)
		}
	}
	return str
}
