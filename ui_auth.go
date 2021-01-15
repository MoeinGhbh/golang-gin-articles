package main

import (
	"fmt"
	"git.zeuz.io/devel/olympuz/sdk/golang/zeuztool/cli"
	"git.zeuz.io/sdk/golang/zeuzsdk"
	"strconv"
	"time"
	"net/url"
)

func AutoLogin() error {
	cmd := &CmdAuthLogin{}
	s := cli.GetCurSettings()
	url, err := url.Parse(s.APIEndpoint)
	if err != nil {
		return err
	}
	cmd.Endpoint = *url
	cmd.PWHash = s.DevPassword
	//cmd.User = s.UserID
	//cmd.Dev = s.DevID
	//cmd.ApiKey = s.APIKey
	cmd.Login = s.DevLogin
	return cmd.Run()
}

type CmdAuthLogin struct {

	cli.CmdEndpoint
	zeuzsdk.AuthLoginIn
	Password string `arg:"required" help:"password"`
	PWHash   string
	tries    int
}

type CmdAuthLogout struct {
	cli.CmdSession
}

func (cmd *CmdAuthLogin) Run() error {
	fmt.Println("login:", cmd.Endpoint.String(), cmd.Login)

	cmd.AuthLoginIn.Nonce = zeuzsdk.IDGenerate(zeuzsdk.IDTypeInvalid)
	cmd.AuthLoginIn.Time = zeuzsdk.TSNow()
	sT := strconv.FormatInt(int64(cmd.AuthLoginIn.Time), 10)
	cmd.AuthLoginIn.Hash = zeuzsdk.StringHash(cmd.AuthLoginIn.Nonce + sT + cmd.PWHash)

	//fmt.Println("IsUser", cmd.IsUser)
	//fmt.Println("Time", cmd.Time)
	//fmt.Println("endpoint", cmd.Endpoint)
	//fmt.Println("password", cmd.Password)
	//fmt.Println("PWHash", cmd.PWHash)
	//fmt.Println("tries", cmd.tries)
	//fmt.Println("Hash", cmd.Hash)
	//fmt.Println("Ignorelastlogin", cmd.IgnoreLastLogin)
	//fmt.Println("Isapi", cmd.IsApi)
	//fmt.Println("login", cmd.Login)
	//fmt.Println("nonce", cmd.Nonce)
	//fmt.Println("authloginin", cmd.AuthLoginIn)
	//fmt.Println("sT", sT)

	client := cmd.Client()
	if cmd.IsUser {
		s := cli.GetCurSettings()
		client.ProjID = s.ProjID
		client.EnvID = s.EnvID
	}




	o, err := zeuzsdk.APIAuthLogin(client, &cmd.AuthLoginIn)

	fmt.Println(err)

	if err != nil {
		if err.Error() == zeuzsdk.RequestExpired && cmd.tries < 3 {
			cmd.tries++
			time.Sleep(time.Millisecond * 5) //settle clock
			return cmd.Run()
		}
		return err
	}


	sh := zeuzsdk.SessionFromAuth(o, cmd.PWHash)
	if sh == nil {
		return nil
	}
	fmt.Println("Login Session:", sh.ID)

	s := cli.GetCurSettings()
	s.DevID = o.Dev
	s.UserID = o.User
	s.SessionID = sh.ID
	s.SessionKey = sh.SessionKey
	s.DevPassword = cmd.PWHash
	s.DevLogin = cmd.Login
	s.APIEndpoint = cmd.Endpoint.String()
	cli.SettingsSave()

	return nil
}

func (cmd *CmdAuthLogout) Run() error {
	err := zeuzsdk.APIAuthLogout(cmd.Client())

	fmt.Println(err)

	s := cli.GetCurSettings()
	s.SessionID = ""
	s.SessionKey = ""
	cli.SettingsSave()

	return err
}