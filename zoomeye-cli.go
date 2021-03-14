package main

import (
  "fmt"
  "github.com/speedata/optionparser"
  "os"
  . "zoomeye-cli/api"
  . "zoomeye-cli/cidr"
  . "zoomeye-cli/noapi"
)

func main() {
  var key, ip, cidr string
  var noapi, info bool
  op := optionparser.NewOptionParser()
  op.On("-k", "--init KEY", "Setup your zoomeye api key", &key)
  op.On("-i", "--ip IP", "Search IP", &ip)
  op.On("-c", "--cidr CIDR", "Search CIDR", &cidr)
  op.On("-f", "--info", "Info about your account", &info)
  op.On("-n", "--noapi", "Search without an api key", &noapi)
  op.On("-h", "--help", "Show this help", help)
  err := op.Parse()
  if err != nil || len(os.Args) == 1 || len(op.Extra) != 0 {
    help()
  }
  if noapi {
    port := true
    if ip == "" {
      fmt.Println("You should provide an ip.\nex: zoomeye --noapi -i 1.1.1.1")
      os.Exit(1)
    }
    PrintVuls(port, ip)
    os.Exit(0)
  }
  if key != "" {
    CreateApiFile(key)
    os.Exit(0)
  }
  if ip != "" {
    x, y := ApiCall(ip)
    if x {
      ParseApi(y)
    }
    os.Exit(0)
  }
  if cidr != "" {
    ips := Cidr_to_ip(cidr)
    for _, ip := range ips {
      fmt.Println("\033[38;5;118mIP:       \t\033[33m" + ip)
      x, y := ApiCall(ip)
      if x {
        ParseApi(y)
      }
      fmt.Println("\n")
    }
  }
  if info {
    GetInfo()
    os.Exit(0)
  }
}

func help() {
  fmt.Println(`usage:
  -k|--init  KEY              Setup your zoomeye api key
  -i|--ip    IP               Search IP
  -c|--cidr  CIDR             Search CIDR
  -f|--info                   Info about your account
  -n|--noapi                  Search without an api key
  -h|--help                   Show this help

exemples:
zoomeye-cli --ip 1.1.1.1
zoomeye-cli --cidr 1.1.1.1/24
zoomeye-cli --noapi -i 1.1.1.1`)
  os.Exit(1)
}
