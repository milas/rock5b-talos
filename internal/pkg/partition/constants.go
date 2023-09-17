// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package partition

// Type in partition table.
type Type = string

// GPT partition types.
//
// TODO: should be moved into the blockdevice library.
const (
	EFISystemPartition  Type = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	BIOSBootPartition   Type = "21686148-6449-6E6F-744E-656564454649"
	LinuxFilesystemData Type = "0FC63DAF-8483-4772-8E79-3D69D8477DE4"
)

// FileSystemType is used to format partitions.
type FileSystemType = string

// Filesystem types.
const (
	FilesystemTypeNone FileSystemType = "none"
	FilesystemTypeXFS  FileSystemType = "xfs"
	FilesystemTypeVFAT FileSystemType = "vfat"
	FilesystemTypeExt4 FileSystemType = "ext4"
)

// Partition default sizes.
const (
	MiB = 1024 * 1024

	EFISize      = 100 * MiB
	BIOSGrubSize = 1 * MiB
	BootSize     = 4000 * MiB
	// BootOffset is the expected start of the /boot partition for U-Boot/extlinux on Rockchip devices.
	//
	// HACK: this is for Rock 5 (RK3588) device support.
	BootOffset = 32768 * 512
	// EFIUKISize is the size of the EFI partition when UKI is enabled.
	// With UKI all assets are stored in the EFI partition.
	// This is the size of the old EFISize + BIOSGrubSize + BootSize.
	EFIUKISize = EFISize + BIOSGrubSize + BootSize
	MetaSize   = 1 * MiB
	StateSize  = 100 * MiB
)
