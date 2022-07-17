package configurator

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"io"
	"libvirt.org/go/libvirtxml"
	"os"
)

const emulator = "/usr/bin/qemu-system-x86_64"
const qemuImages = "/var/lib/libvirt/qemu/imgs/"
const ubuntu = qemuImages + "ubuntu.img"

func copyImage(uuid string) (string, error) {
	uf, err := os.Open(ubuntu)
	if err != nil {
		return "", fmt.Errorf("failed to open ubuntu image: %w", err)
	}
	defer uf.Close()

	f, err := os.Create(qemuImages + uuid + ".img")
	if err != nil {
		return "", fmt.Errorf("failed to create image: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, uf); err != nil {
		return "", fmt.Errorf("failed to copy image: %w", err)
	}

	return f.Name(), nil
}

func VMXml(websocketPort int) (string, error) {
	id := uuid.NewString()
	name, err := copyImage(id)
	if err != nil {
		return "", fmt.Errorf("failed to copy image: %w", err)
	}

	d := libvirtxml.Domain{
		Type:        "kvm",
		Name:        id,
		UUID:        id,
		Title:       "Default kvm vm",
		Description: "Default kvm vm",

		VCPU: &libvirtxml.DomainVCPU{
			Value: 1,
		},
		Memory: &libvirtxml.DomainMemory{
			Value: 512,
			Unit:  "MiB",
		},

		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{Arch: "x86_64", Type: "hvm"},
		},

		Devices: &libvirtxml.DomainDeviceList{
			Emulator: emulator,
			Disks: []libvirtxml.DomainDisk{
				{
					XMLName:  xml.Name{},
					Device:   "",
					RawIO:    "",
					SGIO:     "",
					Snapshot: "",
					Model:    "",
					Driver: &libvirtxml.DomainDiskDriver{
						Name: "qemu",
						Type: "qcow2",
					},
					Auth: nil,
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{File: name},
					},
					BackingStore:  nil,
					BackendDomain: nil,
					Geometry:      nil,
					BlockIO:       nil,
					Mirror:        nil,
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "hda",
					},
					IOTune:     nil,
					ReadOnly:   nil,
					Shareable:  nil,
					Transient:  nil,
					Serial:     "",
					WWN:        "",
					Vendor:     "",
					Product:    "",
					Encryption: nil,
					Boot:       nil,
					ACPI:       nil,
					Alias:      nil,
					Address:    nil,
				},
			},

			Graphics: []libvirtxml.DomainGraphic{
				{
					XMLName: xml.Name{},
					VNC: &libvirtxml.DomainGraphicVNC{
						Socket:        "",
						Port:          websocketPort,
						AutoPort:      "yes",
						WebSocket:     websocketPort,
						SharePolicy:   "",
						Passwd:        "",
						PasswdValidTo: "",
						Connected:     "",
						PowerControl:  "",
						Listen:        "0.0.0.0",
					},
				},
			},
		},
	}

	return d.Marshal()
}

func DiskXML() (string, error) {
	d := libvirtxml.DomainDisk{
		XMLName:       xml.Name{},
		Device:        "",
		RawIO:         "",
		SGIO:          "",
		Snapshot:      "",
		Model:         "",
		Driver:        nil,
		Auth:          nil,
		Source:        nil,
		BackingStore:  nil,
		BackendDomain: nil,
		Geometry:      nil,
		BlockIO:       nil,
		Mirror:        nil,
		Target:        nil,
		IOTune:        nil,
		ReadOnly:      nil,
		Shareable:     nil,
		Transient:     nil,
		Serial:        "",
		WWN:           "",
		Vendor:        "",
		Product:       "",
		Encryption:    nil,
		Boot:          nil,
		ACPI:          nil,
		Alias:         nil,
		Address:       nil,
	}

	return d.Marshal()
}
