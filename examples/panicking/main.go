package main

import (
  "os"
  "fmt"
  "app/config"
  "strings"
)

func main() {
  envVar := "HERO_NAME"
  os.Setenv(envVar, "Jon")

  mainConfig := config.LoadP("examples/config.yaml")

  width := mainConfig.GetP("width")
  height := mainConfig.GetP("heigth")
  fmt.Println(fmt.Sprintf("Width: %v", width))
  fmt.Println(fmt.Sprintf("Height: %v", height))

  width = mainConfig.GetStringP("width")
  height = mainConfig.GetStringP("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  width = mainConfig.GetIntP("width")
  height = mainConfig.GetFloatP("heigth")
  fmt.Println(fmt.Sprintf("Width: %d", width))
  fmt.Println(fmt.Sprintf("Height: %f", height))

  isAwesome := mainConfig.GetBoolP("is_awesome")
  fmt.Println(fmt.Sprintf("IsAwesome: %t", isAwesome))

  overrides := config.LoadSectionP("examples/overrides.yaml", "env_vars")

  height = overrides.GetStringP("heigth")
  fmt.Println(fmt.Sprintf("Height: %s", height))

  mainWithOverrides := mainConfig.Merge(overrides)

  width = mainWithOverrides.GetStringP("width")
  height = mainWithOverrides.GetStringP("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  isAwesome = mainWithOverrides.GetBoolP("is_awesome")
  fmt.Println(fmt.Sprintf("IsAwesome: %t", isAwesome))


  fullConfig := mainWithOverrides.MergeWithEnvVars()

  width = fullConfig.GetStringP("width")
  height = fullConfig.GetStringP("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  envVarLowerCased := strings.ToLower(envVar)
  heroName := fullConfig.GetStringP(envVarLowerCased)
  fmt.Println(fmt.Sprintf("%s: %s", envVarLowerCased, heroName))
}
