package command

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var commands map[string]*Command

func init() {
	commands = make(map[string]*Command)
}

func AddNextCommand[T any](nextCommand *Command, commandExecute func(params T, command *Command)) *Command {
	nextCommand.commandExecute = func(params commandParams, command *Command) {
		var args T
		json.Unmarshal(params.jsonBytes, &args)
		commandExecute(args, command)
	}
	nextCommand.env = make(map[string]interface{})

	return nextCommand
}
func AddSubCommand[T any](childCommand *Command, commandExecute func(params T, command *Command)) *Command {
	childCommand.commandExecute = func(params commandParams, command *Command) {
		var args T
		json.Unmarshal(params.jsonBytes, &args)
		commandExecute(args, command)
	}
	childCommand.env = make(map[string]interface{})
	return childCommand
}

func Add[T any](command *Command, commandExecute func(params T, command *Command)) *Command {
	commands[command.Command] = command
	commands[command.Command].commandExecute = func(params commandParams, command *Command) {
		var args T
		json.Unmarshal(params.jsonBytes, &args)
		commandExecute(args, command)
	}
	commands[command.Command].env = make(map[string]interface{})
	return commands[command.Command]
}

func Execute(args []string) {
	var command string
	if len(args) > 1 {
		command = args[1]
	} else {
		command = "version"
	}

	commandLinked := strings.Join(os.Args[1:], " ")
	log.Printf("Your input args is: %s\n", commandLinked)

	var commandInstance *Command
	for commandName, _ := range commands {
		if commandName == command {
			if _, exists := commands[commandName]; exists {
				commandInstance = commands[commandName]
				break
			}
		}
	}
	if commandInstance != nil {
		if len(args)-1 < commandInstance.MinArgs {
			//display help
			log.Println("at least "+fmt.Sprintf("%d", commandInstance.MinArgs)+" params", commandInstance)
			return
		}
		var argsRelation = make(map[string]interface{})

		commandInstance.env["pwd"], _ = os.Getwd()
		commandInstance.env["os_args"] = os.Args
		commandInstance.env["go_version"] = runtime.Version()
		commandInstance.env["timezone"] = time.Now().Location().String()
		commandInstance.env["cpu_nums"] = runtime.NumCPU()
		commandInstance.env["cpu_thread_num"] = runtime.GOMAXPROCS(-1)
		memInfo, _ := mem.VirtualMemory()
		commandInstance.env["host_name"], _ = os.Hostname()
		commandInstance.env["arch"] = runtime.GOARCH
		commandInstance.env["system"] = runtime.GOOS
		commandInstance.env["memory_size"] = memInfo.Total
		commandInstance.env["next_step_args"] = args[2:]
		commandInstance.env["current_command"] = args[1]

		argsRelation["env"] = commandInstance.env
		for index, arg := range args[2:] {
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
		//commandInstance.Args = argParams
		commandInstance.commandExecute(argParams, commandInstance)
		if commandInstance.NextCommand != nil {
			var baseArgs = BaseArgs{}
			json.Unmarshal(argParamsBytes, &baseArgs)
			baseArgs.ENV.NextStepArgs = baseArgs.ENV.NextStepArgs[1:]
			commandInstance.Next(baseArgs)
		}
	} else {
		log.Println("cant found the command", command)
	}

}
