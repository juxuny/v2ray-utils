package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var mergeArgs = struct {
	Template  string
	VmessFile string
	Config    string
}{}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge vmess file into the template.json",
	Run: func(cmd *cobra.Command, args []string) {
		if mergeArgs.Template == "" {
			er("missing --template argument")
		}
		initTemplate(mergeArgs.Template)
		var data []byte
		var err error
		data, err = getDataFromArgs(mergeArgs.VmessFile)
		if err != nil {
			er(err)
		}
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
				balancerSelector = append(balancerSelector, *obj.Tag)
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
		if err := ioutil.WriteFile(mergeArgs.Config, []byte(ToJson(templateConfig, true)), 0644); err != nil {
			panic(err)
		}
	},
}

func init() {
	mergeCmd.PersistentFlags().StringVar(&mergeArgs.Template, "template", "", "a template json file for v2ray-core")
	mergeCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "display debug output")
	mergeCmd.PersistentFlags().StringVar(&mergeArgs.VmessFile, "vmess", "", "vmess list file")
	mergeCmd.PersistentFlags().StringVar(&mergeArgs.Config, "config", "config.json", "output json file")
	rootCmd.AddCommand(mergeCmd)
}
