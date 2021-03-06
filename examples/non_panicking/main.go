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

  width, _ := mainConfig.Get("width")
  height, _ := mainConfig.Get("heigth")
  fmt.Println(fmt.Sprintf("Width: %v", width))
  fmt.Println(fmt.Sprintf("Height: %v", height))

  width, _ = mainConfig.GetString("width")
  height, _ = mainConfig.GetString("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  width, _ = mainConfig.GetInt("width")
  height, _ = mainConfig.GetFloat("heigth")
  fmt.Println(fmt.Sprintf("Width: %d", width))
  fmt.Println(fmt.Sprintf("Height: %f", height))

  isAwesome, _ := mainConfig.GetBool("is_awesome")
  fmt.Println(fmt.Sprintf("IsAwesome: %t", isAwesome))

  overrides, _ := config.LoadSection("examples/overrides.yaml", "env_vars")

  width, _ = overrides.GetString("width")
  height, _ = overrides.GetString("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  isAwesome, _ = overrides.GetBool("is_awesome")
  fmt.Println(fmt.Sprintf("IsAwesome: %t", isAwesome))

  mainWithOverrides := mainConfig.Merge(overrides)

  width, _ = mainWithOverrides.GetString("width")
  height, _ = mainWithOverrides.GetString("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  fullConfig := mainWithOverrides.MergeWithEnvVars()

  width, _ = fullConfig.GetString("width")
  height, _ = fullConfig.GetString("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  envVarLowerCased := strings.ToLower(envVar)
  heroName, _ := fullConfig.GetString(envVarLowerCased)
  fmt.Println(fmt.Sprintf("%s: %s", envVarLowerCased, heroName))

  // Errors when loading files
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

  // Errors when reading keys
  _, err = fullConfig.Get("whatever")

  if err != nil {
    fmt.Println(err)
  }

  _, err = fullConfig.GetString("whatever")

  if err != nil {
    fmt.Println(err)
  }

  _, err = fullConfig.GetInt("whatever")

  if err != nil {
    fmt.Println(err)
  }

  _, err = fullConfig.GetFloat("whatever")

  if err != nil {
    fmt.Println(err)
  }

  _, err = fullConfig.GetBool("whatever")

  if err != nil {
    fmt.Println(err)
  }
}

