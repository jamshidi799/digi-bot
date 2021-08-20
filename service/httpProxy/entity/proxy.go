package entity

type Proxies struct {
	List []Proxy `njson:"data"`
}

type Proxy struct {
	Ip   string `njson:"ip"`
	Port string `njson:"port"`
}
