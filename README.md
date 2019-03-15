# Go Yaml Config

A tool that knows how to read multiple yaml files and override them with environment variables.

## Usage

```go
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

  overrides, _ := config.LoadSection("examples/overrides.yaml", "env_vars")

  width, _ = overrides.GetString("width")
  height, _ = overrides.GetString("heigth")
  fmt.Println(fmt.Sprintf("Width: %s", width))
  fmt.Println(fmt.Sprintf("Height: %s", height))

  width, _ = mainConfig.GetInt("width")
  height, _ = mainConfig.GetFloat("heigth")
  fmt.Println(fmt.Sprintf("Width: %d", width))
  fmt.Println(fmt.Sprintf("Height: %f", height))

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
}
```

If you're tired of handling errors and are sure that you don't want to continue your program's execution after encountering an error, you can use the `P` functions to panic right away. Concretely they are: `LoadP`, `LoadSectionP`, `GetP` and `GetStringP`.

Note that files are loaded via [packr](github.com/gobuffalo/packr/v2) so you will need to use it to compile your program.

### About types

When you use `Get` it returns `interface{}` and you can type-assert it to anything you want. I find it's easiest to use `GetString` though, especially together with `MergeWithEnvVars` because all env vars are strings anyway so it helps to avoid the problem of working with values of different types depending on whether they are overridden or not. You can always use `strconv` on the resulting value.
