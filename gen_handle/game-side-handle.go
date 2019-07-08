package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var game_side_file_tmp = `
package game_side_helper

import (
	"cell/common/mustang/log"
	"cell/common/protocol/$GAME_NAME$"
	"cell/td"
	"encoding/binary"
)

/*自动生成的文件，不要手动修改*/
$ALL_HANDLE_DEF$

func init() {
	if all_handle==nil{
		all_handle = make(map[uint32]ICanHandle)
	}

	$ALL_HANDLE_IMPL$
}

`

var game_side_file_tmp_DEF = `

type handle_$PACK_NAME$_$MSG_NAME$ func(uint32, uint64, uint32, *$PACK_NAME$.$MSG_NAME$)
type _$PACK_NAME$_$MSG_NAME$ struct {
	handles   []handle_$PACK_NAME$_$MSG_NAME$
	sendIndex uint32
	recvIndex uint32
}

func (m *_$PACK_NAME$_$MSG_NAME$) getSendIndex() uint32{
	return m.sendIndex+($GAME_ID$<<16)
}

func (m *_$PACK_NAME$_$MSG_NAME$) getRecvIndex() uint32{
	return m.recvIndex+($GAME_ID$<<16)
}

func (m *_$PACK_NAME$_$MSG_NAME$) doHandle(gameID uint32, toID uint64, channel uint32, data []byte) {
	for _, v := range m.handles {
		cmd := &$PACK_NAME$.$MSG_NAME${}
		err := cmd.Unmarshal(data)
		if err != nil {
			log.Error("反序列化消息失败 $PACK_NAME$.$MSG_NAME$ %s", err.Error())
			return
		}
		v(gameID, toID, channel, cmd)
	}
}

func (m *_$PACK_NAME$_$MSG_NAME$) Pack(user *td.$APP_USER$, cmd *$PACK_NAME$.$MSG_NAME$) []byte {
	if m.sendIndex==0 {
		log.Error("序列化消息失败 $PACK_NAME$.$MSG_NAME$ 没有定义发送编号")
		return nil
	}
	size := cmd.Size() + 1 + 4 + Get$APP_USER$Size()
	data := make([]byte, size, size)
	data[0] = uint8(AppEventForwardCmd)	//命令类型(1 Byte)
	binary.LittleEndian.PutUint32(data[1:], m.getSendIndex())	//MsgIndex(4 Bytes) 
	Enc$APP_USER$(user, data[5:])	//AppUser()
	cmd.MarshalTo(data[5+Get$APP_USER$Size():])	//消息体
	return data
}

func (m *_$PACK_NAME$_$MSG_NAME$) RegHandle(f handle_$PACK_NAME$_$MSG_NAME$) {
	m.handles = append(m.handles, f)
}

var $UPPER_PACK_NAME$_$MSG_NAME$ = &_$PACK_NAME$_$MSG_NAME${sendIndex: $SEND$, recvIndex: $RECV$}
`

var game_side_file_tmp_IMPL = `
	if $UPPER_PACK_NAME$_$MSG_NAME$.recvIndex!=0{
		all_handle[$UPPER_PACK_NAME$_$MSG_NAME$.getRecvIndex()] = $UPPER_PACK_NAME$_$MSG_NAME$
	}
`

func gen_game_side(gameName string) {
	fmt.Println("===")
	fileName := viper.GetString("game-side." + gameName + ".file")
	packName := viper.GetString("game-side." + gameName + ".pack")
	appUser := viper.GetString("game-side." + gameName + ".user")
	gameIDVar := viper.GetString("game-side." + gameName + ".gameid")

	fileName = strings.Trim(fileName, " ")
	fmt.Println(fileName)

	if len(fileName) == 0 || len(packName) == 0 {
	}

	str_DEF := ""
	str_IMPL := ""
	game_side_handles := viper.Get("game-side." + gameName + ".handles")
	if game_side_handles != nil {
		for _, v := range game_side_handles.([]interface{}) {
			m := v.(map[string]interface{})
			//fmt.Println(m["cmd"].(string))
			msgName := ""
			recv := "0"
			send := "0"

			if v, ok := m["cmd"]; ok {
				msgName = v.(string)
			}
			if v, ok := m["recv"]; ok {
				recv = fmt.Sprintf("%d", v)
			}
			if v, ok := m["send"]; ok {
				send = fmt.Sprintf("%d", v)
			}

			{
				str := strings.Replace(game_side_file_tmp_DEF, "$MSG_NAME$", msgName, -1)
				str = strings.Replace(str, "$UPPER_PACK_NAME$", strings.ToUpper(packName), -1)
				str = strings.Replace(str, "$PACK_NAME$", packName, -1)
				str = strings.Replace(str, "$APP_USER$", appUser, -1)
				str = strings.Replace(str, "$GAME_ID$", gameIDVar, -1)
				str = strings.Replace(str, "$RECV$", recv, -1)
				str = strings.Replace(str, "$SEND$", send, -1)
				str_DEF += str + "\n"
			}
			{
				str := strings.Replace(game_side_file_tmp_IMPL, "$MSG_NAME$", msgName, -1)
				str = strings.Replace(str, "$UPPER_PACK_NAME$", strings.ToUpper(packName), -1)
				str = strings.Replace(str, "$PACK_NAME$", packName, -1)
				str = strings.Replace(str, "$RECV$", recv, -1)
				str = strings.Replace(str, "$SEND$", send, -1)
				str_IMPL += str + "\n"
			}
		}
	}

	str := game_side_file_tmp
	str = strings.Replace(str, "$GAME_NAME$", packName, -1)
	str = strings.Replace(str, "$ALL_HANDLE_DEF$", str_DEF, -1)
	str = strings.Replace(str, "$ALL_HANDLE_IMPL$", str_IMPL, -1)
	str = strings.Replace(str, "$PACK_NAME$", packName, -1)

	saveToFile(fileName, str)
}
