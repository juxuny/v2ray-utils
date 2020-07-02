package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	nbUser       string
	nbPassword   string
	nbLogin      bool
	outVmessList string
)

const (
	host      = "https://nbsd.live"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"
)

// 南部隧道客户端
type NClient struct {
	cookies []*http.Cookie
}

func NewNClient() *NClient {
	ret := NClient{}
	return &ret
}

type LoginResp struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

func (t *NClient) loadCookie() error {
	if data, err := ioutil.ReadFile("cookies.json"); err != nil {
		return errors.Wrap(err, "read cookies.json failed")
	} else {
		if err := json.Unmarshal(data, &t.cookies); err != nil {
			return errors.Wrap(err, "invalid cookies data")
		}
	}
	return nil
}

func (t *NClient) genApi(path string) string {
	return fmt.Sprintf("%s%s", host, path)
}

func (t *NClient) Login(userName string, password string, code ...string) error {
	value := url.Values{}
	value.Add("email", userName)
	value.Add("passwd", password)
	if len(code) > 0 {
		value.Add("code", code[0])
	}
	buf := bytes.NewBufferString(value.Encode())
	req, err := http.NewRequest(http.MethodPost, t.genApi("/auth/login"), buf)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "create request failed ")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "login failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "http code = %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read http body error")
	}
	fmt.Println(string(data))
	var loginResp LoginResp
	if err := json.Unmarshal(data, &loginResp); err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "invalid json")
	}
	fmt.Println(resp.Header)
	fmt.Println("ret = ", loginResp.Ret, "resp: ", loginResp.Msg)
	fmt.Println("cookies: ", t.cookies)
	if data, err := json.Marshal(resp.Cookies()); err != nil {
		fmt.Println("save cookies error: ", err)
	} else {
		if err := ioutil.WriteFile("cookies.json", data, 0644); err != nil {
			fmt.Println("save cookies as json error: ", err)
		}
	}
	t.cookies = resp.Cookies()
	return nil
}

type NodeInfo struct {
	Name       string  `json:"name"`
	Load       int64   `json:"load"`
	Connection int64   `json:"connection"`
	Speed      float64 `json:"speed"`
	Enable     bool    `json:"enable"`
	Remain     int64   `json:"remain"`
	Vmess      string  `json:"vmess"`
}

type NodeInfoList []NodeInfo

func (t NodeInfoList) Filter(f func(n NodeInfo) bool) (ret NodeInfoList) {
	for _, item := range t {
		if f(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func (t NodeInfoList) Less(i, j int) bool {
	a := float64(t[i].Connection)*0.7 + float64(t[i].Remain)*0.2
	b := float64(t[j].Connection)*0.7 + float64(t[j].Remain)*0.2
	return a > b
}

func (t NodeInfoList) Len() int {
	return len(t)
}

func (t NodeInfoList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *NClient) FetchNodeList() (list NodeInfoList, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, t.genApi("/user/node"), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request failed")
	}
	req.Header.Set("User-Agent", UserAgent)
	for _, c := range t.cookies {
		req.AddCookie(c)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code: %v", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io error")
	}
	//fmt.Println(string(data))
	_ = ioutil.WriteFile("tmp.html", data, 0644)
	buf := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return nil, errors.Wrap(err, "parse error")
	}
	doc.Find(".tile-wrap > .tile.tile-collapse").Each(func(i int, selection *goquery.Selection) {
		var item NodeInfo
		value := selection.Find(".tile-inner > div.text-overflow.node-textcolor > span.enable-flag").Text()
		item.Name = strings.Trim(value, "\t\r\n ")
		value = selection.Find(".tile-inner > div.text-overflow.node-textcolor > strong > b > span.node-alive").Text()
		value = NumberFilter(value)
		item.Connection, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		// load
		value = selection.Find(".tile-inner > div.text-overflow.node-textcolor > span.node-load").Text()
		if strings.Contains(value, "N/A") {
			item.Load = -1
		} else {
			value = NumberFilter(value)
			item.Load, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				fmt.Println(err)
			}
		}

		// network traffic
		value = selection.Text()
		value = getNetworkTraffic(value)
		item.Remain, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		// status
		value = selection.Find(".tile-inner > div.text-overflow.node-textcolor > span.node-status").Text()
		if strings.Contains(value, "可用") {
			item.Enable = true
		} else {
			item.Enable = false
		}

		// vmess
		node := selection.Find(".collapsible-region.collapse > .tile-sub > .card.nodetip-table > .card-main > .card-inner > p")
		var ok bool
		value, ok = node.Find("a").Attr("data-clipboard-text")
		if ok {
			item.Vmess = value
		}

		list = append(list, item)
		//tmp, _ := json.Marshal(item)
		//fmt.Println(string(tmp))
	})
	list = list.Filter(func(n NodeInfo) bool {
		return n.Remain > 0 && n.Load > 0 && n.Enable
	})
	sort.Sort(list)
	return list, nil
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawl vmess link from https://nbsd.live",
	Run: func(cmd *cobra.Command, args []string) {
		if nbUser == "" || nbPassword == "" {
			er("username and password can't be empty")
		}
		client := NewNClient()
		if _, err := os.Stat("cookies.json"); os.IsNotExist(err) {
			if err := client.Login(nbUser, nbPassword); err != nil {
				panic(err)
			}
		}

		// load cookies from current directory
		if err := client.loadCookie(); err != nil {
			// login again
			fmt.Printf("%+v", err)
			fmt.Println("load cookie failed, login again")
			if err := client.Login(nbUser, nbPassword); err != nil {
				panic(err)
			}
		}

		fmt.Println("crawl node data")
		list, err := client.FetchNodeList()
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		fmt.Println("the number of node: ", len(list))

		if len(list) > 2 {
			list = list[:2]
		} else if len(list) == 0 {
			fmt.Println("list is empty")
			return
		}

		outData := ""
		for _, item := range list {
			fmt.Println(item.Name, item.Vmess)
			outData += item.Vmess + "\n"
		}
		if err := ioutil.WriteFile(outVmessList, []byte(strings.Trim(outData, "\r\n ")), 0644); err != nil {
			er(err)
		}
	},
}

func init() {
	crawlCmd.PersistentFlags().StringVar(&nbUser, "user", "", "username of nbsd.live")
	crawlCmd.PersistentFlags().StringVar(&nbPassword, "password", "", "password of nbsd.live")
	crawlCmd.PersistentFlags().BoolVar(&nbLogin, "login", true, "re-login")
	crawlCmd.PersistentFlags().StringVar(&outVmessList, "out", "vmess.list", "output file for vmess list")
	crawlCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "display debug output")
	rootCmd.AddCommand(crawlCmd)
}
