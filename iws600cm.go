package iws600cm

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// Iws600cm is control of Iws600cm that TokyoDevices produce and sell.
type Iws600cm struct {
	id      string
	libPath string
	debug   bool
}

// LIST represents "list"
const LIST = "list"

// LOOP represents "loop"
const LOOP = "loop"

// ANY represents "ANY"
const ANY = "ANY"

// NewIws600cm creates new instance with control command of Iws600cm.
func NewIws600cm(libPath string) *Iws600cm {
	return &Iws600cm{libPath: libPath}
}

// NewIws600cmSensor creates new instance with sensor ID and control command of Iws600cm.
func NewIws600cmSensor(ID string, libPath string) *Iws600cm {
	return &Iws600cm{id: ID, libPath: libPath}
}

// SetDebug set debug mode.
func (i Iws600cm) SetDebug(b bool) {
	i.debug = b
}

// List executes command "list". Value iws600cm received is sent to channel send of args.
func (i Iws600cm) List() []string {
	// To execute "iws600cm list" returns "exit status 1" and fail every, so ignore err.
	r, _ := exec.Command(i.libPath, LIST).Output()

	var ret []string

	for _, v := range strings.Split(string(r), "\n") {
		ret = append(ret, strings.Trim(v, " "))
	}

	return ret
}

// Loop executes command "loop ID". Value iws600cm received is sent to channel send of args.
// NOTE: Loop doesn't close send and rcv.
// NOTE: Please run on goroutine.
func (i Iws600cm) Loop(ID string, send chan<- string, rcv <-chan bool) error {
	cmd := exec.Command(i.libPath, LOOP, ID)

	return i.loop(cmd, send, rcv)
}

// LoopAny executes command "loop ANY". Value iws600cm received is sent to channel send of args.
// NOTE: LoopAny doesn't close send and rcv.
// NOTE: Please run on goroutine.
func (i Iws600cm) LoopAny(send chan<- string, rcv <-chan bool) error {
	cmd := exec.Command(i.libPath, LOOP, ANY)

	return i.loop(cmd, send, rcv)
}

func (i Iws600cm) loop(c *exec.Cmd, send chan<- string, rcv <-chan bool) error {
	stdout, err := c.StdoutPipe()

	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	var line string

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		select {
		case r := <-rcv:
			if r {
				return nil
			}

		default:
			line = scanner.Text()
			if i.debug {
				fmt.Println("out is : " + line)
			}

			send <- line
		}
	}
	return nil
}
