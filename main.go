package main

import (
	"net"
	//	"flag"
	"fmt"
	//	"os"
	//	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	debug            bool
	blackTargetList  [1 << 29]uint8
	whiteTargetList  [1 << 29]uint8
	pblackTargetList *[1 << 29]uint8
	pwhiteTargetList *[1 << 29]uint8
	statsListMap     map[string]int
	srcMac           net.HardwareAddr
	dstMac           net.HardwareAddr

//	lock             sync.RWMutex
)

func main() {
	/*
		parseFlags()
		// Exit on invalid parameters
		flagsComplete, errString := flagsComplete()
		if !flagsComplete {
			fmt.Println(errString)
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/
	configMap := initConfig("conf.txt")

	//init nic
	mirrorNetworkDevice := configMap["mirrornetworkdevice"]
	mgtNetworkDevice := configMap["mgtnetworkdevice"]

	//init targets
	targetListPath := configMap["targetlistpath"]
	getTargetList(targetListPath)

	//init nic mac for send packets
	srcMac, dstMac = initMacs(configMap["srcmac"], configMap["dstmac"])
	if (srcMac != nil) && (dstMac != nil) {
		fmt.Println("set src and dst mac: ", srcMac, dstMac)
	}

	initBullets()

	statsListMap = initStats()

	// Open connection
	handleMirror, errHandleMirror := pcap.OpenLive(
		mirrorNetworkDevice, // network device
		int32(65535),
		true,
		time.Microsecond,
	)
	if errHandleMirror != nil {
		fmt.Println("Mirror Handler error", errHandleMirror.Error())
	}

	handleMgt, errHandleMgt := pcap.OpenLive(
		mgtNetworkDevice, // network device
		int32(65535),
		false,
		time.Microsecond,
	)
	if errHandleMgt != nil {
		fmt.Println("MGT Handler error", errHandleMgt.Error())
	}

	//init send packet channel
	c := make(chan [2]gopacket.Packet)
	go sendResetPacket(handleMgt, c)

	//Close when done
	//defer handleMirror.Close()
	//defer handleMgt.Close()

	//Capture Live Traffic
	packetSource := gopacket.NewPacketSource(handleMirror, handleMirror.LinkType())
	for packet := range packetSource.Packets() {
		go analysePacket(packet, handleMgt, c)
	}

}
