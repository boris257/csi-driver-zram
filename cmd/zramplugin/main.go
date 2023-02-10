/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/boris257/csi-driver-zram/pkg/zram"
	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)
}

var (
	endpoint             = flag.String("endpoint", "unix:///tmp/csi.sock", "CSI endpoint")
	nodeID               = flag.String("nodeid", "", "node id")
	driverName           = flag.String("drivername", zram.DefaultDriverName, "name of the driver")
	ver                  = flag.Bool("ver", false, "Print the version and exit.")
	enableGetVolumeStats = flag.Bool("enable-get-volume-stats", true, "allow GET_VOLUME_STATS on agent node")
	enableTopology       = flag.Bool("enable-topology", true, "allow GET_VOLUME_STATS on agent node")
	workingMountDir      = flag.String("working-mount-dir", "/tmp", "working directory for provisioner to mount zram shares temporarily")
)

func main() {
	flag.Parse()
	if *ver {
		info, err := zram.GetVersionYAML(*driverName)
		if err != nil {
			klog.Fatalln(err)
		}
		fmt.Println(info) // nolint
		os.Exit(0)
	}
	if *nodeID == "" {
		// nodeid is not needed in controller component
		klog.Warning("nodeid is empty")
	}
	handle()
	os.Exit(0)
}

func handle() {
	driverOptions := zram.DriverOptions{
		NodeID:               *nodeID,
		DriverName:           *driverName,
		EnableGetVolumeStats: *enableGetVolumeStats,
		WorkingMountDir:      *workingMountDir,
		EnableTopology:       *enableTopology,
	}
	driver := zram.NewDriver(&driverOptions)
	driver.Run(*endpoint, false)
}
