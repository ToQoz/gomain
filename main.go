package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	exitCode    = 0
	defaultTmpl = []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, gomain")
}
`)
)

func reportError(err error) {
	exitCode = 2
	fmt.Fprintln(os.Stderr, err.Error())
}

func main() {
	defer os.Exit(exitCode)

	if err := doMain(nil); err != nil {
		reportError(err)
		return
	}
}

func doMain(tmpl []byte) error {
	out := bufio.NewWriter(os.Stdout)

	tempDir, err := ioutil.TempDir("", "gomain")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// write template contents to main.go
	f, err := os.Create(filepath.Join(tempDir, "main.go"))
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if tmpl == nil {
		f.Write([]byte(defaultTmpl))
	} else {
		f.Write([]byte(tmpl))
	}
	f.Close()

	err = launchEditor(f.Name())
	if err != nil {
		return err
	}

	// show contents of main.go
	writtenCode, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return err
	}
	out.Write(writtenCode)

	// go run main.go && show result
	out.WriteRune('\n')
	out.WriteString("--- Output ---")
	out.WriteRune('\n')
	cmd := exec.Command("go", "run", f.Name())
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Run()

	out.Write([]byte("\nre-edit? y/[N]"))
	out.Flush()

	s, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if isYes(s) {
		return doMain(writtenCode)
	}

	return nil
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

func isYes(txt string) bool {
	txt = strings.Trim(strings.ToUpper(txt), " \n")
	return txt == "Y" || txt == "YES"
}
