package server

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"govps/pkg/configurator"
	"libvirt.org/go/libvirtxml"
	"log"
	"net"
	"time"
)

func NewUsecase(libvirtSocket string) *Usecase {
	c, err := net.DialTimeout("unix", libvirtSocket, 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}

	l := libvirt.NewWithDialer(dialers.NewAlreadyConnected(c))
	if err = l.Connect(); err != nil {
		log.Fatalf("failed to connect to libvirt: %v", err)
	}

	return &Usecase{l}
}

type Usecase struct {
	l *libvirt.Libvirt
}

type Hypervisor struct {
	LibcVersion uint64
	Mem         uint64
	MaxVcpus    int32
	Cpus        int32
	Mhz         int32
	Nodes       int32
	Sockets     int32
	Cores       int32
	Threads     int32

	Model   [32]int8
	Domains []libvirt.Domain
}

func (u *Usecase) Get() (Hypervisor, error) {
	v, err := u.l.ConnectGetLibVersion()
	if err != nil {
		return Hypervisor{}, fmt.Errorf("failed to get libvirt version: %v", err)
	}

	model, mem, cups, mhz, nodes, sockets, cores, threads, err := u.l.NodeGetInfo()
	if err != nil {
		return Hypervisor{}, fmt.Errorf("failed to get node info: %v", err)
	}

	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := u.l.ConnectListAllDomains(1, flags)
	if err != nil {
		return Hypervisor{}, fmt.Errorf("failed to list domains: %w", err)
	}

	return Hypervisor{
		Model:   model,
		Mem:     mem,
		Cpus:    cups,
		Mhz:     mhz,
		Nodes:   nodes,
		Sockets: sockets,
		Cores:   cores,
		Threads: threads,

		LibcVersion: v,
		Domains:     domains,
	}, nil
}

// создание диска
// создание виртуалки
// Saga ?

var wsport = 39000

func (u *Usecase) Create() error {

	//disk, err := configurator.DiskXML()
	//if err != nil {
	//	return fmt.Errorf("failed to create disk: %w", err)
	//}
	//

	dxml, err := configurator.VMXml(wsport)
	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}
	//
	//u.l.
	//
	//	//u.l.DomainDefineXML(domXml)

	//dxml, err := os.ReadFile("./configs/dxml.xml")
	//if err != nil {
	//	return fmt.Errorf("failed to read dxml: %w", err)
	//}

	_, err = u.l.DomainCreateXML(string(dxml), 0)
	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}

	wsport++

	return nil
}

func (u *Usecase) Delete(uuid libvirt.UUID) error {

	rDom, err := u.l.DomainLookupByUUID(uuid)

	if err != nil {
		return fmt.Errorf("failed to lookup domain: %w", err)
	}

	if err = u.l.DomainDestroy(rDom); err != nil {
		return fmt.Errorf("failed to destroy domain: %w", err)
	} // stop

	if err = u.l.DomainUndefine(rDom); err != nil {
		log.Printf("failed to undefine domain: %v", err)
	}

	return nil
}

type VMInfo struct {
	WSVNCPort int

	State       uint8
	Mem, MaxMem uint64
	VCpu        uint16
	CpuTime     uint64

	RawXML string
}

func (u *Usecase) GetVMInfo(uuid libvirt.UUID) (VMInfo, error) {
	rDom, err := u.l.DomainLookupByUUID(uuid)
	if err != nil {
		return VMInfo{}, fmt.Errorf("failed to lookup domain: %w", err)
	}

	a, maxMem, mem, vCpu, cpuTime, err := u.l.DomainGetInfo(rDom)
	if err != nil {
		return VMInfo{}, fmt.Errorf("failed to get domain info: %w", err)
	}

	rXml, err := u.l.DomainGetXMLDesc(rDom, 1)
	if err != nil {
		return VMInfo{}, fmt.Errorf("failed to get domain xml: %w", err)
	}

	xmlDomain := libvirtxml.Domain{}
	if err = xmlDomain.Unmarshal(rXml); err != nil {
		return VMInfo{}, fmt.Errorf("failed to unmarshal domain xml: %w", err)
	}

	wport := 0
	for _, v := range xmlDomain.Devices.Graphics {
		if v.VNC != nil {
			wport = v.VNC.WebSocket
		}
	}

	return VMInfo{
		wport,
		a, mem, maxMem, vCpu, cpuTime,
		rXml,
	}, nil
}
