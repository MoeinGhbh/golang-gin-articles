package main

import (
	"fmt"
	"git.zeuz.io/devel/olympuz/sdk/golang/zeuztool/cli"
	"git.zeuz.io/sdk/golang/zeuzsdk"
)

type CmdProfileUserCreate struct {
	CmdEnvironment
	zeuzsdk.ProfileUserCreateIn
	Password string `arg:"" help:"provide either password or pwhash"`
}

func (cmd *CmdProfileUserCreate) Run() error  {
	if len(cmd.Password) > 0 {
		cmd.PWHash = cli.CalcPWHash(cmd.Login, cmd.Password)
	}
	uid, err := zeuzsdk.APIProfileUsercreate(cmd.Client(), &cmd.ProfileUserCreateIn)
	fmt.Println(err)
	fmt.Println("User created:", uid)
	return err
}
