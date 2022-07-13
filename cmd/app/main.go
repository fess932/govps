package main

import (
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("Hello World")

	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}

	l := libvirt.NewWithDialer(dialers.NewAlreadyConnected(c))
	if err = l.Connect(); err != nil {
		log.Fatalf("failed to connect to libvirt: %v", err)
	}

	v, err := l.ConnectGetLibVersion()
	if err != nil {
		log.Fatalf("failed to get libvirt version: %v", err)
	}
	fmt.Println("Version:", v)

	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := l.ConnectListAllDomains(1, flags)
	if err != nil {
		log.Fatalf("failed to list domains: %v", err)
	}

	fmt.Println("ID\tName\tUUID")
	fmt.Println("----------------------------------------------------")
	for _, d := range domains {
		fmt.Printf("%d\t%s\t%s\n", d.ID, d.Name, d.UUID)
	}

	if od, err := l.DomainLookupByName("demo2"); err != nil {
		log.Printf("failed to lookup domain: %v", err)
	} else {
		if err = l.DomainDestroy(od); err != nil {
			log.Printf("cant destroy domain: %v", err)
		}
	}

	rDom, err := l.DomainCreateXML(dxml, 0)
	if err != nil {
		log.Fatalf("failed to create domain: %v", err)
	}
	fmt.Printf("%d\t%s\t%s\n", rDom.ID, rDom.Name, rDom.UUID)

	if err = l.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect from libvirt: %v", err)
	}

}

var dxml = `
<domain type='kvm'>
  <name>demo2</name>
  <uuid>4dea24b3-1d52-d8f3-2516-782e98a23fa0</uuid>
  <memory>131072</memory>
  <vcpu>1</vcpu>
  <os>
    <type arch="x86_64">hvm</type>
  </os>
  <clock sync="localtime"/>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
    <disk type='file' device='disk'>
     <source file='/var/lib/libvirt/images/demo2.img'/>
     <target dev='hda'/>
    </disk>
    <interface type='network'>
      <source network='default'/>
      <mac address='24:42:53:21:52:45'/>
    </interface>
    <graphics type='vnc' port='-1' keymap='de'/>
  </devices>
</domain>
`
