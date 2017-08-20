package proxy

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Service struct {
	Name string
	Url  string
}

type Config struct {
	UrlPattern string
	Services   []Service
}

func (c *Config) GetService(host string) *Service {
	isDev := regexp.MustCompile(`.*\.dev$`)
	if !isDev.MatchString(host) {
		return nil
	}

	name := strings.Split(host, ".")[0]
	for _, s := range c.Services {
		if s.Name == name {
			return &s
		}
	}

	return nil
}

func LoadConfig(filepath string) *Config {
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

		config := strings.Split(line, " ")
		name, url := config[0], config[1]
		services = append(services, Service{Name: name, Url: url})
	}

	return &Config{
		UrlPattern: `.*\.dev$`,
		Services:   services,
	}

}
