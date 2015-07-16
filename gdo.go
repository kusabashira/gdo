package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

type Matcher struct {
	re *regexp.Regexp
}

func NewMatcher(expr string) (m *Matcher, err error) {
	m = &Matcher{}
	m.re, err = regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Matcher) MatchString(s string) bool {
	return m.re.MatchString(s)
}

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

type Lines struct {
	lines          []string
	matchedLines   []string
	matchedIndexes map[int]bool
}

func NewLines() *Lines {
	return &Lines{
		lines:          []string{},
		matchedLines:   []string{},
		matchedIndexes: make(map[int]bool),
	}
}

func (l *Lines) LoadLines(r io.Reader, m *Matcher) error {
	b := bufio.NewScanner(r)
	for i := 0; b.Scan(); i++ {
		line := b.Text()
		if m.MatchString(line) {
			l.matchedLines = append(l.matchedLines, line)
			l.matchedIndexes[i] = true
		}
		l.lines = append(l.lines, line)
	}
	return b.Err()
}

func (l *Lines) Flush(out io.Writer, p *Processor) error {
	if err := p.Process(l.matchedLines); err != nil {
		return err
	}
	mi := 0
	for li := 0; li < len(l.lines); li++ {
		if l.matchedIndexes[li] {
			fmt.Fprintln(out, l.matchedLines[mi])
			mi++
		} else {
			fmt.Fprintln(out, l.lines[li])
		}
	}
	return nil
}
