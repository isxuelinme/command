package command

import (
	"log"
)

type mainArgs struct {
	BaseArgs
	Command     string `json:"command"`
	ServiceName string `json:"serviceName"`
}
type modelArgs struct {
	BaseArgs
	Path string `json:"path"`
}

func INIT() {
	mainCommand()
}

func mainCommand() {
	Add(&Command{
		Command: "rpc",
		Type:    TypeOfSystem,
		Desc:    "",
		MinArgs: 1,
		ArgsRelations: []string{
			"command",
			"serviceName",
		},
		SubCommands: []*Command{
			_initCommand(),
		},
		//NextCommand: _nextCommand(),
	}, func(params mainArgs, command *Command) {
		log.Printf("in %s command", params.ENV.CurrentCommand, params)
		command.ExecuteSub(params.ENV.NextStepArgs)
		//generate(mainGoDir)
	})
}

func _initCommand() *Command {
	return AddSubCommand(&Command{
		Command: "init",
		Type:    TypeOfSystem,
		SubCommands: []*Command{
			__modelCommand(),
		},
		Desc: "",
	},
		func(params BaseArgs, command *Command) {
			log.Printf("in %s command", params.ENV.CurrentCommand, params)
			command.ExecuteSub(params.ENV.NextStepArgs)
		})
}
func __modelCommand() *Command {
	return AddSubCommand(&Command{
		Command: "model",
		Type:    TypeOfSystem,
		Desc:    "",
		SubCommands: []*Command{
			___modelCommand(),
		},
		ArgsRelations: []string{
			"path",
		},
	},
		func(params modelArgs, command *Command) {
			log.Printf("in %s command", params.ENV.CurrentCommand, params)
			command.ExecuteSub(params.ENV.NextStepArgs)
		})
}
func ___modelCommand() *Command {
	return AddSubCommand(&Command{
		Command: "model_model",
		Type:    TypeOfSystem,
		Desc:    "",
		ArgsRelations: []string{
			"path",
		},
	},
		func(params modelArgs, command *Command) {
			log.Printf("in %s command", params.ENV.CurrentCommand, params)
		})
}

func _nextCommand() *Command {
	return AddNextCommand(&Command{
		Command: "nextCommand",
		Type:    TypeOfSystem,
		Desc:    "",
		MinArgs: 1},
		func(params BaseArgs, command *Command) {
			log.Printf("in %s command", params.ENV.CurrentCommand)
		})
}
