package zram

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"k8s.io/klog/v2"
	"k8s.io/utils/exec"
	"k8s.io/utils/mount"
)

type ZRAMDevice struct {
	id      int
	devPath string
	sysPath string
	mounter *mount.SafeFormatAndMount
}

func NewZRAMDevice() (*ZRAMDevice, error) {
	data, err := ioutil.ReadFile("/sys/class/zram-control/hot_add")
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid device number: %s", err.Error())
	}
	return &ZRAMDevice{
		id:      id,
		devPath: fmt.Sprintf("/dev/zram%d", id),
		sysPath: fmt.Sprintf("/sys/block/zram%d", id),
		mounter: &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: exec.New()},
	}, nil
}

func NewZRAMDeviceFromId(id int) (*ZRAMDevice, error) {
	return &ZRAMDevice{
		id:      id,
		devPath: fmt.Sprintf("/dev/zram%d", id),
		sysPath: fmt.Sprintf("/sys/block/zram%d", id),
		mounter: &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: exec.New()},
	}, nil
}

func NewZRAMDeviceFromMountPath(mountPath string) (*ZRAMDevice, error) {
	mounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: exec.New()}
	dev, _, err := GetDeviceNameFromMountPath(mounter, mountPath)
	if err != nil {
		return nil, fmt.Errorf("faied to get device from mount path: %s", err.Error())
	}
	devName := filepath.Base(dev)
	if !strings.HasPrefix(devName, "zram") {
		return nil, fmt.Errorf("invalid device: %s", dev)
	}
	id, err := strconv.Atoi(strings.TrimPrefix(devName, "zram"))
	if err != nil {
		return nil, fmt.Errorf("invalid device: %s", err.Error())
	}
	return &ZRAMDevice{
		id:      id,
		devPath: fmt.Sprintf("/dev/zram%d", id),
		sysPath: fmt.Sprintf("/sys/block/zram%d", id),
		mounter: &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: exec.New()},
	}, nil
}

func (d *ZRAMDevice) GetId() int {
	return d.id
}

func (d *ZRAMDevice) GetDevPath() string {
	return d.devPath
}

func (d *ZRAMDevice) GetSysPath() string {
	return d.sysPath
}

func (d *ZRAMDevice) GetMounter() *mount.SafeFormatAndMount {
	return d.mounter
}

func (d *ZRAMDevice) Remove() error {
	fd, err := os.OpenFile("/sys/class/zram-control/hot_remove", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = fd.WriteString(strconv.Itoa(d.id))
	if err != nil {
		return err
	}
	return fd.Close()
}

func (d *ZRAMDevice) writeSysFile(name, data string) error {
	fileName := filepath.Join(d.sysPath, name)
	klog.Infof("write %s: %s", fileName, data)
	fd, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	_, err = fd.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return fd.Close()
}

func (d *ZRAMDevice) Reset() error {
	return d.writeSysFile("reset", "1")
}

func (d *ZRAMDevice) SetMaxCompStreams(maxCompStreams int) error {
	return d.writeSysFile("max_comp_streams", strconv.Itoa(maxCompStreams))
}

func (d *ZRAMDevice) SetCompAlgorithm(compAlgorithm string) error {
	return d.writeSysFile("comp_algorithm", compAlgorithm) // lzo-rle is default
}

func (d *ZRAMDevice) SetDiskSize(diskSize int64) error {
	return d.writeSysFile("disksize", strconv.FormatInt(diskSize, 10))
}

func (d *ZRAMDevice) SetMemLimit(memLimit int64) error {
	return d.writeSysFile("mem_limit", strconv.FormatInt(memLimit, 10))
}

func (d *ZRAMDevice) RefCount() (int, error) {
	mps, err := d.mounter.List()
	if err != nil {
		return 0, err
	}

	// Find all references to the device.
	refCount := 0
	for i := range mps {
		if mps[i].Device == d.devPath {
			klog.Infof("dev: %s mount-path: %s", mps[i].Device, mps[i].Path)
			refCount++
		}
	}
	return refCount, nil
}

func (d *ZRAMDevice) UnmountAndCleanup() error {
	mps, err := d.mounter.List()
	if err != nil {
		return fmt.Errorf("failed to list devices: %v", err)
	}

	// Find all references to the device.
	for i := range mps {
		if mps[i].Device == d.devPath {
			mountErr := mount.CleanupMountPoint(mps[i].Path, d.mounter, false)
			if err != nil {
				err = multierror.Append(err, mountErr)
				klog.Errorf("dev: %s unmount: %s error: %v", mps[i].Device, mps[i].Path, err)
			} else {
				klog.Infof("dev: %s unmount: %s success", mps[i].Device, mps[i].Path)
			}
		}
	}
	return err
}

func (d *ZRAMDevice) FormatAndMount(mountPath, fsType string, options []string) error {
	notMnt, err := d.mounter.IsLikelyNotMountPoint(mountPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("heuristic determination of mount point failed: %v", err)
	}
	if !notMnt {
		klog.Infof("zram: %s already mounted", mountPath)
		return nil
	}

	if err := os.MkdirAll(mountPath, 0750); err != nil {
		klog.Errorf("zram: failed to mkdir %s, error", mountPath)
		return err
	}

	err = d.mounter.FormatAndMount(d.devPath, mountPath, fsType, options)
	if err != nil {
		klog.Errorf("zram: failed to mount zram volume %s [%s] to %s, error %v", d.devPath, fsType, mountPath, err)
	}
	return err
}

// GetDeviceNameFromMount given a mnt point, find the device from /proc/mounts
// returns the device name, reference count, and error code.
func GetDeviceNameFromMountPath(mounter mount.Interface, mountPath string) (string, int, error) {
	mps, err := mounter.List()
	if err != nil {
		return "", 0, err
	}

	// Find the device name.
	// FIXME if multiple devices mounted on the same mount path, only the first one is returned.
	device := ""
	// If mountPath is symlink, need get its target path.
	slTarget, err := filepath.EvalSymlinks(mountPath)
	if err != nil {
		slTarget = mountPath
	}
	for i := range mps {
		if mps[i].Path == slTarget {
			device = mps[i].Device
			break
		}
	}

	// Find all references to the device.
	refCount := 0
	for i := range mps {
		if mps[i].Device == device {
			klog.Infof("dev: %s mount-path: %s", mps[i].Device, mps[i].Path)
			refCount++
		}
	}
	return device, refCount, nil
}
