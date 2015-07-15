package qy_test

import (
	"fmt"

	"github.com/bigwhite/gowechat/qy"
)

const (
	accessToken = "wx1234abcd"
	layout      = `{
	    "button": [
			{
				"name":"submenu1",
				"subbutton": [
				    {
						"name": "item1",
						"type":"click",
						"key":"s1-item1"
					},
				    {
						"name": "item2",
						"type":"click",
						"key":"s1-item2"
					}
				]
			},
			{
				"name":"submenu2",
				"subbutton": [
				    {
						"name": "item1",
						"type":"click",
						"key":"s2-item1"
					},
				    {
						"name": "item2",
						"type":"click",
						"key":"s2-item2"
					}
				]
			}
		]
} `
	agentID = "5"
)

func ExampleCreateMenu() {
	err := qy.CreateMenu([]byte(layout), accessToken, agentID)
	fmt.Println(err)
	// Output: https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=wx1234abcd&agentid=5
	// invalid access_token
}
