package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

type Processor struct {
	cmd *exec.Cmd
}

func NewProcessor(name string, arg ...string) (p *Processor, err error) {
	if _, err = exec.LookPath(name); err != nil {
		return nil, err
	}
	p = &Processor{}
	p.cmd = exec.Command(name, arg...)
	return p, nil
}

func (p *Processor) Process(a []string) error {
	in, err := p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer out.Close()

	if err = p.cmd.Start(); err != nil {
		return err
	}
	for _, s := range a {
		fmt.Fprintln(in, s)
	}
	if err = in.Close(); err != nil {
		return err
	}

	b := bufio.NewScanner(out)
	for i := 0; i < len(a) && b.Scan(); i++ {
		a[i] = b.Text()
	}
	return b.Err()
}
