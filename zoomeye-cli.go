package main

import (
  "flag"
  "fmt"
  "os"
  . "zoomeye-cli/api"
  . "zoomeye-cli/cidr"
  . "zoomeye-cli/noapi"
)

func main() {
  key := flag.String("init", "none", "your zoomeye api key")
  ip := flag.String("ip", "none", "Ip to search")
  cidr := flag.String("cidr", "none", "Cidr to scan")
  info := flag.Bool("info", false, "info about your account")
  noapi := flag.Bool("noapi", false, "Use zoomeye without api key")
  flag.Parse()
  if *noapi == true {
    port := true
    if *ip == "none" {
      fmt.Println("You should provide an ip.\nex: zoomeye -noapi -ip 1.1.1.1")
      os.Exit(1)
    }
    PrintVuls(port, *ip)
    os.Exit(0)
  }
  if *key != "none" {
    CreateApiFile(key)
    os.Exit(0)
  }
  if *ip != "none" {
    x, y := ApiCall(*ip)
    if x {
      ParseApi(y)
    }
    os.Exit(0)
  }
  if *cidr != "none" {
    ips := Cidr_to_ip(*cidr)
    for _, ip := range ips {
      fmt.Println("\033[38;5;118mIP:       \t\033[33m" + ip)
      x, y := ApiCall(ip)
      if x {
        ParseApi(y)
      }
      fmt.Println("\n")
    }
  }
  if *info {
    GetInfo()
    os.Exit(0)
  }
}
