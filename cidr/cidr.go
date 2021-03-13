package cidr

import "net"

func inc(ip net.IP) {
  for j := len(ip) - 1; j >= 0; j-- {
    ip[j]++
    if ip[j] > 0 {
      break
    }
  }
}

func Cidr_to_ip(cidr string) []string {
  ip, ipnet, _ := net.ParseCIDR(cidr)
  var ips []string
  for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
    ips = append(ips, ip.String())
  }
  lenIPs := len(ips)
  switch {
  case lenIPs < 1:
    return ips
  default:
    return ips[1 : len(ips)-1]
  }
}
