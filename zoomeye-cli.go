package main

import (
  "flag"
  "fmt"
  "os"
  . "zoomeye-cli/api"
  . "zoomeye-cli/noapi"
)

func main() {
  key := flag.String("init", "none", "your zoomeye api key")
  ip := flag.String("ip", "none", "Ip to search")
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
    x := ApiCall(ip)
    ParseApi(x)
    os.Exit(0)
  }
  if *info {
    GetInfo()
    os.Exit(0)
  }
}
