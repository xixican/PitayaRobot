package node

import (
	"concurrent-test/behaviortree"
	"concurrent-test/config"
	"concurrent-test/logger"
	"concurrent-test/robot"
	"time"
)

const (
	ReconnectSpaceFrequency = 5 * time.Second //退出后重连或传送间隔使用退出时的随机间隔判断

	MoveInterval   = 180 * time.Millisecond
	MoveInterval3D = 20 * time.Millisecond
	MoveFrequency  = 5 * time.Second

	ChangeAnimationFrequency = 20 * time.Second

	ChangeDirectionFrequency = 30 * time.Second

	ChangeMoveModeFrequency = 15 * time.Second

	ChangeImageFrequency = 60 * time.Second

	PrivateFrequency = 2 * time.Minute
	NearFrequency    = 5 * time.Minute
	GlobalFrequency  = 8 * time.Minute

	SendCardFrequency = 10 * time.Second

	FollowFrequency = 15 * time.Second

	ExitFrequency            = 30 * time.Second
	CloseConnectionFrequency = 1 * time.Minute
)

type Condition struct{}

func (c *Condition) ConnectSpaceCondition(robot *robot.Robot, childSLice []*behaviortree.NodeBase) bool {
	if robot.Connected {
		//logger.Log.Infof("%s Connected,spaceId=%s", robot.BaseData.Name, robot.BaseData.SpaceId)
		return false
	}
	if robot.SendConnect {
		return false
	}
	//if time.Since(robot.LastConnectTime) < ReconnectSpaceFrequency {
	//	return false
	//}
	return true
}

func (c *Condition) JoinSpaceCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if robot.Joined {
		//logger.Log.Infof("%s joined,spaceId=%s", robot.BaseData.Name, robot.BaseData.SpaceId)
		return false
	}
	if !robot.Connected {
		logger.RobotLog.Debugf("Must ConnectSpace before JoinSpace")
		return false
	}
	return true
}

func (c *Condition) MoveCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	robot.MoveDataLock.RLock()
	defer robot.MoveDataLock.RUnlock()
	data := robot.MoveData
	//路径为空表示当前robot未处于寻路中，寻路频率拦截
	if (len(data.MoveRoad) == 0 && len(data.MoveRoad3D) == 0) && config.Config.Mode == config.NormalMode && time.Since(data.LastStopTime) < MoveFrequency {
		return false
	}
	//出生后延迟移动
	if time.Since(robot.LastConnectTime) < MoveFrequency {
		return false
	}
	//路径不为空表示robot处于移动寻路中，移动频率拦截
	if len(data.MoveRoad) > 0 && time.Since(data.LastMoveTime) < MoveInterval {
		return false
	}
	if len(data.MoveRoad3D) > 0 && time.Since(data.LastMoveTime) < MoveInterval3D {
		return false
	}
	return true
}
func (c *Condition) ChangeAnimationCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	//robot.MoveDataLock.RLock()
	//defer robot.MoveDataLock.RUnlock()
	data := robot.MoveData
	if time.Since(robot.LastConnectTime) < ChangeAnimationFrequency {
		return false
	}
	if time.Since(data.LastChangeAnimationTime) < ChangeAnimationFrequency {
		return false
	}
	return true
}

func (c *Condition) ChangeDirectionCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	data := robot.MoveData
	if time.Since(robot.LastConnectTime) < ChangeDirectionFrequency {
		return false
	}
	if time.Since(data.LastChangeDirectionTime) < ChangeDirectionFrequency {
		return false
	}
	return true
}

func (c *Condition) ChangeMoveModeCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	//robot.MoveDataLock.RLock()
	//defer robot.MoveDataLock.RUnlock()
	data := robot.MoveData
	if time.Since(robot.LastConnectTime) < ChangeMoveModeFrequency {
		return false
	}
	if time.Since(data.LastChangeMoveModeTime) < ChangeMoveModeFrequency {
		return false
	}
	return true
}

func (c *Condition) ChangeImageCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	//robot.MoveDataLock.RLock()
	//defer robot.MoveDataLock.RUnlock()
	data := robot.MoveData
	if time.Since(robot.LastConnectTime) < ChangeImageFrequency {
		return false
	}
	if time.Since(data.LastChangeImageTime) < ChangeImageFrequency {
		return false
	}
	return true
}

func (c *Condition) PrivateChatCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < PrivateFrequency {
		return false
	}
	if time.Since(robot.ChatData.LastPrivateChatTime) < PrivateFrequency {
		return false
	}
	return true
}
func (c *Condition) NearChatCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < NearFrequency {
		return false
	}
	if time.Since(robot.ChatData.LastNearChatTime) < NearFrequency {
		return false
	}
	return true
}
func (c *Condition) GlobalChatCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < GlobalFrequency {
		return false
	}
	if time.Since(robot.ChatData.LastGlobalChatTime) < GlobalFrequency {
		return false
	}
	return true
}

func (c *Condition) SendCardCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	//需要token,生产环境不允许
	if robot.Token == "" {
		return false
	}
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < SendCardFrequency {
		return false
	}
	if time.Since(robot.CardHolderData.LastSendCardTime) < SendCardFrequency {
		return false
	}
	return true
}

func (c *Condition) FollowCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if config.Config.Mode != config.FollowMode {
		return false
	}
	if time.Since(robot.MoveData.LastFollowTime) < FollowFrequency {
		return false
	}
	return true
}

func (c *Condition) ExitSpaceCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < ExitFrequency {
		return false
	}
	if time.Since(robot.LastExitTime) < ExitFrequency {
		return false
	}
	return true
}

func (c *Condition) CloseConnectionCondition(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if config.Config.Mode == config.SilenceMode {
		return false
	}
	if time.Since(robot.LastConnectTime) < CloseConnectionFrequency {
		return false
	}
	if time.Since(robot.LastCloseConnTime) < CloseConnectionFrequency {
		return false
	}
	return true
}
