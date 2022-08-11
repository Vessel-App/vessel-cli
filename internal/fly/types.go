package fly

import "github.com/umahmood/haversine"

/*****************
 * APP
****************/

type App struct {
	AppName      string       `json:"name"`
	Organization Organization `json:"organization"`
	IpAddresses  IpAddresses  `json:"ipAddresses"`
}

/*****************
 * Machine
****************/

type Machine struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	Region string `json:"regions"`
	Image  string `json:"image"`
}

/*****************
 * IP ADDRESSES
****************/

type IpAddressAllocation struct {
	IpAddress IpAddress `json:"ipAddress"`
}

type IpAddresses struct {
	Nodes []IpAddress `json:"nodes"`
}

type IpAddress struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Region  string `json:"region"`
}

/*****************
 * REGION
****************/

type NearestRegion struct {
	NearestRegion Region `json:"nearestRegion"`
}

type Region struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	GatewayAvailable bool   `json:"gatewayAvailable"`
	Location         haversine.Coord
	UserDistance     float64
}

/*****************
 * USER/ORG
****************/

type User struct {
	Email         string        `json:"email"`
	Organizations Organizations `json:"organizations"`
}

type Organizations struct {
	Nodes []Organization `json:"nodes"`
}

type Organization struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Type string `json:"type"`
	Role string `json:"viewerRole"`
}
