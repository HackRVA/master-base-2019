# master-base
Master Base Station coding

The Base station expects a badge to be connected over usb on `/dev/ttyACM0`

## Environment Variables
An Environment Variable is expected for the Leaderboard API URI
```
export LEADERBOARD_API=<leaderboard API URI>
```

If you do not specify the `LEADERBOARD_API` variable, it will default to `http://localhost:5000/api/` 

execute startup script:
> sh build_and_run.sh