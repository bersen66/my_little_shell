package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bersen66/my_little_shell/pkg/utils"
	ps "github.com/mitchellh/go-ps"
)

func closeDescriptors(f []io.Closer) {
	for _, file := range f {
		file.Close()
	}
}

type Command interface {
	Run()
	ParseArgs(src []string)
	SetStderr(stderr io.Writer)
	SetStdout(stdout io.Writer)
	SetStdin(stdin io.Reader)
	GetStderr() io.Writer
	GetStdout() io.Writer
	GetStdin() io.Reader
	StdoutPipe() (io.ReadCloser, error)
	Output() ([]byte, error)
}

// Shutdown terminal
type Exit struct {
	Code   int
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Run Exit command
func (e *Exit) Run() {
	os.Exit(e.Code)
}

// Parse Exit command arguments
func (e *Exit) ParseArgs(src []string) {
	if len(src) >= 1 {
		e.Code, _ = strconv.Atoi(src[0])
	} else {
		e.Code = 1
	}
}

func (e *Exit) SetStderr(stderr io.Writer) {
	e.Stderr = stderr
}

func (e *Exit) SetStdout(stdout io.Writer) {
	e.Stdout = stdout
}

func (e *Exit) SetStdin(stdin io.Reader) {
	e.Stdin = stdin
}

func (e *Exit) GetStdin() io.Reader {
	return e.Stdin
}

func (e *Exit) GetStdout() io.Writer {
	return e.Stdout
}

func (e *Exit) GetStderr() io.Writer {
	return e.Stderr
}

// Changes current working dir
type Cd struct {
	Config       *Config
	dest         string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	childIoFiles []io.Closer
}

// Set Stderr for Cd command
func (c *Cd) SetStderr(stderr io.Writer) {
	c.Stderr = stderr
}

// Set Stdout for Cd command
func (c *Cd) SetStdout(stdout io.Writer) {
	c.Stdout = stdout
}

// Set Stdin for Cd command
func (c *Cd) SetStdin(stdin io.Reader) {
	c.Stdin = stdin
}

// Get Stdin of Cd command
func (c *Cd) GetStdin() io.Reader {
	return c.Stdin
}

// Get Stdout of Cd command
func (c *Cd) GetStdout() io.Writer {
	return c.Stdout
}

// Get Stderr of Cd command
func (c *Cd) GetStderr() io.Writer {
	return c.Stderr
}

// Run Cd command
func (c *Cd) Run() {
	defer closeDescriptors(c.childIoFiles)
	var err error

	if c.dest == "" {
		err = os.Chdir(os.Getenv("HOME"))
	} else {
		err = os.Chdir(c.dest)
	}

	if err != nil {
		fmt.Fprintf(c.Stderr, "%v\n", err)
		return
	}

	dir, err := os.Getwd()

	if err != nil {
		fmt.Fprintf(c.Stderr, "%v\n", err)
		return
	}
	c.Config.CurrentDir = dir

}

// Parse Cd command arguments
func (c *Cd) ParseArgs(args []string) {
	if len(args) < 1 {
		c.dest = c.Config.CurrentUser.HomeDir
		return
	}
	c.dest = args[0]
}

// Get StdoutPipe of Cd command
func (c *Cd) StdoutPipe() (io.ReadCloser, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	c.Stdout = pw
	c.childIoFiles = append(c.childIoFiles, pw)
	return pr, nil
}

func (c *Cd) Output() ([]byte, error) {
	c.Run()
	return []byte{}, nil
}

// Clear terminal
type Clear struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Run Clear command
func (c *Clear) Run() {
	utils.ClearCmd()
}

// Parse Clear command arguments
func (c *Clear) ParseArgs(args []string) {}

// Path to current dir
type Pwd struct {
	Config *Config
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	childIoFiles []io.Closer
}

// Run Pwd command
func (p *Pwd) Run() {
	defer closeDescriptors(p.childIoFiles)
	fmt.Fprintf(p.Stdout, "%v\n", p.Config.CurrentDir)
}

// Parse Pwd command arguments
func (p *Pwd) ParseArgs(args []string) {}

func (p *Pwd) Output() ([]byte, error) {
	p.Run()
	return []byte(p.Config.CurrentDir), nil
}

// Set Stderr for Pwd command
func (p *Pwd) SetStderr(stderr io.Writer) {
	p.Stderr = stderr
}

// Set Stdout for Pwd command
func (p *Pwd) SetStdout(stdout io.Writer) {
	p.Stdout = stdout
}

// Set Stdin for Echo command
func (p *Pwd) SetStdin(stdin io.Reader) {
	p.Stdin = stdin
}

// Get Stdin of Pwd command
func (p *Pwd) GetStdin() io.Reader {
	return p.Stdin
}

// Get Stdout of Echo command
func (p *Pwd) GetStdout() io.Writer {
	return p.Stdout
}

// Get Stderr of Echo command
func (p *Pwd) GetStderr() io.Writer {
	return p.Stderr
}

// Get StdoutPipe of Echo command
func (p *Pwd) StdoutPipe() (io.ReadCloser, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	p.Stdout = pw
	p.childIoFiles = append(p.childIoFiles, pw)
	return pr, nil
}

// Echo STDIN to STDOUT
type Echo struct {
	args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	childIoFiles []io.Closer
	result       string
}

// Run Echo command
func (e *Echo) Run() {
	defer closeDescriptors(e.childIoFiles)
	e.result = strings.Join(e.args, " ")
	fmt.Fprintln(e.Stdout, e.result)
}

// Parse Echo command arguments
func (e *Echo) ParseArgs(args []string) {
	e.args = args
}

// Set Stderr for Echo command
func (e *Echo) SetStderr(stderr io.Writer) {
	e.Stderr = stderr
}

// Set Stdout for Echo command
func (e *Echo) SetStdout(stdout io.Writer) {
	e.Stdout = stdout
}

// Set Stdin for Echo command
func (e *Echo) SetStdin(stdin io.Reader) {
	e.Stdin = stdin
}

// Get Stdin of Echo command
func (e *Echo) GetStdin() io.Reader {
	return e.Stdin
}

// Get Stdout of Echo command
func (e *Echo) GetStdout() io.Writer {
	return e.Stdout
}

// Get Stderr of Echo command
func (e *Echo) GetStderr() io.Writer {
	return e.Stderr
}

// Get StdoutPipe of Echo command
func (e *Echo) StdoutPipe() (io.ReadCloser, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	e.Stdout = pw
	e.childIoFiles = append(e.childIoFiles, pw)
	return pr, nil
}

func (e *Echo) Output() ([]byte, error) {
	e.Run()
	return []byte(e.result), nil
}

type Ps struct {
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	childIoFiles []io.Closer
	result       string
}

// Run Ps command
func (p *Ps) Run() {
	defer closeDescriptors(p.childIoFiles)
	processList, _ := ps.Processes()

	b := strings.Builder{}
	for x := range processList {
		var process ps.Process
		process = processList[x]
		b.WriteString(fmt.Sprintf("%d\t%s\n", process.Pid(), process.Executable()))
	}

	p.result = b.String()
	fmt.Fprintln(p.Stdout, p.result)
}

// Parse Ps command arguments
func (p *Ps) ParseArgs(args []string) {
}

// Set Stderr for Ps command
func (p *Ps) SetStderr(stderr io.Writer) {
	p.Stderr = stderr
}

// Set Stdout for Ps command
func (p *Ps) SetStdout(stdout io.Writer) {
	p.Stdout = stdout
}

// Set Stdin for Ps command
func (p *Ps) SetStdin(stdin io.Reader) {
	p.Stdin = stdin
}

// Get Stdin of Ps command
func (p *Ps) GetStdin() io.Reader {
	return p.Stdin
}

// Get Stdout of Ps command
func (p *Ps) GetStdout() io.Writer {
	return p.Stdout
}

// Get Stderr of Ps command
func (p *Ps) GetStderr() io.Writer {
	return p.Stderr
}

// Get StdoutPipe of Ps command
func (p *Ps) StdoutPipe() (io.ReadCloser, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	p.Stdout = pw
	p.childIoFiles = append(p.childIoFiles, pw)
	return pr, nil
}

func (p *Ps) Output() ([]byte, error) {
	p.Run()
	return []byte(p.result), nil
}

type Kill struct {
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	childIoFiles []io.Closer
	pid          int64
	invalid      bool
}

// Run Kill command
func (k *Kill) Run() {
	defer closeDescriptors(k.childIoFiles)
	if k.invalid {
		return
	}
	process, _ := os.FindProcess(int(k.pid))
	process.Kill()
}

// Parse Kill command arguments
func (k *Kill) ParseArgs(args []string) {
	pid, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		k.invalid = true
		return
	}
	k.pid = pid
}

// Set Stderr for Kill command
func (k *Kill) SetStderr(stderr io.Writer) {
	k.Stderr = stderr
}

// Set Stdout for Kill command
func (k *Kill) SetStdout(stdout io.Writer) {
	k.Stdout = stdout
}

// Set Stdin for Kill command
func (k *Kill) SetStdin(stdin io.Reader) {
	k.Stdin = stdin
}

// Get Stdin of Kill command
func (k *Kill) GetStdin() io.Reader {
	return k.Stdin
}

// Get Stdout of Kill command
func (k *Kill) GetStdout() io.Writer {
	return k.Stdout
}

// Get Stderr of Ps command
func (k *Kill) GetStderr() io.Writer {
	return k.Stderr
}

// Get StdoutPipe of Ps command
func (k *Kill) StdoutPipe() (io.ReadCloser, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	k.Stdout = pw
	k.childIoFiles = append(k.childIoFiles, pw)
	return pr, nil
}

func (k *Kill) Output() ([]byte, error) {
	k.Run()
	return []byte{}, nil
}

// External command
type NotBuiltin struct {
	cmd *exec.Cmd
}

// Run NotBuiltin command
func (n *NotBuiltin) Run() {
	n.cmd.Start()
}

func (n *NotBuiltin) ParseArgs(args []string) {
	n.cmd = exec.Command(args[0], args[1:]...)
}

// Set Stderr for NotBuiltin command
func (n *NotBuiltin) SetStderr(stderr io.Writer) {
	n.cmd.Stderr = stderr
}

// Set Stdout for NotBuiltin command
func (n *NotBuiltin) SetStdout(stdout io.Writer) {
	n.cmd.Stdout = stdout
}

// Set Stdin for NotBuiltin command
func (n *NotBuiltin) SetStdin(stdin io.Reader) {
	n.cmd.Stdin = stdin
}

// Get Stdin of NotBuiltin command
func (n *NotBuiltin) GetStdin() io.Reader {
	return n.cmd.Stdin
}

// GetStdout Get Stdout of NotBuiltin command
func (n *NotBuiltin) GetStdout() io.Writer {
	return n.cmd.Stdout
}

// Get Stderr of NotBuiltin command
func (n *NotBuiltin) GetStderr() io.Writer {
	return n.cmd.Stderr
}

// Get StdoutPipe of NotBuiltin command
func (n *NotBuiltin) StdoutPipe() (io.ReadCloser, error) {
	return n.cmd.StdoutPipe()
}

func (n *NotBuiltin) Output() ([]byte, error) {
	return n.cmd.Output()
}
