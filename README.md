# master-base
Master Base Station has been built for the [RVASec](https://rvasec.com/) conference by members at [HackRVA](https://hackrva.org).

The Master Base Station communicates with the RVASec conference attendees' badges over IR to schedule LaserTag Games.  Additionally, The Base Station can help attendees transmit code to their badge.

## Connecting A badge over serial

The Base station expects a badge to be connected over usb on `/dev/ttyACM0`

## Config File
A config file can be created here: `/etc/basestation/baseconfig.yaml`

This is were you can override variables such as:
```
leaderBoard_API: "http://192.168.1.2:5000/api/"
serialPort: /dev/ttyACM0
```
note: default values exist if a config file is not present.


## Start Up Script
The start up script will compile and run a binary called `basestation`
```
$ sh build_and_run.sh
```