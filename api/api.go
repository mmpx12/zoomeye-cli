package api

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/itchyny/gojq"
  "io/ioutil"
  "net/http"
  "os"
  "os/user"
  "sort"
  "strconv"
  "strings"
)

func GetApiKey() string {
  usr, _ := user.Current()
  dat, _ := ioutil.ReadFile(usr.HomeDir + "/.zoomeye")
  if dat == nil {
    return "false"
  }
  x := string(dat)
  return x
}

func CreateApiFile(api *string) {
  usr, _ := user.Current()
  f, _ := os.Create(usr.HomeDir + "/.zoomeye")
  defer f.Close()
  f.WriteString(*api)
}

func ApiCall(ip *string) []byte {
  apikey := GetApiKey()
  req, _ := http.NewRequest("GET", "https://api.zoomeye.org/host/search?query=cidr:"+*ip+"/32", nil)
  req.Header.Set("API-KEY", apikey)
  client := &http.Client{}
  resp, _ := client.Do(req)
  content, _ := ioutil.ReadAll(resp.Body)
  checkerr := bytes.Contains(content, []byte(`{"error": "invalid_`))
  if checkerr {
    fmt.Println("Error:")
    fmt.Println(string(content))
    os.Exit(1)
  }
  noport := bytes.Contains(content, []byte(`{"total": 0, "available": 0,`))
  if noport {
    fmt.Println("Nothing found.")
    os.Exit(0)
  }
  return content
}

func ParseApi(raw []byte) {
  var input interface{}
  result := make([]string, 0)
  infores := make([]string, 0)
  json.Unmarshal(raw, &input)
  info, _ := gojq.Parse(`(if .matches[0].geoinfo.isp == "" or .matches[0].geoinfo.isp == null then "unknow" else .matches[0].geoinfo.isp end), (if .matches[0].geoinfo.country.names.en == null or .matches[0].geoinfo.country.names.en == "" then "unknow" else .matches[0].geoinfo.country.names.en end), (.matches |length)`)
  val := info.Run(input)
  for {
    v, ok := val.Next()
    if !ok {
      break
    }
    res := fmt.Sprintf("%v\n", v)
    infores = append(infores, res)
  }
  fmt.Printf("\033[38;5;118mCountry:    \t\033[33m%s\033[38;5;118mISP:        \t\033[33m%s\033[38;5;118mOpened port:\t\033[33m%s", infores[1], infores[0], infores[2])
  // ports
  query, _ := gojq.Parse(`.matches[]|[.portinfo.port, (if .protocol.transport == "" or .protocol.transport == null then "unknow" else .protocol.transport end), (if .portinfo.app == null or .portinfo.app == "" then (if .protocol.application == null or .protocol.application == "" or .protocol.application == "test" then .portinfo.service else (.protocol.application|ascii_downcase) end) else .portinfo.app end),( if .portinfo.version == null or .portinfo.version == "" then "unknow" else .portinfo.version end), .timestamp, if (.portinfo.banner | length) < 200 then .portinfo.banner else (.portinfo.banner|.[0:100])+" ..." end]|@tsv`)
  value := query.Run(input)
  for {
    v, ok := value.Next()
    if !ok {
      break
    }
    res := fmt.Sprintf("%v\n", v)
    result = append(result, res)
  }
  sortRes := make(map[int]string, 0)
  for _, d := range result {
    restt := strings.Split(d, "\t")
    port, err := strconv.Atoi(restt[0])
    if err != nil {
      continue
    }
    sortRes[port] = restt[1] + "\t" + restt[2] + "\t" + restt[3] + "\t" + restt[4] + "\t" + restt[5]
  }
  keys := make([]int, 0, len(result))
  for k := range sortRes {
    keys = append(keys, k)
  }
  sort.Ints(keys)
  fmt.Println("\033[38;5;118mPorts:")
  for _, j := range keys {
    finalres := strings.Split(sortRes[j], "\t")
    if finalres[2] == "unknow" {
      fmt.Printf("   \033[33m%d/%s\033[0m \033[31m%s \033[38;5;5m%s\n\t\033[33m╰─▪ \033[36m%s", j, finalres[0], finalres[1], finalres[3], finalres[4])
    } else {
      fmt.Printf("   \033[33m%d/%s\033[0m \033[31m%s\033[38;5;202m (%s) \033[38;5;5m%s\n\t\033[33m╰─▪ \033[36m%s", j, finalres[0], finalres[1], finalres[2], finalres[3], finalres[4])
    }
  }
}

func GetInfo() {
  apikey := GetApiKey()
  req, _ := http.NewRequest("GET", "https://api.zoomeye.org/resources-info", nil)
  req.Header.Set("API-KEY", apikey)
  client := &http.Client{}
  resp, _ := client.Do(req)
  content, _ := ioutil.ReadAll(resp.Body)
  var input interface{}
  result := make([]string, 0)
  json.Unmarshal(content, &input)
  query, _ := gojq.Parse(`[.plan, .resources.search]|@tsv`)
  value := query.Run(input)
  for {
    v, ok := value.Next()
    if !ok {
      break
    }
    res := fmt.Sprintf("%v\n", v)
    result = append(result, res)
  }
  info := strings.Split(result[0], "\t")
  fmt.Print("Account: ", info[0], "\nCredits: ", info[1])
  os.Exit(0)
}
