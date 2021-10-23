package lib

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func GetCity(ip string) (*string, error) {
	geoipFilePath := "../resources/GeoLite2-City.mmdb"
	db, err := geoip2.Open(geoipFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ipAddr := net.ParseIP(ip)
	record, err := db.City(ipAddr)
	if err != nil {
		return nil, err
	}
	city := record.City.Names["en"]
	country := record.Country.Names["en"]
	c := fmt.Sprintf("%s, %s", city, country)
	return &c, nil
}
