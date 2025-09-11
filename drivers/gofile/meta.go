/*
Package gofile
Author: Da3zKi7<da3zki7@duck.com>

Modifications by ILoveScratch2<ilovescratch@foxmail.com>
*/

package gofile

import (
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
)

type Addition struct {
	driver.RootID
	APIToken string `json:"api_token" required:"true" help:"Get your API token from your Gofile profile page"`
}

var config = driver.Config{
	Name:        "Gofile",
	DefaultRoot: "",
	LocalSort:   false,
	OnlyProxy:   false,
	NoCache:     false,
	NoUpload:    false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Gofile{}
	})
}