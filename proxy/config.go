package proxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Service holds the details of the service (Name and URL)
type Service struct {
	Name string
	URL  string
}

// Empty service means no name or no url
func (s Service) Empty() bool {
	return s.Name == "" || s.URL == ""
}

// Config holds the configuration for the proxy.
type Config struct {
	mutex      sync.Mutex
	lastChange time.Time
	size       int64

	Port       string
	Domain     string
	Verbose    bool
	Services   map[string]Service
	ConfigFile string
}

// Defines the name of ergo env variable for configuration.
const (
	PortEnv       = "ERGO_PORT"
	DomainEnv     = "ERGO_DOMAIN"
	VerboseEnv    = "ERGO_VERBOSE"
	ConfigFileEnv = "ERGO_CONFIG_FILE"

	PortDefault           = "2000"
	DomainDefault         = ".dev"
	ConfigFilePathDefault = "./.ergo"
)

// NewConfig gets the defaults config.
func NewConfig() *Config {
	var config = &Config{
		Port:       PortDefault,
		Domain:     DomainDefault,
		ConfigFile: ConfigFilePathDefault,
		Verbose:    os.Getenv(VerboseEnv) != "",
		Services:   make(map[string]Service),
	}

	port, isPortPresent := os.LookupEnv(PortEnv)
	if isPortPresent {
		config.Port = port
	}

	domain, isDomainPresent := os.LookupEnv(DomainEnv)
	if isDomainPresent {
		config.Domain = domain
	}

	configFile, isConfigFilePresent := os.LookupEnv(ConfigFileEnv)
	if isConfigFilePresent {
		config.ConfigFile = configFile
	}

	return config
}

// OverrideBy makes sure that it sets the correct config based on
// the defaults and the passed by argument
func (c *Config) OverrideBy(new *Config) {
	if new.Port != "" {
		c.Port = new.Port
	}

	if new.Domain != "" {
		c.Domain = new.Domain
	}

	if new.Verbose != c.Verbose {
		c.Verbose = new.Verbose
	}

	if new.ConfigFile != "" {
		c.ConfigFile = new.ConfigFile
	}
}

var once sync.Once
var domainPattern *regexp.Regexp

// GetService gets the service for the given host.
func (c *Config) GetService(host string) Service {
	once.Do(func() {
		domainPattern = regexp.MustCompile(`((\w*\:\/\/)?.+)(` + c.Domain + `)`)
	})

	parts := domainPattern.FindAllStringSubmatch(host, -1)

	if len(parts) < 1 {
		return Service{}
	}

	return c.Services[parts[0][1]]
}

// GetProxyPacURL returns the correct url for the pac file
func (c *Config) GetProxyPacURL() string {
	return "http://127.0.0.1:" + c.Port + "/proxy.pac"
}

// AddService add a service using the correct key
func (c *Config) AddService(service Service) error {
	if service.Empty() {
		return fmt.Errorf("Service is invalid")
	}

	c.Services[service.Name] = service
	return nil
}

// LoadServices loads the services from filepath, returns an error
// if the configuration could not be parsed
func (c *Config) LoadServices() error {
	services, err := readServicesFromFile(c.ConfigFile)
	fmt.Println("services", services)
	if err != nil {
		return err
	}

	updatedServices := make(map[string]Service)
	for _, s := range services {
		if !s.Empty() {
			updatedServices[s.Name] = s
		}
	}

	c.mutex.Lock()
	{
		c.Services = updatedServices
	}
	c.mutex.Unlock()

	return nil
}

// WatchConfigFile listen for file changes and updates the config services
func (c *Config) WatchConfigFile(tickerChan <-chan time.Time) {
	for range tickerChan {
		info, err := os.Stat(c.ConfigFile)
		if err != nil {
			log.Printf("Error reading config file: %s\r\n", err.Error())
			continue
		}

		if info.ModTime().Before(c.lastChange) || info.Size() != c.size {
			c.size = info.Size()
			c.lastChange = info.ModTime()

			err = c.LoadServices()
			if err != nil {
				log.Printf("Error reading the modified config file")
				continue
			}
		}
	}
}

// readServicesFromFile reads the given path and parse it into services
func readServicesFromFile(filepath string) ([]Service, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("file error: %v", err)
	}
	defer file.Close()

	services := []Service{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")

		pair := strings.Split(line, " ")

		if len(pair) != 2 {
			return nil, fmt.Errorf("File error: invalid format `%v` expected `{NAME} {URL}`", line)
		}

		urlPattern := regexp.MustCompile(`(\w+\:\/\/)?([^ ]+):(\d+)`)

		var name, urlWithPort string
		first := urlPattern.MatchString(pair[0])
		if first {
			name, urlWithPort = pair[1], pair[0]
		} else {
			name, urlWithPort = pair[0], pair[1]
		}

		if name == "" || urlWithPort == "" {
			return nil, fmt.Errorf("File error: invalid format `%v` expected `{NAME} {URL}`", line)
		}

		services = append(services, Service{Name: name, URL: urlWithPort})
	}

	return services, nil
}

// AddService adds new service to the filepath
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
