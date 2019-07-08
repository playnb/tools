package main

import (
	"cell/common/tools/protoc-gen-lrpc-ts/generator"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
)

/*
protoc echo.proto --gogofast_out=../echo --lrpc_out=../out
go build cell\common\tools\protoc-gen-lrpc

protoc echo.proto --gogofast_out=C:\ltp\code\game-server-go\server\src\cell\common\tools\protoc-gen-lrpc-go\echo\ --lrpc-go_out=C:\ltp\code\game-server-go\server\src\cell\common\tools\protoc-gen-lrpc-go\echo\ --lrpc-ts_out=C:\ltp\code\h5\egret-hjc\src\draw\view\scene\net\
*/
func main() {
	g := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.Generate()

	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
