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

  mainConfig, _ := config.Load("examples/config.yaml")

  fmt.Println(fmt.Sprintf("Width: %v", mainConfig.Get("width")))
  fmt.Println(fmt.Sprintf("Height: %v", mainConfig.Get("height")))

  fmt.Println(fmt.Sprintf("Width: %s", mainConfig.GetString("width")))
  fmt.Println(fmt.Sprintf("Height: %s", mainConfig.GetString("height")))

  overrides, _ := config.LoadSection("examples/overrides.yaml", "env_vars")

  fmt.Println(fmt.Sprintf("Width: %s", overrides.GetString("width")))
  fmt.Println(fmt.Sprintf("Height: %s", overrides.GetString("height")))

  mainWithOverrides := mainConfig.Merge(overrides)

  fmt.Println(fmt.Sprintf("Width: %s", mainWithOverrides.GetString("width")))
  fmt.Println(fmt.Sprintf("Height: %s", mainWithOverrides.GetString("height")))

  fullConfig := mainWithOverrides.MergeWithEnvVars()

  fmt.Println(fmt.Sprintf("Width: %s", fullConfig.GetString("width")))
  fmt.Println(fmt.Sprintf("Height: %s", fullConfig.GetString("height")))

  envVarLowerCased := strings.ToLower(envVar)
  fmt.Println(fmt.Sprintf("%s: %s", envVarLowerCased, envVarLowerCased))

  _, err := config.Load("non_existing.yaml")

  if err != nil {
    fmt.Println(err)
  }

  _, err = config.LoadSection("non_existing.yaml", "something_not_there")

  if err != nil {
    fmt.Println(err)
  }

  _, err = config.LoadSection("examples/overrides.yaml", "something_not_there")

  if err != nil {
    fmt.Println(err)
  }
}
