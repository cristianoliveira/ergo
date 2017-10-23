package proxy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	modTime time.Time
	size    int64
)

var configChan = make(chan []Service, 1)

//Service holds the details of the service (Name and URL)
type Service struct {
	Name string
	URL  string
}

//Config holds the configuration for the proxy.
type Config struct {
	Port       string
	Domain     string
	URLPattern string
	Verbose    bool
	Services   []Service
	ConfigFile string
}

//GetService gets the service for the given host.
func (c *Config) GetService(host string) *Service {
	domainPattern := regexp.MustCompile(`(\w*\:\/\/)?(.+)` + c.Domain)
	parts := domainPattern.FindAllString(host, -1)
	//we must lock the access as the configuration can be dynamically loaded
	select {
	case srv := <-configChan:
		c.Services = srv
	default:
	}
	for _, s := range c.Services {
		if len(parts) > 0 && s.Name+c.Domain == parts[0] {
			return &s
		}
	}

	return nil
}

//NewConfig gets the new config.
func NewConfig() *Config {
	return &Config{
		Port:       "2000",
		Domain:     ".dev",
		URLPattern: `.*\.dev$`,
		Services:   nil,
	}
}

//NewService gets the new service.
func NewService(name, url string) Service {
	return Service{
		Name: name,
		URL:  url,
	}
}

//LoadServices loads the services from filepath, returns an error
//if the configuration could not be parsed
func LoadServices(filepath string) ([]Service, error) {

	info, err := os.Stat(filepath)

	if err != nil {
		return nil, err
	}

	size = info.Size()
	modTime = info.ModTime()

	file, e := os.Open(filepath)

	if e != nil {
		return nil, fmt.Errorf("file error: %v", e)
	}

	defer file.Close()

	log.Println("Just received ", filepath)

	services := []Service{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		declaration := regexp.MustCompile(`(\S+)`)
		config := declaration.FindAllString(line, -1)
		if config == nil {
			continue
		}
		if len(config) != 2 {
			return nil, fmt.Errorf("file error: invalid format `%v` expected `{NAME} {URL}`", line)
		}
		name, url := config[0], config[1]
		services = append(services, Service{Name: name, URL: url})
	}

	return services, nil
}

//AddService adds new service to the filepath
func AddService(filepath string, service Service) error {
	file, e := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if e != nil {
		log.Printf("File error: %v\n", e)
		return e
	}

	defer file.Close()

	serviceStr := service.Name + " " + service.URL + "\n"
	_, err := file.WriteString(serviceStr)
	return err
}
