# golang-with-mongodb-sample

This repository is demo how to use mongodb with golang

## sample

install mongo driver on golang

```shell
go get go.mongodb.org/mongo-driver/mongo
```

## setup run tools on Taskfile.yml

```yaml
version: '3'

tasks:
  default:
    cmds:
      - echo "This is task cmd"
    silent: true
  
  build:
    cmds:
      - CGO_ENABLED=0 GOOS=linux go build -o bin/main cmd/main.go
    silent: true
  run:
    cmds:
      - ./bin/main
    deps:
      - build
    silent: true

  build-mage:
    cmds:
      - CGO_ENABLED=0 GOOS=linux go build -o ./mage mage-tools/mage.go
    silent: true
  
  build-gg:
    cmds:
      - ./mage -d mage-tools -compile ../gg
    deps:
      - build-mage
    silent: true

  coverage:
    cmds:
      - go test -v -cover ./...
    silent: true
  test:
    cmds:
      - go test -v ./...
    silent: true
```

## setup test script with mage lib

### main execute file

```golang
//go:build ignore
// +build ignore

package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() { os.Exit(mage.Main()) }
```

### setup run script on magefile.go

```golang
//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// clean the build binary
func Clean() error {
	return sh.Rm("bin")
}

// Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "./bin/server", "./cmd/main.go")
}

// start the server
func Launch() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/server")
	if err != nil {
		return err
	}
	return nil
}

// run the test
func Test() error {
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		return err
	}
	return nil
}
```