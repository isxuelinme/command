package command

import (
	"encoding/json"
	"fmt"
	"log"
)

type NullArgs struct {
}
type BaseArgs struct {
	ENV struct {
		HostName       string   `json:"host_name"`
		System         string   `json:"system"`
		Arch           string   `json:"arch"`
		GoVersion      string   `json:"go_version"`
		Timezone       string   `json:"timezone"`
		CpuNums        int      `json:"cpu_nums"`
		CpuThreadNum   int      `json:"cpu_thread_num"`
		MemorySize     int      `json:"memory_size"`
		PWD            string   `json:"pwd"`
		OsArgs         []string `json:"os_args"`
		CurrentCommand string   `json:"current_command"`
		NextStepArgs   []string `json:"next_step_args"`
	} `json:"env"`
}

type Type int

const (
	TypeOfSystem Type = iota
	TypeOfExport
	TypeOfFix
)

type commandParams struct {
	jsonBytes []byte
}
type ArgRelation struct {
	Index int
	Name  string
}
type Command struct {
	Command string
	Type    Type
	Desc    string
	MinArgs int
	//ArgsRelations  []ArgRelation
	ArgsRelations  []string
	NextCommand    *Command
	SubCommands    []*Command
	commandExecute func(params commandParams, command *Command)
	env            map[string]interface{}
}

func (c *Command) Next(args BaseArgs) {
	if c.NextCommand != nil {
		c.executeNext(args.ENV.NextStepArgs)
	} else {
		log.Println("NextCommand is nil", args)
	}
}
func (c *Command) ExecuteSub(args []string) {
	c.executeSub(args)
}

func (c *Command) executeNext(args []string) {
	commandInstance := c.NextCommand

	if len(args) < commandInstance.MinArgs {
		//display help
		log.Println("at least "+fmt.Sprintf("%d", commandInstance.MinArgs)+" params", commandInstance)
		return
	}
	var argsRelation = make(map[string]interface{})
	c.env["current_command"] = args[0]
	c.env["next_step_args"] = args
	argsRelation["env"] = c.env
	for index, arg := range args[1:] {
		for paramIndex, paramName := range commandInstance.ArgsRelations {
			if paramIndex == index {
				argsRelation[paramName] = arg
			}
		}
	}
	argParamsBytes, _ := json.Marshal(argsRelation)
	argParams := commandParams{
		jsonBytes: argParamsBytes,
	}
	commandInstance.env = c.env
	commandInstance.commandExecute(argParams, commandInstance)
	if commandInstance.NextCommand != nil {
		var baseArgs = BaseArgs{}
		json.Unmarshal(argParamsBytes, &baseArgs)
		baseArgs.ENV.NextStepArgs = baseArgs.ENV.NextStepArgs[1:]
		commandInstance.Next(baseArgs)
	}

}

func (c *Command) executeSub(args []string) {
	var command = args[0]
	var commandInstance *Command
	for i, child := range c.SubCommands {
		if child.Command == command {
			commandInstance = c.SubCommands[i]
			break

		}
	}
	if commandInstance != nil {
		if len(args)-1 < commandInstance.MinArgs {
			//display help
			log.Println("at least "+fmt.Sprintf("%d", commandInstance.MinArgs)+" params", commandInstance)
			return
		}
		var argsRelation = make(map[string]interface{})

		c.env["next_step_args"] = args[1:]
		c.env["current_command"] = args[0]
		argsRelation["env"] = c.env

		for index, arg := range args[1:] {
			for paramIndex, paramName := range commandInstance.ArgsRelations {
				if paramIndex == index {
					argsRelation[paramName] = arg
				}
			}
		}

		argParamsBytes, _ := json.Marshal(argsRelation)
		argParams := commandParams{
			jsonBytes: argParamsBytes,
		}
		//json.Unmarshal(argParams, argParams)
		//commandInstance. = argParams
		commandInstance.env = c.env
		commandInstance.commandExecute(argParams, commandInstance)
		if commandInstance.NextCommand != nil {
			var baseArgs = BaseArgs{}
			json.Unmarshal(argParamsBytes, &baseArgs)
			baseArgs.ENV.NextStepArgs = baseArgs.ENV.NextStepArgs[1:]
			commandInstance.Next(baseArgs)
		}
	} else {
		log.Printf("cant found the command %s\n", command)
	}

}
