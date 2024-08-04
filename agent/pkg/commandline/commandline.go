package commandline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Doer interface {
	do(r *bufio.Reader, w io.Writer) error
	getName() string
}

type SelectCmd struct {
	LoopExit bool
	arr      []Doer
	Name     string
}

func (c *SelectCmd) Add(arr ...Doer)                          { c.arr = append(c.arr, arr...) }
func (c *SelectCmd) Start(r *bufio.Reader, w io.Writer) error { return c.do(r, w) }
func (c *SelectCmd) getName() string                          { return c.Name }

func (c *SelectCmd) do(r *bufio.Reader, w io.Writer) error {
	var help string

	for i := range c.arr {
		help += fmt.Sprintf("[%d] - %s\n", i+1, c.arr[i].getName())
	}

	help += fmt.Sprintln("[0] - назад")

	for {
		fmt.Fprintf(w, "[%s]\n", strings.ToUpper(c.Name))
		fmt.Fprint(w, help)
		fmt.Fprint(w, "\nВыбор: ")

		res, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read res: %w", err)
		}

		answer, err := strconv.Atoi(strings.TrimSpace(res))
		if err != nil {
			fmt.Fprintln(w, "неверный ввод")

			continue
		}

		if answer > len(c.arr) {
			fmt.Fprintf(w, "неверный выбор")

			continue
		}

		if answer == 0 {
			return nil
		}

		if err := c.arr[answer-1].do(r, w); err != nil {
			fmt.Fprintf(w, "err: %v\n", err)

			continue
		}

		if c.LoopExit {
			return nil
		}
	}
}

type ExecCmd struct {
	Fn     func(c *ExecCmd) error
	UserIN []string
	buf    []byte
	Name   string
}

func (c *ExecCmd) Get() [][]byte {
	res := bytes.Split(c.buf, []byte("\n"))

	return res[:len(res)-1]
}

func (c *ExecCmd) do(r *bufio.Reader, w io.Writer) error {
	for i := range c.UserIN {
		if _, err := w.Write([]byte(c.UserIN[i])); err != nil {
			return err
		}

		ans, err := r.ReadBytes('\n')
		if err != nil {
			return err
		}

		c.buf = append(c.buf, ans...)
	}

	if c.Fn != nil {
		return c.Fn(c)
	}

	return nil
}

func (c *ExecCmd) getName() string { return c.Name }
