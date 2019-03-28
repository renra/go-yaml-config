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
var primaryHeight float64 = 200.5
var primaryLength int = 200
var numbers [3]string = [3]string{"one", "two", "three"}
var primaryIsAwesome bool = false
var primaryIsTerrible bool = true

var secondaryFileName string = "./secondaryFile.yaml"
var secondaryWidth int = 400
var secondaryHeight int = 400

var tertiaryHeight int = 600
var heroName string = "Jon"
var section string = "env_vars"
var tertiaryIsAwesome bool = true
var tertiaryIsTerrible bool = false

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
    "WIDTH": primaryWidth,
    "height": primaryHeight,
    "length": primaryLength,
    "numbers": numbers,
    "is_awesome": primaryIsAwesome,
    "is_terrible": primaryIsTerrible,
  })

  writeYaml(secondaryFileName, map[string]interface{}{
    section: map[string]interface{}{
      "width": secondaryWidth,
      "HEIGHT": secondaryHeight,
    },
  })

  os.Setenv("HERO_NAME", heroName)
  os.Setenv("HEIGHT", fmt.Sprintf("%d", tertiaryHeight))

  os.Setenv("IS_AWESOME", fmt.Sprintf("%t", tertiaryIsAwesome))
  os.Setenv("IS_TERRIBLE", fmt.Sprintf("%t", tertiaryIsTerrible))
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

func TestLoad(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%.1f", primaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedLength := fmt.Sprintf("%d", primaryLength)
  lengthFromConfig, err := config.GetString("length")

  if lengthFromConfig != expectedLength {
    t.Errorf("Expected %s, got %s", expectedLength, lengthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: length")
  }

  expectedIsAwesome := fmt.Sprintf("%t", primaryIsAwesome)
  isAwesomeFromConfig, err := config.GetString("is_awesome")

  if isAwesomeFromConfig != expectedIsAwesome {
    t.Errorf("Expected %s, got %s", expectedIsAwesome, isAwesomeFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_awesome")
  }

  expectedIsTerrible := fmt.Sprintf("%t", primaryIsTerrible)
  isTerribleFromConfig, err := config.GetString("is_terrible")

  if isTerribleFromConfig != expectedIsTerrible {
    t.Errorf("Expected %s, got %s", expectedIsTerrible, isTerribleFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_terrible")
  }

  expectedValue := ""
  unexistingValue, err := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }

  if err == nil {
    t.Errorf("Expected not to find key: unexisting")
  }
}

func TestLoadSection(t *testing.T) {
  config, _ := config.LoadSection(fmt.Sprintf("test/%s", secondaryFileName), section)

  expectedWidth := fmt.Sprintf("%d", secondaryWidth)
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%d", secondaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedValue := ""
  unexistingValue, err := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }

  if err == nil {
    t.Errorf("Expected not to find key: unexisting")
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
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%d", secondaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedLength := fmt.Sprintf("%d", primaryLength)
  lengthFromConfig, err := config.GetString("length")

  if lengthFromConfig != expectedLength {
    t.Errorf("Expected %s, got %s", expectedLength, lengthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: length")
  }

  expectedIsAwesome := fmt.Sprintf("%t", primaryIsAwesome)
  isAwesomeFromConfig, err := config.GetString("is_awesome")

  if isAwesomeFromConfig != expectedIsAwesome {
    t.Errorf("Expected %s, got %s", expectedIsAwesome, isAwesomeFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_awesome")
  }

  expectedIsTerrible := fmt.Sprintf("%t", primaryIsTerrible)
  isTerribleFromConfig, err := config.GetString("is_terrible")

  if isTerribleFromConfig != expectedIsTerrible {
    t.Errorf("Expected %s, got %s", expectedIsTerrible, isTerribleFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_terrible")
  }

  expectedValue := ""
  unexistingValue, err := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }

  if err == nil {
    t.Errorf("Expected not to find key: unexisting")
  }
}

func TestMergeWithEnvVars(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config = config.MergeWithEnvVars()

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%d", tertiaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedHeroName := heroName
  heroNameFromConfig, err := config.GetString("hero_name")

  if heroNameFromConfig != expectedHeroName {
    t.Errorf("Expected %s, got %s", expectedHeroName, heroNameFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: hero_name")
  }

  expectedIsAwesome := fmt.Sprintf("%t", tertiaryIsAwesome)
  isAwesomeFromConfig, err := config.GetString("is_awesome")

  if isAwesomeFromConfig != expectedIsAwesome {
    t.Errorf("Expected %s, got %s", expectedIsAwesome, isAwesomeFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_awesome")
  }

  expectedIsTerrible := fmt.Sprintf("%t", tertiaryIsTerrible)
  isTerribleFromConfig, err := config.GetString("is_terrible")

  if isTerribleFromConfig != expectedIsTerrible {
    t.Errorf("Expected %s, got %s", expectedIsTerrible, isTerribleFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_terrible")
  }
}

func TestGet(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedValue := numbers
  valueFromConfig, err := config.Get("numbers")

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

  if err != nil {
    t.Errorf("Expected to find key: numbers")
  }
}

func TestLoadP(t *testing.T) {
  config := config.LoadP(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%.1f", primaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedLength := fmt.Sprintf("%d", primaryLength)
  lengthFromConfig, err := config.GetString("length")

  if lengthFromConfig != expectedLength {
    t.Errorf("Expected %s, got %s", expectedLength, lengthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: length")
  }

  expectedIsAwesome := fmt.Sprintf("%t", primaryIsAwesome)
  isAwesomeFromConfig, err := config.GetString("is_awesome")

  if isAwesomeFromConfig != expectedIsAwesome {
    t.Errorf("Expected %s, got %s", expectedIsAwesome, isAwesomeFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_awesome")
  }

  expectedIsTerrible := fmt.Sprintf("%t", primaryIsTerrible)
  isTerribleFromConfig, err := config.GetString("is_terrible")

  if isTerribleFromConfig != expectedIsTerrible {
    t.Errorf("Expected %s, got %s", expectedIsTerrible, isTerribleFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_terrible")
  }

  expectedValue := ""
  unexistingValue, err := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }

  if err == nil {
    t.Errorf("Expected not to find key: unexisting")
  }
}

func TestLoadSectionP(t *testing.T) {
  config := config.LoadSectionP(fmt.Sprintf("test/%s", secondaryFileName), section)

  expectedWidth := fmt.Sprintf("%d", secondaryWidth)
  widthFromConfig, err := config.GetString("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight := fmt.Sprintf("%d", secondaryHeight)
  heightFromConfig, err := config.GetString("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %s, got %s", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: height")
  }

  expectedValue := ""
  unexistingValue, err := config.GetString("unexisting")

  if unexistingValue != expectedValue {
    t.Errorf("Expected %s, got %s", expectedValue, unexistingValue)
  }

  if err == nil {
    t.Errorf("Expected not to find key: unexisting")
  }
}

func TestLoadPUnexistingFile(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config.LoadP("whatever.yaml")
}

func TestLoadSectionPUnexistingFile(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config.LoadSectionP("whatever.yaml", "whatever")
}

func TestLoadSectionPUnexistingSection(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config.LoadSectionP(secondaryFileName, "whatever")
}

func TestGetP(t *testing.T) {
  config := config.LoadP(fmt.Sprintf("test/%s", mainFileName))

  expectedValue := numbers
  valueFromConfig := config.GetP("numbers")

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

func TestGetPUnexistingKey(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config := config.LoadP(fmt.Sprintf("test/%s", mainFileName))
  config.GetP("whatever")
}

func TestGetStringP(t *testing.T) {
  config := config.LoadP(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := fmt.Sprintf("%d", primaryWidth)
  widthFromConfig := config.GetStringP("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %s, got %s", expectedWidth, widthFromConfig)
  }
}

func TestGetStringPUnexistingKey(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config := config.LoadP(fmt.Sprintf("test/%s", mainFileName))
  config.GetStringP("whatever")
}

func TestGetInt(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := primaryWidth
  widthFromConfig, err := config.GetInt("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %d, got %d", expectedWidth, widthFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedWidth = 0
  widthFromConfig, err = config.GetInt("unexisting_width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %d, got %d", expectedWidth, widthFromConfig)
  }

  if err == nil {
    t.Errorf("Expected to see error here")
  }
}

func TestGetIntP(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedWidth := primaryWidth
  widthFromConfig := config.GetIntP("width")

  if widthFromConfig != expectedWidth {
    t.Errorf("Expected %d, got %d", expectedWidth, widthFromConfig)
  }
}

func TestGetIntPUnexistingKey(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config.GetIntP("unexisting_width")
}

func TestGetFloat(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedHeight := primaryHeight
  heightFromConfig, err := config.GetFloat("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %f, got %f", expectedHeight, heightFromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: width")
  }

  expectedHeight = 0
  heightFromConfig, err = config.GetFloat("unexisting_height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %f, got %f", expectedHeight, heightFromConfig)
  }

  if err == nil {
    t.Errorf("Expected to see error here")
  }
}

func TestGetFloatP(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))

  expectedHeight := primaryHeight
  heightFromConfig := config.GetFloatP("height")

  if heightFromConfig != expectedHeight {
    t.Errorf("Expected %f, got %f", expectedHeight, heightFromConfig)
  }
}

func TestGetFloatPUnexistingKey(t *testing.T) {
  defer func(){
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()

  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config.GetFloatP("unexisting_height")
}

func TestGetBool(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config = config.MergeWithEnvVars()

  expected := tertiaryIsAwesome
  fromConfig, err := config.GetBool("is_awesome")

  if fromConfig != expected {
    t.Errorf("Expected %t, got %t", expected, fromConfig)
  }

  if err != nil {
    t.Errorf("Expected to find key: is_awesome")
  }

  expected = false
  fromConfig, err = config.GetBool("unexisting_flag")

  if fromConfig != expected {
    t.Errorf("Expected %t, got %t", expected, fromConfig)
  }

  if err == nil {
    t.Errorf("Expected to see error here")
  }
}

func TestGetBoolP(t *testing.T) {
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config = config.MergeWithEnvVars()

  expected := tertiaryIsTerrible
  fromConfig := config.GetBoolP("is_terrible")

  if fromConfig != expected {
    t.Errorf("Expected %t, got %t", expected, fromConfig)
  }
}

func TestGetBoolPUnexistingKey(t *testing.T) {
  defer func() {
    r := recover()

    if r == nil {
      t.Errorf("Expected to be recovering from a panic here")
    }
  }()
  config, _ := config.Load(fmt.Sprintf("test/%s", mainFileName))
  config = config.MergeWithEnvVars()

  config.GetBoolP("is_whatever")
}
