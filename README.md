This is a things for clocked in/out to Jobcan, then it results are send to Chatwork.

"jclockedio" means Jobcan clocked in/out.

## How to use
### 1. Configure
```
jclockedio configure
```
You can change output format.  
Please check the help.


### 2. Clocked in / out
```
jclockedio adit
```
IF YOU WANT TO CONFIRM IT CAN BE LOGIN ONLY, YOU CAN USE `--no-adit` OPTION.


## Use docker
jclockedio is depends on chrome and chromedriver.  
If you do not install those software on host machine, you can clocked in/out on docker.

### 1. Build image
```
docker build -t jclockedio .
```

### 2. Run container
Execute clocked in/out on Jobcan.
```
docker run --rm -it jclockedio
```
Default timezone is `Asia/Tokyo`

You want to set other timezone, specifiy -e option.
```
docker run --rm -it -e TZ=UTC jclockedio
```

For example:  
It can be specified settings files, like below.
```
docker run --rm -it -v "$HOME"/.jclockedio:/root/.jclockedio jclockedio configure
```


## NOTE
IF YOU CLOUDN'T CLOCKED IN/OUT, I CAN'T HAVE ANY RESPONSIBIRITY.  
YOU MUST CONFIRM THIS EXECUTION PROCESS AND EXECUTED RESULTS, IF YOU WANT TO USE IT.  
I WISH YOU GOOD LUCK FOR YOUR REMOTE WORK!🌸   

