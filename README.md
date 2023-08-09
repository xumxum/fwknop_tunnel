# fwknop_tunnel

## Description
fwknop_tunnel is a simple Tcp Tunnel written in go, that will forward a local TCP connection to a remote server, but only after running a shell script, usually a fwknop command to open the port on the remote server.  
Since most of the VPNs,clients don't have a native support for fwknop, redirecting them to use this tunnel will add transparent support.

No matter how updated and secure you keep your server, the best way to keept it secure is by not having open ports and only opening them when needed before connection from a client. 

Use cases: 
- bypass company firewalls which limit outgoing ports
- hide any tcp based VPN ports (ex tinc) 
- hide any services (ssh or vpn) behind a working https port. 

## fwknop
fwknop in a nutshell, is a tool that sends Single Packet Authorization to a remote server's fwknop server, which will trigger to open up a port temporarily(that is otherwise firewalled) to be able to connect to it. 

Excelent tutorial about fwknop:  
http://www.cipherdyne.org/fwknop/docs/fwknop-tutorial.html

## Build
`go build`

## Install
`GOBIN=/usr/local/bin/ go install`

## Usage
```
Usage: ./fwknop_tunnel [options...]

Options:

  -bind-address string
    	bind address (default "127.0.0.1")
  -cmd string
    	The fwknop command to execute before connecting to remote
  -delay int
    	Time to wait in ms after running the cmd, before it tries to connect (default 1000)
  -local-port int
    	local port (default 6000)
  -mode string
    	Tunnel operation mode (values: single, multi) (default "multi")
  -remote-host string
    	remote host
  -remote-port int
    	remote port
  -verbose
    	verbose logs

```

### Options:
Most of the options are clear with some notes:
- `cmd` - the name of the bash file to run that has the fwknop command inside. Check examples subfolder.   
  Make sure this bash file is executable, and test it out that it openes the port before running it with fwknop_tunnel.  
  If no `cmd` specified, it will act as a regular tcp tunnel, it will also not wait the `delay`

- `delay` - Since it takes maybe couple of hundreds of miliseconds for the SPA packet to arrive and be processed on the remote server, this is the delay it will wait before trying to connect to the remote side to give it time for the packet to be processed and port opened. For safety 1000ms, but you can probably reduce it if this is an issue to you.
  
- `mode` :
  - `single` - It will allow only one tcp connection at a time. Only after the tunnel is closed by either side, the listening port will be reopened and new client request will be acceped. Usually VPN have one tcp at a time(ex tinc), so mostly for that
  - `multi` - Support as many tcp connections as you want. Before each connection, it will rerun the cmd again. For ex for ssh connections where you can have multiple ssh connections at a time. 

Check also `examples` subfolder

## Benchmark
iperf results through the fwknop_tunnel, first result is through fwknop_tunnel, second is directly, on the localhost. 
```
[  1] local 127.0.0.1 port 5001 connected with 127.0.0.1 port 43278
[ ID] Interval       Transfer     Bandwidth
[  1] 0.0000-10.0033 sec  44.0 GBytes  37.8 Gbits/sec
[  2] local 127.0.0.1 port 5001 connected with 127.0.0.1 port 54686
[ ID] Interval       Transfer     Bandwidth
[  2] 0.0000-10.0015 sec  52.8 GBytes  45.3 Gbits/sec
[  3] local 127.0.0.1 port 5001 connected with 127.0.0.1 port 51100
```