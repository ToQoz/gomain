package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	exitCode = 0
	tmpl     = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, gomain")
}
`
)

func reportError(err error) {
	exitCode = 2
	fmt.Fprintln(os.Stderr, err.Error())
}

func main() {
	defer os.Exit(exitCode)

	out := bufio.NewWriter(os.Stdout)

	tempDir, err := ioutil.TempDir("", "gomain")
	if err != nil {
		reportError(err)
		return
	}
	defer os.RemoveAll(tempDir)

	// write template contents to main.go
	f, err := os.Create(filepath.Join(tempDir, "main.go"))
	if err != nil {
		reportError(err)
		return
	}
	defer os.Remove(f.Name())
	f.Write([]byte(tmpl))
	f.Close()

	err = launchEditor(f.Name())
	if err != nil {
		reportError(err)
		return
	}

	// show contents of main.go
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		reportError(err)
		return
	}
	out.Write(data)

	// go run main.go && show result
	out.WriteRune('\n')
	out.WriteString("--- Output ---")
	out.WriteRune('\n')
	cmd := exec.Command("go", "run", f.Name())
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Run()

	out.Flush()
}

func launchEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
