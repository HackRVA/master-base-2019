# master-base
Master Base Station has been built for the [RVASec](https://rvasec.com/) conference by members at [HackRVA](https://hackrva.org).

The Master Base Station communicates with the RVASec conference attendees' badges over IR to schedule LaserTag Games.  Additionally, The Base Station can help attendees transmit code to their badge.

## Connecting A badge over serial

The Base station expects a badge to be connected over usb on `/dev/ttyACM0`

## Environment Variables
An Environment Variable is expected for the Leaderboard API URI
```
export LEADERBOARD_API=<leaderboard API URI>
```

If you do not specify the `LEADERBOARD_API` variable, it will default to `http://localhost:5000/api/` 

execute startup script:
> sh build_and_run.sh