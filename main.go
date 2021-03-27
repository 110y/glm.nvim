package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/110y/glm/pkg/glm"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleFunction(&plugin.FunctionOptions{Name: "StartGLMWorker"}, runGLMWorker)
		return nil
	})
}

func runGLMWorker(args []string) (string, error) {
	l := len(args)
	if l != 2 {
		return "", errors.New("invalid number of arguments")
	}

	glmfile := args[0]

	duration := args[1]
	d, err := time.ParseDuration(duration)
	if err != nil {
		return "", errors.New("invalid duration format")
	}

	go func() {
		<-time.After(3 * time.Second)

		for {
			output, err := glm.GetImportablePackages()
			if err != nil {
				log.Fatalf("failed to execute glm: %s\n", err.Error())
			}

			file, err := os.OpenFile(glmfile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0o755)
			if err != nil {
				return
			}
			defer file.Close()

			fmt.Fprint(file, string(output))

			<-time.After(d)
		}
	}()

	return "", nil
}
