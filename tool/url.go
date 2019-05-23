package tool

var domain string

func InitDomain(env string) {
	if env == "dev" {
		domain = "http://devmanager.wb-intra.com"
	} else {
		domain = "http://testmanager.wb-intra.com"
	}
}

func GetDomain() string {
	return domain
}
