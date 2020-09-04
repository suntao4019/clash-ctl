package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Dreamacro/clash-ctl/common"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
)

func HandleServerCommand(args []string) {
	if len(args) == 0 {
		return
	}

	cfg, err := common.ReadCfg()
	if err != nil {
		return
	}

	switch args[0] {
	case "ls":
		servers := cfg.Servers
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Address", "Port", "Secret", "HTTPS"})

		rows := []table.Row{}
		for name, s := range servers {
			rows = append(rows, []interface{}{name, s.Host, s.Port, s.Secret, s.HTTPS})
		}

		t.AppendRows(rows)
		t.Render()
	case "add":
		form := []common.Field{
			{
				Name: "name",
				Prompt: promptui.Prompt{
					Label: "server name",
					Validate: func(in string) error {
						if len(in) == 0 {
							return errors.New("name is required")
						} else if _, ok := cfg.Servers[in]; ok {
							return errors.New("name is exist")
						}
						return nil
					},
				},
			},
			{
				Name: "host",
				Prompt: promptui.Prompt{
					Label: "server address",
					Validate: func(in string) error {
						if len(in) == 0 {
							return errors.New("address is required")
						}
						return nil
					},
				},
			},
			{
				Name: "port",
				Prompt: promptui.Prompt{
					Label: "server port",
					Validate: func(in string) error {
						_, err := strconv.Atoi(in)
						if err != nil {
							return errors.New("port must be int")
						}

						return nil
					},
				},
			},
			{
				Name: "secret",
				Prompt: promptui.Prompt{
					Label:    "server secret",
					Validate: func(in string) error { return nil },
				},
			},
			{
				Name: "https",
				Prompt: promptui.Prompt{
					Label: "API is HTTPS?[y/N]",
					Validate: func(in string) error {
						in = strings.ToLower(in)
						if in != "y" && in != "n" && in != "" {
							return errors.New("value must be y, n or empty(n)")
						}
						return nil
					},
				},
			},
		}

		ret, err := common.ReadMap(form)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		cfg.Servers[ret["name"]] = common.Server{
			Host:   ret["host"],
			Port:   ret["port"],
			Secret: ret["secret"],
			HTTPS:  strings.ToLower(ret["https"]) == "y",
		}

		if err := common.SaveCfg(cfg); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("write server success")
	}
}
