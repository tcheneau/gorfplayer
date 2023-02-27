# An API bridge for the RFPlayer

The [RFPlayer](https://www.gce-electronics.com/fr/produits-radio/1777-rf-player-3770008041004.html) is a device built by GCE Electronics similar to the RFXCom. 

This is a device that can talk many wireless protocols in the 433Mghz/868Mghz band. It comes with a java based tool that can send raw commands. The device itself is controled by a serial interface and AT-like commands.

Full specification of the serial protocol can be found [here](https://github.com/gce-electronics/HA_RFPlayer/blob/main/rfplayer_api_v1.8.pdf).

One of my pet project is to automate (through Home Assistant or OpenHAB) the opening/closing of a Somfy RTS Rollershutter. The implementation choice I made here is to make the device available through an HTTP API. One of the min advantage is that the device can be shared with multiple tools.

The device is able to receive events (from Oregon Scientific probes for example). This is currently not supported as the API focuses on sending commands for now.

## How to install

```shell
go install cmd/gorfplayer/gorfplayer.go
```

## How to use

Show the status of the device:
```shell
curl http://localhost:8000/v1/status --output -
```

Read from the device. It has no use for now, but could prove handy if I add rflink support:
```shell
curl http://localhost:8000/v1/read --output -
```

Reply "PONG" and checks that the device still replies:
```shell
curl http://localhost:8000/v1/read --output -
```

The three many commands for opening/closing a Somfy RTS Rollershutter. Please not that A1 is the address I associated the device upon. You can get you RFplayer associated through the java tool (or use the ASSOC commands - though I haven't tested it).

```shell
curl http://localhost:8000/v1/ping --output -
curl -H "Content-Type: application/json" -d '{"order":"off", "protocol": "RTS", "address": "A1"}' http://localhost:8000/v1/command --output -
curl -H "Content-Type: application/json" -d '{"order":"on", "protocol": "RTS", "address": "A1", "burst": "3"}' http://localhost:8000/v1/command --output -
curl -H "Content-Type: application/json" -d '{"order":"dim", "protocol": "RTS", "address": "A1", "burst": "3", "percent" : "100"}' http://localhost:8000/v1/command --output -
```

## Security

This software needs to be run in a *trusted network*.

CI/CD does not provide any signing. You should build your own version, after you read the source code.


## TODO

- add "bearer token" support so as to protect the API from malicious requests
- implement a locking mecanism (serial protocol cannot handle multiple requests in parallel)
