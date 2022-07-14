package server

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
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
	Domains     []libvirt.Domain
}

func (u *Usecase) Get() (Hypervisor, error) {
	v, err := u.l.ConnectGetLibVersion()
	if err != nil {
		return Hypervisor{}, fmt.Errorf("failed to get libvirt version: %v", err)
	}

	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := u.l.ConnectListAllDomains(1, flags)
	if err != nil {
		return Hypervisor{}, fmt.Errorf("failed to list domains: %w", err)
	}

	return Hypervisor{
		LibcVersion: v,
		Domains:     domains,
	}, nil
}

func (u *Usecase) Create() {
	//if od, err := u.l.DomainLookupByName("demo2"); err != nil {
	//	log.Printf("failed to lookup domain: %v", err)
	//} else {
	//	if err = u.l.DomainDestroy(od); err != nil {
	//		log.Printf("cant destroy domain: %v", err)
	//	}
	//}

	//rDom, err := l.DomainCreateXML(dxml, 0)
	//if err != nil {
	//	log.Fatalf("failed to create domain: %v", err)
	//}
	//fmt.Printf("%d\t%s\t%s\n", rDom.ID, rDom.Name, rDom.UUID)
	//
	//if err = l.Disconnect(); err != nil {
	//	log.Fatalf("failed to disco	nnect from libvirt: %v", err)
	//}
}
