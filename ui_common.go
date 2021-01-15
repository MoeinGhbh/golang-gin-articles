package main // import "git.zeuz.io/devel/olympuz/sdk/golang/zeuztool/cli"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"git.zeuz.io/sdk/golang/zeuzsdk"
)

type CmdEndpoint struct {
	Endpoint url.URL `arg:"-endpoint,category" default:"http://127.0.0.1:8080/api/v1" help:"API endpoint URL"`
}

type CmdSession struct {
	CmdEndpoint
	SessionID  zeuzsdk.SessionID `arg:"required,category" help:"session id"`
	SessionKey string            `arg:"required,category" help:"session key"`
}

type CmdProject struct {
	CmdSession
	ProjID zeuzsdk.ProjID `arg:"required,category" help:"project id"`
}

type CmdEnvironment struct {
	CmdProject
	EnvID zeuzsdk.EnvID `arg:"-env,required,category" help:"environment-id"`
}

type CmdOutFile struct {
	Out string `help:"save output to file"`
}
type CmdInFile struct {
	In string `arg:"required" help:"input file or json"`
}
type CmdInFileOpt struct {
	In string `help:"input file or json"`
}

func (cmd *CmdEndpoint) Client() *zeuzsdk.Client {
	return zeuzsdk.NewClient(cmd.Endpoint)
}

func (cmd *CmdSession) Client() *zeuzsdk.Client {
	c := cmd.CmdEndpoint.Client()
	c.SessionID = cmd.SessionID
	c.SessionKey = cmd.SessionKey
	return c
}

func (cmd *CmdProject) Client() *zeuzsdk.Client {
	c := cmd.CmdSession.Client()
	c.ProjID = cmd.ProjID
	s := GetCurSettings()
	s.ProjID = c.ProjID
	return c
}

func (cmd *CmdEnvironment) Client() *zeuzsdk.Client {
	c := cmd.CmdProject.Client()
	c.EnvID = cmd.EnvID
	s := GetCurSettings()
	s.EnvID = c.EnvID
	return c
}

func (cmd *CmdOutFile) HandleOutput(out interface{}) error {
	b, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		return err
	}

	if cmd.Out != "" {
		return ioutil.WriteFile(cmd.Out, b, 0644)
	}

	fmt.Println(string(b))
	return nil
}

func HandleInput(filename string, in interface{}) error {
	if filename == "" {
		return nil
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		err = json.Unmarshal([]byte(filename), in)
		if err != nil {
			return errors.New(filename + " is not a file or cannot be parsed as JSON")
		}
		return nil
	}
	return json.Unmarshal(b, in)
}

func (cmd *CmdInFile) HandleInput(in interface{}) error {
	return HandleInput(cmd.In, in)
}

func (cmd *CmdInFileOpt) HandleInput(in interface{}) error {
	return HandleInput(cmd.In, in)
}
