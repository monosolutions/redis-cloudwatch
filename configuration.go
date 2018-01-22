package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"reflect"
	"regexp"

	"github.com/urfave/cli"
)

type Configuration struct {
	Password  string `json:"password"`
	Address   string `json:"address"`
	Port      string `json:"port"`
	SecretKey string `json:"secretkey"`
	AccessKey string `json:"accesskey"`
}
type ConfigurationService struct{}

func (config *ConfigurationService) GetConfigFile() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%v/.redis-cloudwatch", usr.HomeDir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0644)
	}
	return fmt.Sprintf("%v/config.json", path)
}

func set(configuration *Configuration, field string, value string) {
	v := reflect.ValueOf(configuration).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
}
func (config *ConfigurationService) Set(field string, value string) error {
	configuration := config.Read()
	conf := &configuration
	set(conf, field, value)
	err := config.Write(conf)

	return err

}
func (config *ConfigurationService) Write(configuration *Configuration) error {
	file := config.GetConfigFile()
	data, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(data), 0644)

}
func (config *ConfigurationService) Read() Configuration {
	file := config.GetConfigFile()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	result := Configuration{}
	json.Unmarshal(data, &result)

	return result
}
func (config *ConfigurationService) Run(c *cli.Context) error {

	reader := bufio.NewReader(os.Stdin)
	re := regexp.MustCompile(`\r?\n`)

	fmt.Print("Enter password:")
	password, _ := reader.ReadString('\n')
	password = re.ReplaceAllString(password, "")

	fmt.Print("Enter Address:")
	address, _ := reader.ReadString('\n')
	address = re.ReplaceAllString(address, "")

	fmt.Print("Enter Port:")
	port, _ := reader.ReadString('\n')
	port = re.ReplaceAllString(port, "")

	fmt.Print("Enter Access Key:")
	accesskey, _ := reader.ReadString('\n')
	accesskey = re.ReplaceAllString(accesskey, "")

	fmt.Print("Enter Secret Key:")
	secretkey, _ := reader.ReadString('\n')
	secretkey = re.ReplaceAllString(secretkey, "")

	configuration := &Configuration{
		Password:  password,
		Address:   address,
		Port:      port,
		AccessKey: accesskey,
		SecretKey: secretkey,
	}
	return config.Write(configuration)

}
