package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

		glmfile := args[0]

		duration := args[1]
		d, err := time.ParseDuration(duration)
		if err != nil {
			return
			// return "", errors.New("invalid duration format")
		}

		for {
			modifle, err := os.OpenFile("./go.mod", os.O_RDONLY, 0755)
			if err != nil {
				log.Fatalf("failed to open original modifle: %s\n", err.Error())
				return
			}
			defer modifle.Close()

			copiedModfilePath := fmt.Sprintf("%s/%s", filepath.Dir(glmfile), "go.mod")
			copiedModfile, err := os.OpenFile(copiedModfilePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				// TODO:
				log.Fatalf("failed to open copied modifle: %s\n", err.Error())
				// panic(err)
				return
			}
			defer copiedModfile.Close()

			_, err = io.Copy(copiedModfile, modifle)
			if err != nil {
				log.Fatalf("failed to copy: %s\n", err.Error())
				// panic(err)
				return
			}

			if err := os.Setenv("MODFILE", copiedModfilePath); err != nil {
				return
			}
			if err := os.Setenv("GO_EXECUTABLE", "go1.14beta1"); err != nil {
				return
			}

			output, _ := glm.GetImportablePackages()

			file, err := os.OpenFile(glmfile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				return
				// return "", err
			}
			defer file.Close()

			fmt.Fprint(file, string(output))

			<-time.After(d)
		}
	}()

	return "", nil
}
