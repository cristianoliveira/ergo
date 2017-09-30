package proxy

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Service struct {
	Name string
	URL  string
}

type Config struct {
	Port       string
	Domain     string
	URLPattern string
	Verbose    bool
	Services   []Service
}

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

func NewConfig() *Config {
	return &Config{
		Port:       "2000",
		Domain:     ".dev",
		URLPattern: `.*\.dev$`,
		Services:   nil,
	}
}

func LoadServices(filepath string) []Service {
	file, e := os.Open(filepath)
	defer file.Close()
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	services := []Service{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		declaration := regexp.MustCompile(`(\S+)`)
		config := declaration.FindAllString(line, -1)
		if config != nil {
			name, url := config[0], config[1]
			services = append(services, Service{Name: name, URL: url})
		}
	}

	return services
}
