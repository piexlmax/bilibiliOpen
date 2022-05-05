package module

type DanmuRes struct {
	Data struct {
		FansMedalLevel         int    `json:"fans_medal_level"`
		FansMedalName          string `json:"fans_medal_name"`
		FansMedalWearingStatus bool   `json:"fans_medal_wearing_status"`
		GuardLevel             int    `json:"guard_level"`
		Msg                    string `json:"msg"`
		Timestamp              int    `json:"timestamp"`
		Uid                    int    `json:"uid"`
		Uname                  string `json:"uname"`
		Uface                  string `json:"uface"`
		RoomId                 int    `json:"room_id"`
	} `json:"data"`
	Cmd string `json:"cmd"`
}
