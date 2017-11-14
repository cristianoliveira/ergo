package proxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

const pollIntervall = 500

//Service holds the details of the service (Name and URL)
type Service struct {
	Name string
	URL  string
}

//Config holds the configuration for the proxy.
type Config struct {
	mutex      *sync.Mutex
	lastChange time.Time
	size       int64

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

	for _, s := range c.Services {
		if len(parts) > 0 && s.Name+c.Domain == parts[0] {
			return &s
		}
	}

	return nil
}

//GetProxyPacURL returns the correct url for the pac file
func (c *Config) GetProxyPacURL() string {
	return "http://127.0.0.1:" + c.Port + "/proxy.pac"
}

//NewConfig gets the new config.
func NewConfig() Config {
	return Config{
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
func (c *Config) LoadServices() error {
	services, err := LoadServicesFromConfig(c.ConfigFile)
	if err != nil {
		return err
	}

	c.getMutex().Lock()
	defer c.getMutex().Unlock()
	c.Services = services

	return nil
}

func (c *Config) getMutex() *sync.Mutex {
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}

	return c.mutex
}

// ListenServices updates the services for each message in a given channel
func (c *Config) ListenServices(servicesSignal <-chan []Service) {
	for {
		select {
		case services := <-servicesSignal:
			c.getMutex().Lock()
			c.Services = services
			c.mutex.Unlock()
		}
	}
}

// WatchConfigFile listen for file changes and sends signal for updates
func (c *Config) WatchConfigFile(servicesChan chan []Service) {
	ticker := time.NewTicker(pollIntervall * time.Millisecond)
	quit = make(chan struct{})

	for {
		select {
		case <-ticker.C:
			info, err := os.Stat(c.ConfigFile)
			if err != nil {
				log.Printf("Error reading config file: %s\r\n", err.Error())
				continue
			}

			if info.ModTime().Before(c.lastChange) || info.Size() != c.size {
				services, err := LoadServicesFromConfig(c.ConfigFile)
				if err != nil {
					log.Printf("Error reading the modified config file: %s\r\n", err.Error())
					continue
				}

				c.size = info.Size()
				c.lastChange = info.ModTime()
				servicesChan <- services
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// LoadServicesFromConfig reads the given path and parse it into services
func LoadServicesFromConfig(filepath string) ([]Service, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("file error: %v", err)
	}
	defer file.Close()

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
		return fmt.Errorf("File error: %v", e)
	}

	defer file.Close()

	serviceStr := service.Name + " " + service.URL + "\n"
	_, err := file.WriteString(serviceStr)
	return err
}

// RemoveService removes a service from the filepath
func RemoveService(filepath string, service Service) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("File error: %v\n", err)
		return err
	}

	serviceRegex := regexp.MustCompile(service.Name + "\\s+" + service.URL + "\n")

	file = serviceRegex.ReplaceAll(file, []byte("\n"))

	ioutil.WriteFile(filepath, file, 0755)

	return nil
}
