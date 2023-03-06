package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/joho/godotenv"
	"github.com/phlx-ru/hatchet/cli"
	pkgConfig "github.com/phlx-ru/hatchet/config"
	"github.com/phlx-ru/hatchet/jwt"
	"github.com/phlx-ru/hatchet/logger"

	"storage/internal/conf"
)

var (
	// Name is the name of the compiled software.
	Name = `storage-jwt`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env.local")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(pkgConfig.EnvReplaceDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	logger.SetGlobalDefaultLogger(`warn`) // Silent config load
	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	token := jwt.Make(Name, bc.Auth.Jwt.Secret)

	fmt.Println(cli.ColorGreen + `JWT token generated:` + cli.ColorReset)
	fmt.Println(cli.ColorBlue + token + cli.ColorReset)
	return nil
}
