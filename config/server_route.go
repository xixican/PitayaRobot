package config

const (
	SendMyPositionRoute        = "SendMyPosition"        //发送客户端的加入场景坐标位置,pb:SendMyPosition
	NewPositionRoute           = "NewPosition"           //某个人出现在我的视野范围内，B参数为true为代表上新加入的，客户端对应出生动画,pb:Position
	MissPositionRoute          = "MissPosition"          //某个人消失在我的视野范围内,pb:Pid
	OnMoveRoute                = "OnMove"                //某个人在我的视野范围内移动,pb:MovePosition
	OnDirectionRoute           = "Direction"             //我附近的人朝向发生了变化,pb:Direction
	SendPrivateMessageRoute    = "SendPrivateMessage"    //我收到了一条私聊消息,pb:UserMessage
	RollbackRoute              = "Rollback"              //我收到这个路由，表示我的位置出错了，要回退到响应坐标,pb:Rollback
	CaptainPositionUpdateRoute = "CaptainPositionUpdate" //跟随的队长不在视野内时位置更新
)
