package config

import (
  "os"
  "fmt"
  "strconv"
  "strings"
  "path/filepath"
  "errors"
  "gopkg.in/yaml.v2"
  "github.com/gobuffalo/packr/v2"
)

type ConfigData map[string]interface{}

type Config struct {
  Data ConfigData
}

func (c *Config) Get(key string) (interface{}, error) {
  v, found := c.Data[key]

  if found {
    return v, nil
  } else {
    return v, errors.New(fmt.Sprintf("Could not read key: %s", key))
  }
}

func (c *Config) GetP(key string) interface{} {
  v, err := c.Get(key)

  if err != nil {
    panic(err)
  }

  return v
}

func (c *Config) GetString(key string) (string, error) {
  value, e := c.Get(key)

  if value == nil {
    value = ""
  }

  return fmt.Sprintf("%v", value), e
}

func (c *Config) GetStringP(key string) string {
  v, err := c.GetString(key)

  if err != nil {
    panic(err)
  }

  return v
}

func (c *Config) GetInt(key string) (int, error) {
  value, e := c.GetString(key)

  if e != nil {
    return 0, e
  }

  return strconv.Atoi(value)
}

func (c *Config) GetIntP(key string) int {
  value, e := c.GetInt(key)

  if e != nil {
    panic(e)
  }

  return value
}

func (c *Config) GetFloat(key string) (float64, error) {
  value, e := c.GetString(key)

  if e != nil {
    return 0, e
  }

  return strconv.ParseFloat(value, 64)
}

func (c *Config) GetFloatP(key string) float64 {
  value, e := c.GetFloat(key)

  if e != nil {
    panic(e)
  }

  return value
}

func (this *Config) Merge(that *Config) *Config {
  data := ConfigData{}

  for k, v := range this.Data {
    data[k] = v
  }

  for k, v := range that.Data {
    data[k] = v
  }

  return &Config{Data: data}
}

func (c *Config) MergeWithEnvVars() *Config {
  data := ConfigData{}

  for k, v := range c.Data {
    data[k] = v
  }

  for _, envVar := range os.Environ() {
    pair := strings.Split(envVar, "=")

    k := strings.ToLower(pair[0])
    v := pair[1]

    data[k] = v
  }

  return &Config{Data: data}
}

func Load(path string) (*Config, error) {
  configData, err := loadConfigData(path)

  if configData != nil {
    return &Config{Data: *configData}, err
  } else {
    return nil, err
  }
}

func LoadP(path string) *Config {
  c, err := Load(path)

  if err != nil {
    panic(err)
  }

  return c
}

func LoadSection(path string, section string) (*Config, error) {
  configData, err := loadConfigDataWithSubSection(path, section)

  if configData != nil {
    return &Config{Data: *configData}, err
  } else {
    return nil, err
  }
}

func LoadSectionP(path string, section string) *Config {
  c, err := LoadSection(path, section)

  if err != nil {
    panic(err)
  }

  return c
}


// Path is split into two to prevent creating boxes with unnecessary files
//  for example packs.New("Whatever", "./") would compile all files in the project
//  and include it in the binary
func loadConfigData(path string) (*ConfigData, error) {
  pathToDir := filepath.Dir(path)
  fileName := filepath.Base(path)

  box := packr.New(fmt.Sprintln("Config - %s", pathToDir), pathToDir)

  configInYaml, err := box.FindString(fileName)

  if err != nil {
    return nil, err
  }

  configData := ConfigData{}
  err = yaml.Unmarshal([]byte(configInYaml), &configData)

  if err != nil {
    return nil, err
  }

  return &configData, nil
}

func loadConfigDataWithSubSection(path string, subSection string) (*ConfigData, error) {
  configData, err := loadConfigData(path)

  if err == nil {
    return configData.SubSection(subSection)
  } else {
    return nil, err
  }
}

func (c *ConfigData) SubSection(name string) (*ConfigData, error) {
  result := make(ConfigData)

  subSection, ok := (*c)[name]

  if ok == false {
    return nil, errors.New(fmt.Sprintf("Could not read sub-section: %s", name))
  }

  if subSection == nil {
    subSection = make(map[interface{}]interface{})
  }

  for k, v := range subSection.(map[interface{}]interface{}) {
    result[k.(string)] = v
  }

  return &result, nil
}

