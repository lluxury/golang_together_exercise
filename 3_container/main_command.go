package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	//"github.com/xianlubird/mydocker/container"
	"github.com/lluxury/golang_together_exercise/3_container/container"
	"github.com/lluxury/golang_together_exercise/3_container/cgroups/subsystems"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
						mydocker run -ti [command]`,
	Action: func(context *cli.Context) error {
							if len(context.Args()) < 1{
								return fmt.Errorf("Missing container command")
							}
							var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
							//cmd := context.Args().Get(0)
							tty := context.Bool("ti")
							resConf := &subsystems.ResourceConfig{
								MemoryLimit: "m",
							}
							
							Run(tty, cmdArray, resConf)
							return nil
	},
	Flags:  []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
	},
}


var initCommand = cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside",
        Action: func(context *cli.Context) error{
		log.Infof("init come on")
		// cmd := context.Args().Get(0)
		// log.Infof("command %s", cmd)
		// err := container.RunContainerInitProcess(cmd, nil) 逻辑错误
		err := container.RunContainerInitProcess()
		return err
	},
}
