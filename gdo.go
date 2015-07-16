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
	matcher        *Matcher
	processor      *Processor
	lines          []string
	matchedLines   []string
	matchedIndexes map[int]bool
}

func NewLines(m *Matcher, p *Processor) *Lines {
	return &Lines{
		matcher:        m,
		processor:      p,
		lines:          []string{},
		matchedLines:   []string{},
		matchedIndexes: make(map[int]bool),
	}
}

func (l *Lines) LoadLines(r io.Reader) error {
	b := bufio.NewScanner(r)
	for i := 0; b.Scan(); i++ {
		line := b.Text()
		if l.matcher.MatchString(line) {
			l.matchedLines = append(l.matchedLines, line)
			l.matchedIndexes[i] = true
		}
		l.lines = append(l.lines, line)
	}
	return b.Err()
}
