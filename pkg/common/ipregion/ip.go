package ipregion

import (
	"embed"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
)

//go:embed db/*
var xdbFile embed.FS

var searcher *xdb.Searcher

type IpRegion struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"` //
	Ip       string `json:"ip"`
}

func Ip2regionSearch(ip string) (IpRegion, error) {
	ir := IpRegion{}
	if searcher == nil {
		file, err2 := xdbFile.ReadFile("db/ip2region.xdb")

		if err2 != nil {
			return ir, err2
		}
		searcher, _ = xdb.NewWithBuffer(file)
	}

	region, err := searcher.SearchByStr(ip)
	if err != nil {
		return ir, err
	}
	split := strings.Split(region, "|")
	ir.Country = split[0]
	ir.Province = split[2]
	ir.City = split[3]
	ir.ISP = split[4]
	ir.Ip = ip
	return ir, nil
}
