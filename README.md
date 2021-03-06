# Zoomeye cli

This can be use with **or** without api key.

For usage without api chrome or chromiu is required (will run in headless mode), It also tke more time (aroude 20/25 seconds)

hostname search will come later.

## Usage

```txt
  -k|--init    KEY             Setup your zoomeye api key
  -i|--ip      IP              Search IP
  -d|--domain  DOMAIN          Search DOMAIN
  -c|--cidr    CIDR            Search CIDR
  -f|--info                    Info about your account
  -n|--noapi                   Search without an api key
  -h|--help                    Show this help message and exit
```

## API usages:

![api](.github/api.png)

#### Init

first add your api key with:

```sh 
zoomeye-cli --init "<API KEY>"

```

API key will be stored in `~/.zoomeye`

#### Info

You can check your credit and role with:

```sh
zoomeye-cli --info
Account: user
Credits: 8559
```

#### Search

##### ip 

```sh
zoomeye-cli -i 1.1.1.1
```

##### Domain

```sh
zoomeye-cli -d 'archlinux.org'
```

![domai](.github/domain.png)


##### cidr

```sh
zoomeye-cli -c 1.1.1.1/24
```


## WIthout API

![noapi](.github/noapi.png)

You need to have chrome or chromium for that, it will use chrome in headless mode.

It still have some bugs sometimes and it's lot longer than with the api (around 20/25seconds).


Only ip search is supported now...


Noapi add [seebug]("https://www.seebug.org") vuln history.

- green: low
- yellow: medium 
- red: hight

for check details go to "https://www.seebug.org/vuldb/ssvid-<ID>"



#### Search

```
zoomeye-cli --noapi -i <ip>
```

## Installation

```sh
git clone git@github.com:mmpx12/zoomeye-cli.git
cd zoomeye-cli

make 
sudo make install 
# or 
sudo make all
```



