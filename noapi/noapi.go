package noapi

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/chromedp/cdproto/page"
  "github.com/chromedp/chromedp"
  "github.com/chromedp/chromedp/kb"
  "github.com/itchyny/gojq"
  "log"
  "sort"
  "strconv"
  "strings"
  "time"
  . "zoomeye-cli/useragent"
)

var vuls string
var x string
var c string
var detail string

func vulns(ip string) ([]string, []string, []string) {
  usergnt, size := GetUserAgent()
  opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.DisableGPU,
    chromedp.NoFirstRun,
    chromedp.NoDefaultBrowserCheck,
    chromedp.Flag("headless", true),
    chromedp.Flag("ignore-certificate-errors", true),
    chromedp.UserAgent(usergnt),
    chromedp.Flag("window-size", size),
  )
  allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
  defer cancel()
  taskCtx, cancel := chromedp.NewContext(allocCtx)
  defer cancel()

  ct, cancel := context.WithTimeout(taskCtx, time.Second*45)
  defer cancel()
  ctx, cancel := chromedp.NewContext(ct)
  defer cancel()

  const script = `(function(w, n, wn) {
  Object.defineProperty(n, 'webdriver', {
    get: () => false,
  });
  Object.defineProperty(n, 'plugins', {
    get: () => [1, 2, 3],
  });
  Object.defineProperty(n, 'languages', {
    get: () => ['en-US', 'en'],
  });
  w.chrome = {
    runtime: {},
  };
})(window, navigator, window.navigator);`

  var scriptID page.ScriptIdentifier
  chromedp.Run(
    ctx,
    chromedp.ActionFunc(func(ctx context.Context) error {
      scriptID, _ = page.AddScriptToEvaluateOnNewDocument(script).Do(ctx)
      return nil
    }),
    chromedp.Navigate("https://www.zoomeye.org/"),
    chromedp.SendKeys("input.ant-input", kb.End+"ip:"+ip+kb.Enter, chromedp.BySearch),
    chromedp.WaitVisible(`a.search-result-item-title`, chromedp.BySearch),
    chromedp.Evaluate(`window.location.href = "/host/vuls/?ip="+function y(){var i = 0;while (document.querySelectorAll('a.search-result-item-title')[i].text != (document.getElementsByClassName("ant-input")[0].value).split(":")[1] ){ i++;} while (document.querySelectorAll('a.search-result-item-title')[i].href.split("=")[2] == null){i++}var x = document.querySelectorAll('a.search-result-item-title')[i].href.split("=")[2];return x;}()+"&query_type=1"`, &x),
    chromedp.Text("pre", &vuls),
    chromedp.Evaluate(`window.location.href = "/host/details/"+String(document.URL).split("/")[5].split("=")[1].split("&")[0]+"?from=detail"`, &x),
    chromedp.Text("pre", &detail),
  )
  var input interface{}
  json.Unmarshal([]byte(vuls), &input)
  query, _ := gojq.Parse("[.vuls[]]| sort_by(.level | tonumber) | .[] | [.level, .seebug_id]|@csv")
  value := query.Run(input)
  vulns := make([]string, 0)
  for {
    v, ok := value.Next()
    if !ok {
      break
    }
    res := fmt.Sprintf("%v", v)
    vulns = append(vulns, res)
  }
  var input2 interface{}
  dets := make([]string, 0)
  json.Unmarshal([]byte(detail), &input2)
  query2, err := gojq.Parse(`.ports[]|[.port, (if .service == null or .service == "" then "NULL" else .service end),( if .product == null or .product == "" then "NULL" else .product end), (if (.banner | length) < 200 then .banner else (.banner|.[0:200])+" ..." end), (if .timestamp == null or .timestamp == "" then "unknow\n" else .timestamp end)]|@tsv`)
  if err != nil {
    log.Fatal(err)
  }
  value2 := query2.Run(input2)
  for {
    v, ok := value2.Next()
    if !ok {
      break
    }
    res := fmt.Sprintf("%v\n", v)
    dets = append(dets, res)
  }
  qinfo, _ := gojq.Parse(`.geoinfo| .country.names.en, (if .isp == "" then "Unknow" else .isp end), .organization`)
  vinfo := qinfo.Run(input2)
  info := make([]string, 0)
  for {
    v, ok := vinfo.Next()
    if !ok {
      break
    }
    vres := fmt.Sprintf("%v", v)
    info = append(info, vres)
  }

  return dets, vulns, info
}

func PrintVuls(port bool, ip string) {
  x, y, info := vulns(ip)
  fmt.Println("\033[38;5;118mCountry:      \t\033[33m", info[0])
  fmt.Println("\033[38;5;118mIsp:          \t\033[33m", info[1])
  fmt.Println("\033[38;5;118mOrganization: \t\033[33m", info[2])
  fmt.Printf("\033[38;5;118mSeebug id:       ")
  for _, rest := range y {
    turn := strings.Split(rest, ",")
    switch t, _ := strconv.Atoi(turn[0]); {
    case t < 4:
      c = "\033[32m"
    case t < 7:
      c = "\033[33m"
    case t > 7:
      c = "\033[31m"
    }
    fmt.Printf("%s%s  ", c, turn[1])
  }
  if port {
    sortRes := make(map[int]string, 0)
    for _, d := range x {
      restt := strings.Split(d, "\t")
      port, _ := strconv.Atoi(restt[0])
      sortRes[port] = restt[1] + "\t" + restt[2] + "\t" + restt[3] + "\t" + strings.TrimSuffix(restt[4], "\n")
    }
    keys := make([]int, 0, len(x))
    for k := range sortRes {
      keys = append(keys, k)
    }
    sort.Ints(keys)
    fmt.Println("\n\033[38;5;118mPorts:")
    for _, j := range keys {
      finalres := strings.Split(sortRes[j], "\t")
      if finalres[1] == "NULL" {
        fmt.Printf("   \033[33m%d\033[0m/\033[31m%s\033[38;5;5m  %s\n    \033[33m╰─▪ \033[36m%s\n", j, finalres[0], finalres[3], finalres[2])
      } else {
        fmt.Printf("   \033[33m%d\033[0m/\033[31m%s \033[38;5;5m  %s\n    \033[33m╰─▪ \033[36m%s\n", j, finalres[1], finalres[3], finalres[2])
      }
    }
  }
}
