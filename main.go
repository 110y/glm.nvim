package main

import (
	"fmt"
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
	// TODO: filename

	go func() {
		l := len(args)
		if l != 2 {
			return
			// return "", errors.New("invalid number of arguments")
		}

		duration := args[1]
		d, err := time.ParseDuration(duration)
		if err != nil {
			return
			// return "", errors.New("invalid duration format")
		}

		for {
			<-time.After(d)
			output, _ := glm.GetImportablePackages()

			file, err := os.OpenFile(args[0], os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				return
				// return "", err
			}

			fmt.Fprint(file, string(output))
		}
	}()

	return "", nil
}
