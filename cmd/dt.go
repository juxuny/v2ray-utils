package cmd

import (
	"fmt"
	"strconv"
)

type H map[string]interface{}

func NewString(s interface{}) *string {
	var ret string
	if s == nil {
		return nil
	}
	switch s.(type) {
	case string:
		ret = s.(string)
	default:
		ret = fmt.Sprintf("%v", s)
	}
	return &ret
}

func NewInt(v interface{}) *int {
	var ret int
	if v == nil {
		return nil
	}
	switch v.(type) {
	case int:
		ret = v.(int)
	default:
		value, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
		ret = int(value)
	}
	return &ret
}

type User struct {
	User     *string `json:"user,omitempty"`
	Pass     *string `json:"pass,omitempty"`
	Id       *string `json:"id,omitempty"`
	AlterId  *int    `json:"alterId,omitempty"`
	Security *string `json:"security,omitempty"`
	Level    *int    `json:"level,omitempty"`
}

type Server struct {
	Address *string `json:"address,omitempty"`
	Port    *int    `json:"port,omitempty"`
	Users   []User  `json:"users,omitempty"`
}

type Vnext struct {
	Users   []User  `json:"users,omitempty"`
	Address *string `json:"address,omitempty"`
	Port    *int    `json:"port,omitempty"`
}

type StreamSettings struct {
	Network    *string `json:"network,omitempty"`
	Security   *string `json:"security,omitempty"`
	WsSettings struct {
		Path    *string `json:"path,omitempty"`
		Headers H       `json:"headers,omitempty"`
	} `json:"wsSettings,omitempty"`
	QuicSettings struct {
		Key      *string `json:"key,omitempty"`
		Security *string `json:"security,omitempty"`
		Header   H       `json:"header,omitempty"`
	} `json:"quicSettings,omitempty"`
	TlsSettings struct {
		AllowInsecure        bool     `json:"allowInsecure"`
		Alpn                 []string `json:"alpn,omitempty"`
		ServerName           *string  `json:"serverName,omitempty"`
		AllowInsecureCiphers bool     `json:"allowInsecureCiphers"`
	} `json:"tlsSettings,omitempty"`
	HttpSettings struct {
		Path *string  `json:"path,omitempty"`
		Host []string `json:"host,omitempty"`
	}
	TcpSettings struct {
		Header H `json:"header,omitempty"`
	} `json:"tcpSettings,omitempty"`
	KcpSettings struct {
		Header           H    `json:"header,omitempty"`
		Mtu              int  `json:"mtu,omitempty"`
		Congestion       bool `json:"congestion,omitempty"`
		Tti              int  `json:"tti,omitempty"`
		UplinkCapacity   int  `json:"uplinkCapacity,omitempty"`
		WriteBufferSize  int  `json:"writeBufferSize,omitempty"`
		ReadBufferSize   int  `json:"readBufferSize,omitempty"`
		DownlinkCapacity int  `json:"downlinkCapacity,omitempty"`
	} `json:"kcpSettings,omitempty"`
}

type Rule struct {
	Type        *string  `json:"type,omitempty"`
	OutboundTag *string  `json:"outboundTag,omitempty"`
	Ip          []string `json:"ip,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	Port        *string  `json:"port,omitempty"`
	BalancerTag *string  `json:"balancerTag,omitempty"`
}

type Routing struct {
	Name           *string `json:"name,omitempty"`
	DomainStrategy *string `json:"domainStrategy,omitempty"`
	Rules          []Rule  `json:"rules,omitempty"`
	Balancers      []struct {
		Tag      *string  `json:"tag,omitempty"`
		Selector []string `json:"selector"`
		Strategy *string  `json:"strategy,omitempty"`
		Interval *int     `json:"interval,omitempty"`
		Timeout  *int     `json:"timeout,omitempty"`
	} `json:"balancers,omitempty"`
}

type Outbound struct {
	Tag      *string `json:"tag,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	Settings struct {
		DomainStrategy *string `json:"domainStrategy,omitempty"`
		Response       struct {
			Type *string `json:"type,omitempty"`
		} `json:"response,omitempty"`

		Network *string  `json:"network,omitempty"`
		Address *string  `json:"address,omitempty"`
		Port    *int     `json:"port,omitempty"`
		Servers []Server `json:"servers,omitempty"`
		Vnext   []Vnext  `json:"vnext,omitempty"`
	} `json:"settings,omitempty"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
}

type Inbound struct {
	Listen   *string `json:"listen,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	Tag      *string `json:"tag,omitempty"`
	Settings struct {
		Auth    *string `json:"auth,omitempty"`
		Udp     bool    `json:"udp,omitempty"`
		Ip      *string `json:"ip,omitempty"`
		Timeout int
	} `json:"settings,omitempty"`
}

type Dns struct {
	Servers []string `json:"servers,omitempty"`
}

type Config struct {
	Log *struct {
		LogLevel *string `json:"logLevel,omitempty"`
	} `json:"log,omitempty"`
	Outbounds []Outbound `json:"outbounds,omitempty"`
	Inbounds  []Inbound  `json:"inbounds,omitempty"`
	Dns       *Dns       `json:"dns"`
	Routing   *Routing   `json:"routing,omitempty"`
}
