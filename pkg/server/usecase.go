package server

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"log"
	"net"
	"os"
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

func (u *Usecase) Create() error {
	dxml, err := os.ReadFile("./configs/dxml.xml")
	if err != nil {
		return fmt.Errorf("failed to read dxml: %w", err)
	}

	_, err = u.l.DomainCreateXML(string(dxml), 0)
	if err != nil {
		return fmt.Errorf("failed to create domain: %w", err)
	}

	return nil
}

func (u *Usecase) Delete(id int32) error {
	rDom, err := u.l.DomainLookupByID(id)
	if err != nil {
		return fmt.Errorf("failed to lookup domain: %w", err)
	}

	if err = u.l.DomainDestroy(rDom); err != nil {
		return fmt.Errorf("failed to destroy domain: %w", err)
	}

	return nil
}
