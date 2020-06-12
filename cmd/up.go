package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

var upArgs = struct {
	Template     string
	VmessFile    string
	Cache        bool // true: check cache data for subscribe data
	SubscribeUrl string
	Config       string
	IgnoreTag    string
	IgnoreAddr   string
}{}

func init() {
	upCmd.PersistentFlags().StringVar(&upArgs.Template, "template", "", "a template json file for v2ray-core")
	upCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "display debug output")
	upCmd.PersistentFlags().StringVar(&upArgs.VmessFile, "vmess", "", "vmess list file")
	upCmd.PersistentFlags().BoolVar(&upArgs.Cache, "cache", false, "cache data from subscribe URL")
	upCmd.PersistentFlags().StringVar(&upArgs.SubscribeUrl, "url", "", "subscribe URL")
	upCmd.PersistentFlags().StringVar(&upArgs.IgnoreTag, "ignore-tag", "", "a list of keyword, used to filter by tag")
	upCmd.PersistentFlags().StringVar(&upArgs.IgnoreAddr, "ignore-addr", "", "a list of keyword, used to filter by addr")
	upCmd.PersistentFlags().StringVar(&upArgs.Config, "config", "config.json", "output json file")
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "update v2ray subscription",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updating from subscription URL")
		initTemplate()
		var data []byte
		var err error
		if upArgs.VmessFile == "" {
			data, err = get(upArgs.SubscribeUrl)
			if err != nil {
				panic(err)
			}
		} else {
			data, err = getDataFromArgs(upArgs.VmessFile)
		}

		logf("response data: %s", string(data))
		data, err = base64Decode(string(data))
		if err != nil {
			panic(err)
		}
		l := strings.Split(string(data), "\n")
		log("server list: ", l)
		outbounds := make([]Outbound, 0)
		balancerSelector := make([]string, 0)
		for _, i := range l {
			if i == "" {
				continue
			}
			obj := parse(i)
			if obj.Settings.Address != nil && *obj.Settings.Port != 0 {
				outbounds = append(outbounds, obj)
				if !isIgnore(*obj.Tag, strings.Split(upArgs.IgnoreTag, ",")) && !isIgnore(*obj.Settings.Address, strings.Split(upArgs.IgnoreAddr, ",")) {
					balancerSelector = append(balancerSelector, *obj.Tag)
				} else {
					log("ignore server: ", *obj.Tag)
				}
			}
		}
		templateConfig.Outbounds = append(templateConfig.Outbounds, outbounds...)
		if len(templateConfig.Routing.Balancers) > 0 {
			for i := range templateConfig.Routing.Balancers {
				if templateConfig.Routing.Balancers[i].Tag != nil && *templateConfig.Routing.Balancers[i].Tag == "Balancer" {
					templateConfig.Routing.Balancers[i].Selector = append(templateConfig.Routing.Balancers[i].Selector, balancerSelector...)
				}
			}
		}
		if err := ioutil.WriteFile(upArgs.Config, []byte(ToJson(templateConfig, true)), 0644); err != nil {
			panic(err)
		}
	},
}

var (
	templateConfig Config
)

func get(u string) ([]byte, error) {
	if upArgs.Cache {
		if data, err := ioutil.ReadFile(".cache.txt"); err == nil {
			return data, nil
		}
	}
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("bad request, http status code: %v", resp.StatusCode))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		_ = ioutil.WriteFile(".cache.txt", data, 0644)
	}
	return data, err
}

func parseVmess(data []byte) Outbound {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		panic(err)
	}
	ret := Outbound{
		Tag:            NewString(obj["ps"].(string)),
		Protocol:       NewString("vmess"),
		StreamSettings: &StreamSettings{},
	}
	ret.Settings.Address = NewString(obj["add"])
	ret.Settings.Port = NewInt(obj["port"])
	ret.Settings.Vnext = []Vnext{
		{
			Address: NewString(obj["add"]),
			Port:    NewInt(obj["port"]),
			Users:   []User{{Id: NewString(obj["id"]), AlterId: NewInt(obj["aid"]), Security: NewString("none"), Level: NewInt(0)}},
		},
	}
	ret.Tag = NewString(obj["ps"])
	ret.StreamSettings.WsSettings.Path = NewString(obj["path"])
	ret.StreamSettings.WsSettings.Headers = H{
		"Host": NewString(obj["add"]),
	}
	ret.StreamSettings.Network = NewString(obj["net"])
	return ret
}

func parse(uri string) (ret Outbound) {
	log(uri)
	uri = strings.Trim(uri, " ")
	sl := strings.Split(uri, ":")
	protocol := sl[0]
	content := strings.Trim(sl[1], "/")
	c, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		panic(err)
	}
	log(protocol, string(c))
	switch protocol {
	case "vmess":
		return parseVmess(c)
	default:
		panic("unknown protocol: " + protocol)
	}
	return
}

func base64Decode(s string) ([]byte, error) {
	if strings.Contains(s, "=") {
		return base64.StdEncoding.DecodeString(s)
	}
	return base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(s)
}

func getDataFromArgs(vmessFile string) (ret []byte, err error) {
	data, err := ioutil.ReadFile(vmessFile)
	if err != nil {
		panic(err)
	}
	ret = make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(ret, data)
	return
}

func initTemplate() {
	if upArgs.Template == "" {
		er("missing --template argument")
	}
	if upArgs.SubscribeUrl == "" {
		er("missing --url argument")
	}
	data, err := ioutil.ReadFile(upArgs.Template)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &templateConfig); err != nil {
		panic(err)
	}
	//log(ToJson(templateConfig, true))
}
