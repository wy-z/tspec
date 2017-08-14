package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wy-z/tspec/tspec"
)

type cliOpts struct {
	PkgPath  string
	TypeExpr string
}

func exit(code int, msg string, args ...interface{}) {
	log.Errorf(msg, args...)
	os.Exit(code)
}

func main() {
	app := cli.NewApp()
	app.Name = "TSpec"
	app.Usage = "Parse golang data structure into json schema."

	opts := new(cliOpts)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "package, p",
			Usage:       "package dir",
			Value:       ".",
			Destination: &opts.PkgPath,
		},
		cli.StringFlag{
			Name:        "expression, expr",
			Usage:       "type expression",
			Destination: &opts.TypeExpr,
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		if c.NArg() > 0 {
			opts.TypeExpr = c.Args().Get(0)
		}

		wd, err := os.Getwd()
		if err != nil {
			exit(1, "failed to get working dir: %s", err)
		}
		importPkg, err := build.Import(opts.PkgPath, wd, build.ImportComment)
		if err != nil {
			msg := fmt.Sprintf("failed to import pkg %s: %s", opts.PkgPath, err)
			err = cli.NewExitError(msg, 1)
			return
		}
		parser := tspec.NewParser()
		pkg, err := parser.ParseDir(opts.PkgPath, importPkg.Name)
		if err != nil {
			msg := fmt.Sprintf("failed to parse pkg %s: %s", opts.PkgPath, err)
			err = cli.NewExitError(msg, 1)
			return
		}
		_, err = parser.Parse(pkg, opts.TypeExpr)
		if err != nil {
			msg := fmt.Sprintf("failed to parse type expr %s: %s", opts.TypeExpr, err)
			err = cli.NewExitError(msg, 1)
			return
		}
		defs := parser.Definitions()
		bytes, err := json.MarshalIndent(defs, "", "\t")
		if err != nil {
			msg := fmt.Sprintf("failed to marshal definitions: %s", err)
			err = cli.NewExitError(msg, 1)
			return
		}
		fmt.Println(string(bytes))
		return
	}

	app.Run(os.Args)
}
