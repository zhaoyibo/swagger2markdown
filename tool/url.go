package tool

var domain string
var project string

func InitDomain(env string) {
	if env == "dev" {
		domain = "http://devmanager.wb-intra.com"
	} else {
		domain = "http://testmanager.wb-intra.com"
	}
}

func SetDomain(d string) {
	domain = d
}

func SetProject(p string) {
	project = p
}

func GetDomain() string {
	return domain
}

func GetProject() string {
	return project
}
