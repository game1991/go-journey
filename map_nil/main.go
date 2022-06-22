package main

import (
	"encoding/json"
	"fmt"
)

var m = map[string]interface{}{
	"msg": 0,
	"str": "",
}

type StateData struct {
	LocationURL string    `json:"location_url"`
	UserOAuth   string    `json:"user_oauth"` // 用户在应用授权界面勾选的选项
	LoginLog    *LoginLog `json:"login_log"`
}

// LoginLog 登录信息采集
type LoginLog struct {
	// 登录类型：1、web端 2、客户端
	ClientID   string `form:"client_id" json:"client_id,omitempty"`     // 应用的accessid
	ClientName string `form:"client_name" json:"client_name,omitempty"` // 应用名称

	UserID       string `form:"user_id" json:"user_id,omitempty"`              // User表主键ID
	UnionID      string `form:"union_id" json:"union_id,omitempty"`            // UnionID(应用所属公司颁发的id)
	OpenID       string `form:"open_id" json:"open_id,omitempty"`              // OpenID(应用颁发给用户的id)
	ThirdUnionID string `form:"third_unionid" json:"third_union_id,omitempty"` // 三方登录的UnionID(例如微信登录)
}

func main() {
	// m["client"] = models.NewGame()

	// m["client"] = nil

	// v, has := m["client"]
	// if !has || v == nil {
	// 	panic("不存在")
	// }
	// client := v.(*models.Game)

	// fmt.Println(client)

	str := `{
		"location_url": "",
		"login_log": {
			"client_id": "sdasd2121",
			"third_union_id": "121212sadadasd"
		}
	}`

	js(str)
}

func js(src string) error {
	if src != "" {
		var extend *StateData
		if err := json.Unmarshal([]byte(src), &extend); err != nil {
			fmt.Println("第三方平台回调PlatformLoginCallback json.Unmarshal StateData extend失败", err)
			return err
		}
		fmt.Printf("11111111111----------%+v-----%+v\n----111111", *extend, extend.LoginLog)

		if extend != nil && extend.LoginLog != nil {
			extend.LoginLog.ThirdUnionID = "after"
		}
		extendAfter, _ := json.Marshal(extend)

		fmt.Println(string(extendAfter))
	}

	return nil
}
