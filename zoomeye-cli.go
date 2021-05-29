package main

import (
	"fmt"
	"github.com/mmpx12/optionparser"
	"os"
	"strings"
	. "zoomeye-cli/api"
	. "zoomeye-cli/cidr"
	. "zoomeye-cli/noapi"
)

func main() {
	var key, domain, ip, cidr string
	var noapi, info bool
	op := optionparser.NewOptionParser()
	op.Banner = "Zoomeye-cli\n\nUsage:"
	op.On("-k", "--init KEY", "Setup your zoomeye api key", &key)
	op.On("-i", "--ip IP", "Search IP", &ip)
	op.On("-d", "--domain DOMAIN", "Search Domain", &domain)
	op.On("-c", "--cidr CIDR", "Search CIDR", &cidr)
	op.On("-f", "--info", "Info about your account", &info)
	op.On("-n", "--noapi", "Search without an api key", &noapi)
	op.Exemple("zoomeye-cli --ip 1.1.1.1")
	op.Exemple("zoomeye-cli --cidr 1.1.1.1/24")
	op.Exemple("zoomeye-cli -i 1.1.1.1 --noapi")
	err := op.Parse()
	if err != nil || len(os.Args) == 1 || len(op.Extra) != 0 {
		op.Help()
		os.Exit(1)
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
		t := "singleip"
		x, y := ApiCall(t, ip)
		if x {
			ParseApi(t, y)
		}
		os.Exit(0)
	}
	if cidr != "" {
		ips := Cidr_to_ip(cidr)
		t := "singleip"
		for _, ip := range ips {
			fmt.Println("\033[38;5;118mIP:       \t\033[33m" + ip)
			x, y := ApiCall(t, ip)
			if x {
				ParseApi(t, y)
			}
			fmt.Println("\n")
		}
	}
	if domain != "" {
		t := "domain"
		x, y := ApiCall(t, domain)
		if x {
			doms := DomainList(y)
			domains := strings.Split(strings.TrimSuffix(doms[0], "\n"), "\t")
			for _, j := range domains {
				if j != "" {
					fmt.Println("\033[38;5;118mDomain:   \t\033[33m." + j + ".")
					ParseApi(j, y)
					fmt.Println("\n")
				}
			}
		}
	}
	if info {
		GetInfo()
		os.Exit(0)
	}
}
