This is a things for clocked in/out to Jobcan, then it results are send to Chatwork.

"jclockedio" means Jobcan clocked in/out.

## How to use
### 1. Build image
```
$ docker build -t jclockedio .
```

### 2. Run container
Set your environments, if you think requirement it.

.zprofile
```
JCLOCKEDIO_JOBCAN_USERNAME="Your Jobcan username"
JCLOCKEDIO_JOBCAN_PASSWORD="Your Jobcan password"
JCLOCKEDIO_CHATWORK_API_TOKEN="Your Chatwork api token"
JCLOCKEDIO_CHATWORK_ROOM_ID=your room id
```

Execute clocked in/out on Jobcan.
```
$ docker run \
    -e "TZ=Asia/Tokyo" \
    -e "JCLOCKEDIO_JOBCAN_DEBUG=true" \
    -e "JCLOCKEDIO_CHATWORK_DEBUG=true" \
    -e "JCLOCKEDIO_JOBCAN_USERNAME=${JCLOCKEDIO_JOBCAN_USERNAME}" \
    -e "JCLOCKEDIO_JOBCAN_PASSWORD=${JCLOCKEDIO_JOBCAN_PASSWORD}" \
    -e "JCLOCKEDIO_CHATWORK_API_TOKEN=${JCLOCKEDIO_CHATWORK_API_TOKEN}" \
    -e "JCLOCKEDIO_CHATWORK_ROOM_ID=${JCLOCKEDIO_CHATWORK_ROOM_ID}" \
    jclockedio
```
You can setting your credentials and room id for sending to chatwork.
IF YOU WANT TO USE SERIOUSLY, THEN REMOVE DEBUG ENVIRONMENT PARAMETERS `JCLOCKEDIO_JOBCAN_DEBUG` and `JCLOCKEDIO_CHATWORK_DEBUG`.


## How to execute edited source code
If you want to execute edited code, run command below.

### Start docker container & Run container
```
$ docker run -it --name jclockedio -e "TZ=Asia/Tokyo" -v $(pwd):/go/src/app jclockedio go run /go/src/app/cmd/jclockedio/main.go
```
If you are not build jclockedio yet, build jclockedio. (Please see first step)

### Stop container & Remove container
```
$ docker stop jclockedio && docker rm jclockedio
```

## NOTE
IF YOU CLOUDN'T CLOCKED IN/OUT, I'M CAN'T HAVE ANY RESPONSIBIRITY.
YOU MUST CONFIRM THIS EXECUTION PROCESS AND EXECUTED RESULTS, IF YOU WANT TO USE IT.
I WISH YOU GOOD LUCK FOR YOUR REMOTE WORK!
