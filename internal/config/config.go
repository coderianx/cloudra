package config

var ServerURL = "http://localhost:8080"

func SetServerURL(url string) {
	ServerURL = url
}
