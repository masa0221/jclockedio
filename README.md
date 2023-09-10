This utility is for clocking in and out on Jobcan, and it sends the results to Chatwork.

"jclockedio" means Jobcan clocked in/out.

## Installation
### macOS
```shell
brew tap masa0221/tap
```

```shell
brew install jclockedio
```

Note: jclockedio requires Chrome and Chromedriver.
```shell
brew install --cask chromedriver
```
Alternatively, you can use Docker.


## Usage
### 1. Configuration
```shell
jclockedio configure
```
You can modify the output format. For more details, run `jclockedio configure --help`.


### 2. Clocking in / out
```shell
jclockedio adit
```
To only verify login without clocking in/out, use the `--no-adit` option.


## Docker Usage
If you don't have Chrome and Chromedriver installed, you can use Docker. To pull the Docker image, check [this link](https://github.com/masa0221/jclockedio/pkgs/container/jclockedio).

```shell
docker build -t jclockedio .
docker run --rm -it jclockedio
```
Default timezone is Asia/Tokyo;to change it, use `-e TZ=UTC`.

### Example
```shell
docker run --rm -it -v "$HOME"/.jclockedio:/root/.jclockedio jclockedio configure
```


## Disclaimer
USE AT YOUR OWN RISK. Make sure to verify the execution process and the results before relying on it.
Good luck with your remote work!🌸   
