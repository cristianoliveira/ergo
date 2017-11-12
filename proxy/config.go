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
	mutex      *sync.Mutex
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

//LoadServices loads the services from filepath, returns an error
//if the configuration could not be parsed
func (c *Config) LoadServices() error {
	c.getMutex().Lock()
	defer c.getMutex().Unlock()

	var err error
	c.Services, err = LoadServicesFromConfig(c.ConfigFile)

	return err
}

func (c *Config) getMutex() *sync.Mutex {
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}

	return c.mutex
}

//Sync updates the services for each message in a given channel
func (c *Config) Sync(servicesSignal <-chan []Service) {
	for {
		select {
		case services := <-servicesSignal:
			fmt.Println("signal", services)
			c.getMutex().Lock()
			c.Services = services
			c.mutex.Unlock()
		}
	}
}

// LoadServicesFromConfig reads the given path and parse it into services
func LoadServicesFromConfig(filepath string) ([]Service, error) {

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

			if info.ModTime().Before(modTime) || info.Size() != size {
				services, err := LoadServicesFromConfig(c.ConfigFile)
				if err != nil {
					log.Printf("Error reading the modified config file: %s\r\n", err.Error())
					continue
				}
				servicesChan <- services
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}
