package main

import (
	"flag"
	"fmt"
	"fwknop_tunnel/version"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

type OperationMode int

//const version = "1.0.0"
//const buildDateTime = "2023/08/08"

const (
	SingleTunnel OperationMode = iota
	MultiTunnel
)

var gTunnelId uint32 = 1
var gOM OperationMode = MultiTunnel

func (s OperationMode) String() string {
	switch s {
	case SingleTunnel:
		return "Single Tunnel"
	case MultiTunnel:
		return "Multi Tunnel"
	}
	return "Unknown"
}

var bindAddress = flag.String("bind-address", "127.0.0.1", "bind address")
var localPort = flag.Int("local-port", 6000, "local port")
var remotePort = flag.Int("remote-port", 0, "remote port")
var remoteHost = flag.String("remote-host", "", "remote host")
var delay = flag.Int("delay", 1000, "Time to wait in ms after running the cmd, before it tries to connect")
var verbose = flag.Bool("verbose", false, "verbose logs")
var versionFlag = flag.Bool("version", false, "print version")
var fwknopCmd = flag.String("cmd", "", "The fwknop command to execute before connecting to remote")
var mode = flag.String("mode", "multi", "Tunnel operation mode (values: single, multi)")

func main() {

	log.SetFlags(log.Lmicroseconds)

	//log.Printf("Version: %s, Built: %s", version, buildDateTime)
	log.Printf("Version: %s", version.BuildVersion())

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...]\n", os.Args[0])
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Options:\n")
		fmt.Fprint(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		os.Exit(0)
	}

	switch {
	case *mode == "single":
		gOM = SingleTunnel
	case *mode == "multi":
		gOM = MultiTunnel
	default:
		flag.PrintDefaults()
		log.Fatal("Invalid Operation mode!")
	}

	log.Printf("Operation mode: %s", gOM)

	mainLoop()

}

func mainLoop() {

	var listener net.Listener = nil
	//make sure we close the listener socket as well
	defer func() {
		if listener != nil {

			if *verbose {
				fmt.Println("Closing listener socket")
			}
			listener.Close()
		}
	}()

	for {

		var err error

		if listener == nil {

			listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddress, *localPort))
			if *verbose {
				log.Println("Creating listener socket")
			}
			if err != nil {
				log.Fatal("Failed to listen on local port: ", err)
			}

		}

		log.Printf("Listening on %s:%d", *bindAddress, *localPort)

		localConn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection: ", err)
		}

		if gOM == SingleTunnel {
			listener.Close()
			listener = nil
		}

		if *verbose == true {
			log.Println("S######################################## ID:" + fmt.Sprintf("%08X", gTunnelId))
			log.Printf("Connection from: %s", localConn.RemoteAddr())
		}

		if gOM == SingleTunnel {
			//Single mode, we will run the tunnel and not accept new connections until old tunnel is closed
			createAndRunTunnel(localConn, gTunnelId)
		} else {
			//Multi mode, run the tunnel in a goroutine so we can create many parrallel tunnels
			go createAndRunTunnel(localConn, gTunnelId)
		}

		gTunnelId++

	}
}
func runCmd() bool {
	if *verbose == true {
		log.Printf("Running command %s", *fwknopCmd)
	}

	cmd := exec.Command(*fwknopCmd)

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			//return exitError.ExitCode()
			log.Printf("Command exited with exit code: %d", exitError.ExitCode())
			return false
		} else {

			log.Println("WARNING: Cmd could not be run")
			return false
		}
	} else if *verbose == true {
		log.Println("Cmd run successfully")
	}
	return true
}

func createAndRunTunnel(localConn net.Conn, tunnelId uint32) {

	var wg sync.WaitGroup

	//defer if *verbose == true { log.Println("E######################################## " + fmt.Sprintf("%08X", tunnelId)) }
	defer func() {
		if *verbose == true {
			log.Println("E######################################## ID:" + fmt.Sprintf("%08X", tunnelId))
		}
	}()

	if *fwknopCmd != "" {
		if runCmd() == false {
			localConn.Close()
			log.Println("WARNING: Local socket closed as we could not run cmd")
			return
		}
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}

	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *remoteHost, *remotePort))

	if err != nil {
		localConn.Close()
		log.Println("WARNING: Failed to connect to remote host: ", err)
		return
	}

	if *verbose == true {
		log.Printf("Connected successfully to %s:%d", *remoteHost, *remotePort)
	}

	log.Println("Tunnel created succesfully ID:" + fmt.Sprintf("%08X", tunnelId))

	//we will start 2 go routines, one for each direction
	wg.Add(2)

	go readLoop(localConn, remoteConn, &wg)
	go readLoop(remoteConn, localConn, &wg)

	//wait until both goroutines closed
	wg.Wait()

	if *verbose == true {
		log.Println("Tunnel closed succesfully")

	}
}

//Fast Copy all data from socket 'from' to 'to'
func readLoop(from net.Conn, to net.Conn, wg *sync.WaitGroup) {

	defer wg.Done()
	defer from.Close()
	defer to.Close()

	if _, err := io.Copy(from, to); err != nil {
		return
	}
}
