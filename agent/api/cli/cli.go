package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	cl "github.com/AndreyVLZ/curly-octo/agent/pkg/commandline"
	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type localStorager interface {
	GetData(ctx context.Context, dataID string) (model.Data, error)
	SaveData(ctx context.Context, data model.Data) error
	List(ctx context.Context) ([]*model.Data, error)
}

type iClient interface {
	Login(ctx context.Context, login, pass string) error
	Registration(ctx context.Context, login, pass string) error
	Send(ctx context.Context) error
	Recv(ctx context.Context) error
}

type Cli struct {
	client iClient
	store  localStorager
}

func New(client iClient, store localStorager) *Cli {
	return &Cli{
		client: client,
		store:  store,
	}
}

func (c *Cli) registration(ctx context.Context) func(*cl.ExecCmd) error {
	return func(cmd *cl.ExecCmd) error {
		vals := cmd.Get()

		err := c.client.Registration(ctx, string(vals[0]), string(vals[1]))
		if err != nil {
			return fmt.Errorf("reg: %w", err)
		}

		return nil
	}
}

func (c *Cli) login(ctx context.Context) func(*cl.ExecCmd) error {
	return func(cmd *cl.ExecCmd) error {
		vals := cmd.Get()

		err := c.client.Login(ctx, string(vals[0]), string(vals[1]))
		if err != nil {
			return fmt.Errorf("login: %w", err)
		}

		return nil
	}
}

func (c *Cli) buidlAddCmd(ctx context.Context, stdOut io.Writer) *cl.SelectCmd {
	logPassUserIN := []string{"имя записи: ", "описание: ", "логин/пароль: "}
	addSelect := &cl.SelectCmd{Name: "добавить запись", LoopExit: true}

	addLogPassCmd := &cl.ExecCmd{
		Name:   "логин/пароль",
		UserIN: logPassUserIN,
		Fn: func(cmd *cl.ExecCmd) error {
			vals := cmd.Get()

			data, err := model.NewLogPassData(string(vals[0]), vals[1], vals[2])
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			if err := c.store.SaveData(ctx, *data); err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Fprintln(stdOut, "запись сохранена в локальное хранилище")

			return nil
		},
	}

	addBinaryCmd := &cl.ExecCmd{
		Name:   "бинарные данные",
		UserIN: []string{"Название: ", "Описание: ", "Путь до файла: "},
		Fn: func(cmd *cl.ExecCmd) error {
			vals := cmd.Get()

			data, err := model.NewBinaryData(string(vals[0]), vals[1], vals[2])
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			if err := c.store.SaveData(ctx, *data); err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Fprintln(stdOut, "запись сохранена в локальное хранилище")

			return nil
		},
	}

	addSelect.Add(addLogPassCmd, addBinaryCmd)

	return addSelect
}

func (c *Cli) buidlEditCmd(ctx context.Context, data *model.Data) *cl.SelectCmd {
	switch data.Type() {
	case model.LogPassData:
		return c.buildLogPassEditCmd(ctx, data)
	case model.BinaryData:
		return c.buildBinaryEditCmd(ctx, data)
	}

	return nil
}

func (c *Cli) buildSendCmd(ctx context.Context) *cl.ExecCmd {
	sendCmd := &cl.ExecCmd{
		Name:   "отправить данные на сервер",
		UserIN: nil,
		Fn: func(_ *cl.ExecCmd) error {
			return c.client.Send(ctx)
		},
	}

	return sendCmd
}

func (c *Cli) buildRecvCmd(ctx context.Context) *cl.ExecCmd {
	recvCmd := &cl.ExecCmd{
		Name:   "получить данные с сервера",
		UserIN: nil,
		Fn: func(_ *cl.ExecCmd) error {
			return c.client.Recv(ctx)
		},
	}

	return recvCmd
}

func (c *Cli) buildBinaryEditCmd(ctx context.Context, data *model.Data) *cl.SelectCmd {
	item := &cl.SelectCmd{Name: data.Name(), LoopExit: true}
	item.Add(
		&cl.ExecCmd{
			Name:   "обновить",
			UserIN: []string{"Название: ", "Описание: ", "Путь до файла: "},
			Fn: func(_ *cl.ExecCmd) error {
				fmt.Printf("обновление записи: %s\n", data.Name())

				return nil
			},
		},
		&cl.ExecCmd{
			Name:   "удалить",
			UserIN: nil,
			Fn: func(_ *cl.ExecCmd) error {
				fmt.Printf("удаление записи: %s\n", data.Name())

				return nil
			},
		},
	)

	return item
}

func (c *Cli) buildLogPassEditCmd(ctx context.Context, data *model.Data) *cl.SelectCmd {
	item := &cl.SelectCmd{Name: data.Name(), LoopExit: true}
	item.Add(
		&cl.ExecCmd{
			Name:   "обновить",
			UserIN: []string{"описание: ", "логин/пароль: "},
			Fn: func(_ *cl.ExecCmd) error {
				fmt.Printf("обновление записи: %s\n", data.Name())

				return nil
			},
		},
		&cl.ExecCmd{
			Name:   "удалить",
			UserIN: nil,
			Fn: func(_ *cl.ExecCmd) error {
				fmt.Printf("удаление записи: %s\n", data.Name())

				return nil
			},
		},
	)

	return item
}

func (c *Cli) Start(ctx context.Context) error {
	stdIn, stdOut := bufio.NewReader(os.Stdin), os.Stdout

	listCmd := &cl.ExecCmd{Name: "список", UserIN: nil,
		Fn: func(cmd *cl.ExecCmd) error {
			updCmd := &cl.SelectCmd{Name: "обновить запись", LoopExit: true}
			arr, err := c.store.List(ctx)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			if len(arr) == 0 {
				fmt.Fprintln(stdOut, "Нет данных в локальном хранилище.")

				return nil
			}

			for i := range arr {
				updCmd.Add(c.buidlEditCmd(ctx, arr[i]))
			}

			return updCmd.Start(stdIn, stdOut)
		},
	}

	main := &cl.SelectCmd{Name: "Главная", LoopExit: false}
	addSelect := c.buidlAddCmd(ctx, stdOut)

	main.Add(
		addSelect,
		listCmd,
		c.buildRecvCmd(ctx),
		c.buildSendCmd(ctx),
	)

	auth := c.buildAuthCmd(ctx)

	if err := auth.Start(stdIn, stdOut); err != nil {
		return fmt.Errorf("start login err: %w\n", err)
	}

	fmt.Fprintln(stdOut, "\nВход выполнен.")

	if err := main.Start(stdIn, stdOut); err != nil {
		fmt.Printf("start list err: %v\n", err)

		return nil
	}

	return nil
}

func (c *Cli) buildAuthCmd(ctx context.Context) *cl.SelectCmd {
	auth := &cl.SelectCmd{Name: "вход", LoopExit: true}
	loginUserIN := []string{"Логин: ", "Пароль: "}

	auth.Add(
		&cl.ExecCmd{
			Name:   "регистрация",
			UserIN: loginUserIN,
			Fn:     c.registration(ctx),
		},
		&cl.ExecCmd{
			Name:   "вход",
			UserIN: loginUserIN,
			Fn:     c.login(ctx),
		},
	)

	return auth
}

/*
func Credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(1)
	fmt.Print("Enter Username: ")
	fmt.Println(2)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Println(3)

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
*/
