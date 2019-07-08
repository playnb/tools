package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var fileTemp = `
package gate

import (
	"cell/common/mustang/log"
	"cell/common/mustang/network"
	"cell/common/protocol/msg"
	"cell/service/app-gate/inf"
)

/*
* 自动生成文件，请勿手动修改
*/

func RegisterAllAutoGenHandle(serv network.ICanRegisterHandler) {
	$ALL_REGISTER_HANDLE$
}


`

func saveToFile(fileName string, data string) {
	f, err := os.OpenFile(fileName, os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("[SaveToFile] Error  " + fileName + " " + err.Error())
	} else {
		f.WriteString(data)
	}
}

func main() {
	pflag.String("config", "", "生成的配置文件")
	pflag.String("out_dir", "", "输出文件位置")

	pflag.String("game_side_out_dir", "", "输出文件位置")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(viper.GetString("config"))
	//viper.AddConfigPath(viper.GetString("config"))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	_register_str := ""
	_register_temp := `
	serv.Register(&msg.$MSG_NAME${}, func(agent network.IAgent, msgData interface{}, msgIndex uint16, data interface{}) {
		/*
		* 需要实现函数 func handle_$MSG_NAME$(agent inf.GateAgent, cmd *msg.$MSG_NAME$) {}
		*/
		client := agent.(inf.GateAgent)
		if !filterCmd(client, msgData, msgIndex) {
			client.Terminate()
			log.Trace("[App_Gate] 消息请求失败，断开连接 %s", client)
			return
		}
		handle_$MSG_NAME$(client, msgData.(*msg.$MSG_NAME$))
	})`
	app_gate_handles := viper.Get("app-gate.handles")
	if app_gate_handles != nil {
		for _, v := range app_gate_handles.([]interface{}) {
			m := v.(map[string]interface{})
			//fmt.Println(m["cmd"].(string))
			msgName := ""

			if v, ok := m["cmd"]; ok {
				msgName = v.(string)
			}

			str := strings.Replace(_register_temp, "$MSG_NAME$", msgName, -1)
			_register_str += str + "\n"
		}
	}

	_register_str = strings.Replace(fileTemp, "$ALL_REGISTER_HANDLE$", _register_str, -1)
	//fmt.Println(_register_str)

	saveToFile(viper.GetString("out_dir")+"/auto_gen_handle.go", _register_str)

	gen_game_side("zt2")
	gen_game_side("common")
	gen_game_side("zt2")
	gen_game_side("zthj")
	gen_game_side("ztls")
}
