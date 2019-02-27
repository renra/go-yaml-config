package main

import (
  "os"
  "fmt"
  "testing"
  "io/ioutil"
  "app/config"
  "gopkg.in/yaml.v2"
)

var mainFileName string = "./mainFile.yaml"
var primaryWidth int = 200
var primaryHeight int = 200
var primaryLength int = 200
var numbers [3]string = [3]string{"one", "two", "three"}

var secondaryFileName string = "./secondaryFile.yaml"
var secondaryWidth int = 400
var secondaryHeight int = 400

var tertiaryHeight int = 600
var heroName string = "Jon"
var section string = "env_vars"

func writeYaml(path string, data map[string]interface{}) {
  contents, err := yaml.Marshal(&data)

  if err != nil {
    panic(err)
  }

  err = ioutil.WriteFile(path, contents, 0644)

  if err != nil {
    panic(err)
  }
}

func setup() {
  writeYaml(mainFileName, map[string]interface{}{
    "width": primaryWidth,
    "height": primaryHeight,
    "length": primaryLength,
    "numbers": numbers,
  })

  writeYaml(secondaryFileName, map[string]interface{}{
    section: map[string]interface{}{
      "width": secondaryWidth,
      "height": secondaryHeight,
    },
  })

  os.Setenv("HERO_NAME", heroName)
  os.Setenv("HEIGHT", fmt.Sprintf("%d", tertiaryHeight))
}

func teardown() {
  os.Remove(mainFileName)
  os.Remove(secondaryFileName)
}

func TestMain(m *testing.M) {
  setup()

  code := m.Run()

  teardown()
  os.Exit(code)
}

func TestLoadFile(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  expectedHeight := fmt.Sprintf("%d", primaryHeight)
  heightFromConfig := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  expectedLength := fmt.Sprintf("%d", primaryLength)
  lengthFromConfig := config.GetString("length")

  if lengthFromConfig != expectedLength {
    t.Errorf("Expected %s, got %s", expectedLength, lengthFromConfig)
  }

  expectedValue := ""
  unexistingValue := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }
}

func TestLoadSection(t *testing.T) {
  config, _ := config.LoadSection(fmt.Sprintf("test/%s", secondaryFileName), section)

  expectedWidth := fmt.Sprintf("%d", secondaryWidth)
  widthFromConfig := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  expectedHeight := fmt.Sprintf("%d", secondaryHeight)
  heightFromConfig := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  expectedValue := ""
  unexistingValue := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }
}

func TestGet(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedValue := numbers
  valueFromConfig := config.Get("numbers")

  switch valueType := valueFromConfig.(type) {
    case []interface{}:
    default:
      t.Errorf("Expected type []interface{}, got %v", valueType)
  }

  retypedValueFromConfig := valueFromConfig.([]interface{})

  for i, value := range retypedValueFromConfig {
    if value != expectedValue[i] {
      t.Errorf("Expected %v at index %v, got %v", expectedValue[i], i, value)
    }
  }
}

func TestLoadUnexistingFile(t *testing.T) {
  filename := "whatever.yaml"
  config, err := config.Load(filename)

  if err == nil {
    t.Errorf("Expected to get an error after trying to load: %s", filename)
  }

  if config != nil {
    t.Errorf("Expected config to be nil")
  }
}

func TestLoadSectionUnexistingFile(t *testing.T) {
  filename := "whatever.yaml"
  config, err := config.LoadSection(filename, "whatever")

  if err == nil {
    t.Errorf("Expected to get an error after trying to load: %s", filename)
  }

  if config != nil {
    t.Errorf("Expected config to be nil")
  }
}

func TestLoadSectionUnexistingSection(t *testing.T) {
  section := "whatever"
  config, err := config.LoadSection(secondaryFileName, section)

  if err == nil {
    t.Errorf("Expected to get an error after trying to load section: %s", section)
  }

  if config != nil {
    t.Errorf("Expected config to be nil")
  }
}

func TestMerge(t *testing.T) {
  c1, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  c2, _ := config.LoadSection(fmt.Sprintf("test/%s", secondaryFileName), section)

  config := c1.Merge(c2)

  expectedWidth := fmt.Sprintf("%d", secondaryWidth)
  widthFromConfig := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  expectedHeight := fmt.Sprintf("%d", secondaryHeight)
  heightFromConfig := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  expectedLength := fmt.Sprintf("%d", primaryLength)
  lengthFromConfig := config.GetString("length")

  if lengthFromConfig != expectedLength {
    t.Errorf("Expected %s, got %s", expectedLength, lengthFromConfig)
  }

  expectedValue := ""
  unexistingValue := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }
}

func TestMergeWithEnvVars(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config = config.MergeWithEnvVars()

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  expectedHeight := fmt.Sprintf("%d", tertiaryHeight)
  heightFromConfig := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  expectedHeroName := heroName
  heroNameFromConfig := config.GetString("hero_name")

  if heroNameFromConfig != expectedHeroName {
    t.Errorf("Expected %s, got %s", expectedHeroName, heroNameFromConfig)
  }
}
