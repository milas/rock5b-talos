// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rock5a

import (
	"fmt"
	"github.com/siderolabs/go-procfs/procfs"
	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime/v1alpha1/board/rk3588"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"

	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

const (
	sectorSize        = 512
	ubootOffset int64 = sectorSize * 0x40
)

var ubootImage = fmt.Sprintf("/usr/install/arm64/u-boot/%s/u-boot.img", constants.BoardRock5a)

// Rock5a represents the Radxa Rock 5A board.
//
// Reference: https://docs.radxa.com/en/rock5/rock5a
type Rock5a struct{}

// Name implements the runtime.Board.
func (r *Rock5a) Name() string {
	return constants.BoardRock5a
}

// Install implements the runtime.Board.
func (r *Rock5a) Install(disk string) (err error) {
	var f *os.File

	if f, err = os.OpenFile(disk, os.O_RDWR|unix.O_CLOEXEC, 0o666); err != nil {
		return err
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(ubootImage)
	if err != nil {
		return err
	}
	uboot = uboot[ubootOffset:]

	log.Printf("writing %s (%d) at offset %d", ubootImage, len(uboot), ubootOffset)

	var n int

	n, err = f.WriteAt(uboot, ubootOffset)
	if err != nil {
		return err
	}

	log.Printf("wrote %d bytes", n)

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

// KernelArgs implements the runtime.Board.
func (r *Rock5a) KernelArgs() procfs.Parameters {
	return []*procfs.Parameter{
		procfs.NewParameter("earlycon").Append("uart8250,mmio32,0xfeb50000"),
		procfs.NewParameter("console").Append("ttyFIQ0").Append("tty1"),
		procfs.NewParameter("sysctl.kernel.kexec_load_disabled").Append("1"),
		procfs.NewParameter(constants.KernelParamDashboardDisabled).Append("1"),
		procfs.NewParameter("modprobe.blacklist").Append("pgdrv"),
		procfs.NewParameter("irqchip.gicv3_pseudo_nmi").Append("0"),
		procfs.NewParameter("switolb").Append("1"),
		procfs.NewParameter("coherent_pool").Append("2M"),
		procfs.NewParameter("cgroup_enable").Append("cpuset").Append("memory"),
		procfs.NewParameter("cgroup_memory").Append("1"),
		procfs.NewParameter("swapaccount").Append("1"),
	}
}

// PartitionOptions implements the runtime.Board.
func (r *Rock5a) PartitionOptions() *runtime.PartitionOptions {
	return nil
}

func (r *Rock5a) DeviceTreeBlobPath() string {
	return filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "rk3588s-rock-5a.dtb")
}

func (r *Rock5a) DeviceTreeOverlaysPath() []string {
	return rk3588.DeviceTreeOverlays
}
