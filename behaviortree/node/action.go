package node

import (
	"bytes"
	"concurrent-test/behaviortree"
	"concurrent-test/config"
	"concurrent-test/logger"
	"concurrent-test/model"
	"concurrent-test/pb"
	"concurrent-test/robot"
	"concurrent-test/space"
	"concurrent-test/util"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Action struct {
}

func (a *Action) ConnectSpace(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Client.Connected {
		//logger.RobotLog.Errorf("%s lost connection", r.BaseData.Name)
		return false
	}
	result := verifyCondition(r, childSlice)
	//已经Connected则条件不满足，但是当前节点return true
	if !result {
		return true
	}
	//选择进入的space
	spaceIds := robot.SpaceIdSlice
	if len(spaceIds) <= 0 {
		return false
	}
	var enterSpaceId string
	if len(spaceIds) == 1 {
		enterSpaceId = spaceIds[0]
	} else {
		//大概率(60%)进入第一个房间
		//rand.Seed(time.Now().Unix())
		//percent := rand.Intn(10)
		//if percent <= 5 {
		//	enterSpaceId = spaceIds[0]
		//} else {
		//	//%40概率在活动下的其他房间中随机一个
		//	order := rand.Intn(len(spaceIds)-1) + 1
		//	enterSpaceId = spaceIds[order]
		//}
		spaceNum := len(spaceIds)
		index := rand.Intn(spaceNum) //[0,n)
		enterSpaceId = spaceIds[index]
	}
	r.BaseData.SpaceId = enterSpaceId
	msg := &pb.ConnectSpaceReq{
		EventId: r.BaseData.EventId,
		Pid:     r.BaseData.Uuid,
		SpaceId: r.BaseData.SpaceId,
	}
	data := util.PbMarshal(msg)
	r.Client.SendRequest(config.ConnectSpace, data)
	r.SendConnect = true
	r.LastConnectTime = time.Now()
	logger.RobotLog.Infof("%s ConnectSpace to %s ", r.BaseData.Name, enterSpaceId)
	return true
}

func (a *Action) JoinSpace(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	//logger.Log.Errorf("action point:%d,robot point:%d", a, r)
	if !r.Client.Connected {
		//logger.RobotLog.Errorf("%s lost connection", r.BaseData.Name)
		return false
	}
	result := verifyCondition(r, childSlice)
	//已经Joined则条件不满足，但是当前节点return true
	if !result {
		return true
	}
	logger.RobotLog.Infof("%s do JoinSpace,%p", r.BaseData.Name, a)
	mapInfo, ok := robot.MapInfoMap.Load(r.BaseData.SpaceId)
	if !ok {
		logger.RobotLog.Errorf("space %s ,MapInfo not exist", r.BaseData.SpaceId)
		return false
	}
	msg := &pb.JoinReq{
		EventId:  r.BaseData.EventId,
		SpaceId:  r.BaseData.SpaceId,
		Pid:      r.BaseData.Uuid,
		RealName: r.BaseData.Name,
		I:        util.GenerateImage(mapInfo.(*model.MapInfo).Type),
	}
	//logger.RobotLog.Debugf("%s image=%s", r.BaseData.Name, msg.I)
	data := util.PbMarshal(msg)
	r.Client.SendRequest(config.Join, data)
	r.Joined = true
	return true
}

func (a *Action) Move(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//上次移动时间小于移动间隔,移动频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	mapInfoAny, ok := robot.MapInfoMap.Load(r.BaseData.SpaceId)
	if !ok || mapInfoAny == nil {
		logger.RobotLog.Errorf("space %s ,MapInfo not exist", r.BaseData.SpaceId)
		return false
	}
	mapInfo := mapInfoAny.(*model.MapInfo)
	//移动可能会在收到rollback消息时被打断，修改移动数据，故加锁
	r.MoveDataLock.Lock()
	defer r.MoveDataLock.Unlock()
	data := r.MoveData
	//------------------------------------------------------------------------------------------------------------------
	if mapInfo.Type == config.Map3D {
		// 3D random 寻路
		if len(data.MoveRoad3D) == 0 {
			// 360度内随机一个角度
			angle := int32(rand.Intn(360))
			rLine := space.RLineSearchRoad(data.X, data.Y, data.Z, angle, mapInfo)
			if len(rLine) == 0 {
				return false
			}
			data.Angle = angle
			data.MoveRoad3D = rLine
			//打印路径
			var roadBuilder strings.Builder
			for i := 0; i < len(rLine); i++ {
				roadBuilder.WriteString("(")
				roadBuilder.WriteString(strconv.Itoa(int(rLine[i].X)))
				roadBuilder.WriteString(",")
				roadBuilder.WriteString(strconv.Itoa(int(rLine[i].Y)))
				roadBuilder.WriteString(",")
				roadBuilder.WriteString(strconv.Itoa(int(rLine[i].Z)))
				roadBuilder.WriteString(")")
				roadBuilder.WriteString(" ")
			}
			logger.RobotLog.Debugf("%s RLine road: {%s}", r.BaseData.Name, roadBuilder.String())
			// 发送开始移动动画
			changeToMoveAnimationReq := &pb.ChangeAnimationReq{A: pb.AnimationType_WALK_B}
			moveAnimationData := util.PbMarshal(changeToMoveAnimationReq)
			r.Client.SendNotify(config.ChangeAnimation, moveAnimationData)
			logger.RobotLog.Infof("%s send walk 3D animation %v", r.BaseData.Name, changeToMoveAnimationReq.A)
		}
		// 发送移动消息
		next3DPoint := data.MoveRoad3D[0]
		move3DReq := &pb.MoveReq{
			X:     next3DPoint.X,
			Y:     next3DPoint.Y,
			Z:     next3DPoint.Z,
			Angle: data.Angle,
		}
		move3DReqData := util.PbMarshal(move3DReq)
		r.Client.SendNotify(config.MoveAngle, move3DReqData)
		logger.RobotLog.Debugf("%s move 3D  current:%d,%d,%d -> target:%d,%d,%d", r.BaseData.Name, data.X, data.Y, data.Z, next3DPoint.X, next3DPoint.Y, next3DPoint.Z)
		// 更新数据
		data.MoveRoad3D = data.MoveRoad3D[1:]
		data.X = next3DPoint.X
		data.Y = next3DPoint.Y
		data.Z = next3DPoint.Z
		data.LastMoveTime = time.Now()
		//停止移动
		if len(data.MoveRoad3D) == 0 {
			data.LastStopTime = time.Now()
			// 发送停止动作
			changeToStandAnimationReq := &pb.ChangeAnimationReq{A: pb.AnimationType_STAND_B}
			standAnimationData := util.PbMarshal(changeToStandAnimationReq)
			r.Client.SendNotify(config.ChangeAnimation, standAnimationData)
			logger.RobotLog.Infof("%s send stop 3D animation %v", r.BaseData.Name, changeToStandAnimationReq.A)
		}
		return true
	}
	// -----------------------------------------------------------------------------------------------------------------

	// 2D A star 寻路
	//跟随模式A*选取路径
	if r.MoveData.Following {
		var (
			captainX, captainY int32
			captainDir         pb.DirectionType
		)
		position, exist := r.SpaceData.PlayerMap.Load(data.FollowPid)
		if !exist {
			//视野内不存在
			if r.SpaceData.CaptainPosition != nil {
				captainX, captainY = r.SpaceData.CaptainPosition.X, r.SpaceData.CaptainPosition.Y
				captainDir = r.SpaceData.CaptainPosition.D
			} else {
				logger.RobotLog.Debugf("%s Can not follow captain %s, not found in view", r.BaseData.Name, r.MoveData.FollowPid)
				return false
			}
		} else {
			captainX, captainY = position.(*pb.Position).X, position.(*pb.Position).Y
			captainDir = position.(*pb.Position).D
		}
		end := space.CalculateFollowPosition(captainX, captainY, captainDir, mapInfo, r.SpaceData.MapGridPlayerMap)
		if end == nil {
			logger.RobotLog.Debugf("%s Can not follow captain %s, no point", r.BaseData.Name, r.MoveData.FollowPid)
			return false
		}
		//已经在跟随的位置
		if end.X == r.MoveData.X && end.Y == r.MoveData.Y {
			return true
		}
		logger.RobotLog.Debugf("%s follow,captain position->%d,%d , dir->%v, target position->%d,%d, current position->%d,%d",
			r.BaseData.Name, captainX, captainY, captainDir, end.X, end.Y, r.MoveData.X, r.MoveData.Y)
		//地图内玩家站位
		//todo cc sync.map fatal occur on iteration and write遍历和写时报错--原因是赋值操作会创建副本，值类型（sync.Map结构体）赋值都是深拷贝 ，但引用类型（sync.Map内部的m.read.Load().(readOnly).m）都是浅拷贝, 换成*sync.Map
		mapGridPlayerInfo := r.SpaceData.MapGridPlayerMap
		road := space.AStarSearchRoad(data.X, data.Y, end, mapGridPlayerInfo, mapInfo)
		if len(road) > 0 {
			data.MoveRoad = road
			//打印路径
			var roadBuilder strings.Builder
			for i := 0; i < len(road); i++ {
				roadBuilder.WriteString(strconv.Itoa(int(road[i].X)))
				roadBuilder.WriteString(",")
				roadBuilder.WriteString(strconv.Itoa(int(road[i].Y)))
				roadBuilder.WriteString(" ")
			}
			logger.RobotLog.Debugf("%s A* road: {%s}", r.BaseData.Name, roadBuilder.String())
		} else {
			logger.RobotLog.Debugf("%s Can not search road to target point", r.BaseData.Name)
			return false
		}
	} else {
		//随机目的地A*选取路径
		if len(data.MoveRoad) == 0 {
			//地图内玩家站位
			mapGridPlayerInfo := r.SpaceData.MapGridPlayerMap
			road := space.AStarSearchRoad(data.X, data.Y, nil, mapGridPlayerInfo, mapInfo)
			if len(road) > 0 {
				data.MoveRoad = road
				////打印路径
				//var roadBuilder strings.Builder
				//for i := 0; i < len(road); i++ {
				//	roadBuilder.WriteString(string(road[i].X))
				//	roadBuilder.WriteString(",")
				//	roadBuilder.WriteString(string(road[i].Y))
				//	roadBuilder.WriteString(" ")
				//}
				//logger.RobotLog.Debugf("%s A* road: {%s}", r.BaseData.Name, roadBuilder.String())
			} else {
				logger.RobotLog.Debugf("%s Can not search road to target point", r.BaseData.Name)
				return false
			}
		}
	}
	//移动
	lastIndex := len(data.MoveRoad) - 1
	nextPoint := data.MoveRoad[lastIndex]
	logger.RobotLog.Debugf("%s move  current:%d,%d -> target:%d,%d", r.BaseData.Name, data.X, data.Y, nextPoint.X, nextPoint.Y)

	var route string
	//上
	if nextPoint.Y < data.Y {
		route = config.MoveUp
		data.Dir = pb.DirectionType_UP
	}
	//下
	if nextPoint.Y > data.Y {
		route = config.MoveDown
		data.Dir = pb.DirectionType_DOWN
	}
	//左
	if nextPoint.X < data.X {
		route = config.MoveLeft
		data.Dir = pb.DirectionType_LEFT
	}
	//右
	if nextPoint.X > data.X {
		route = config.MoveRight
		data.Dir = pb.DirectionType_RIGHT
	}
	//发送移动1.0
	moveReqMsg := &pb.MoveReq{X: nextPoint.X, Y: nextPoint.Y}
	moveReqData := util.PbMarshal(moveReqMsg)
	r.Client.SendNotify(route, moveReqData)

	//更新数据
	data.MoveRoad = data.MoveRoad[:lastIndex]
	data.X = nextPoint.X
	data.Y = nextPoint.Y
	data.LastMoveTime = time.Now()
	//停止移动
	if len(data.MoveRoad) == 0 {
		data.LastStopTime = time.Now()
	}
	return true
}

func (a *Action) ChangeAnimation(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//改变动作频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	showAnimation := r.MoveData.ShowAnimation
	ChangeAnimationMsg := &pb.ChangeAnimationReq{}
	if showAnimation {
		ChangeAnimationMsg.A = pb.AnimationType_STAND_F
		r.MoveData.ShowAnimation = false
	} else {
		rand.Seed(time.Now().Unix())
		animationCode := rand.Intn(15) //[0,n)
		ChangeAnimationMsg.A = pb.AnimationType(animationCode)
		r.MoveData.ShowAnimation = true
	}
	msgData := util.PbMarshal(ChangeAnimationMsg)
	r.Client.SendNotify(config.ChangeAnimation, msgData)
	r.MoveData.LastChangeAnimationTime = time.Now()
	logger.RobotLog.Infof("%s changeAnimation: %v", r.BaseData.Name, ChangeAnimationMsg.A)
	return true
}

func (a *Action) ChangeDirection(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//改变朝向频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	//在移动时不执行
	if len(r.MoveData.MoveRoad) > 0 {
		return false
	}
	rand.Seed(time.Now().Unix())
	directionCode := rand.Intn(4) //[0,n)
	changeDirectionMsg := &pb.ChangeDirectionReq{D: pb.DirectionType(directionCode)}
	msgData := util.PbMarshal(changeDirectionMsg)
	r.Client.SendNotify(config.ChangeDirection, msgData)
	r.MoveData.LastChangeDirectionTime = time.Now()
	logger.RobotLog.Infof("%s changeDirection: %v", r.BaseData.Name, changeDirectionMsg.D)
	return true
}

func (a *Action) ChangeMoveMode(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//改变移动模式频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	isFly := r.MoveData.IsFlyMode
	if isFly {
		normalModeMsg := &pb.ChangeMoveModeReq{}
		normalModeMsg.Mode = pb.MoveMode_NORMAL_MODE
		data := util.PbMarshal(normalModeMsg)
		r.Client.SendRequest(config.ChangeMoveMode, data)
		logger.RobotLog.Infof("%s changeMoveMode: %v", r.BaseData.Name, pb.MoveMode_NORMAL_MODE)
		//更新数据
		r.MoveData.IsFlyMode = false
		r.MoveData.LastChangeMoveModeTime = time.Now()
	} else {
		//1、周围（以自身为中心的九宫格）有人则开启穿行模式(需要遍历视野中的人)--暂时不用
		//2、被服务器rollBack强拉，则开启穿行
		if r.MoveData.RollBacked {
			flyModeMsg := &pb.ChangeMoveModeReq{}
			flyModeMsg.Mode = pb.MoveMode_FLY_MODE
			data := util.PbMarshal(flyModeMsg)
			r.Client.SendRequest(config.ChangeMoveMode, data)
			logger.RobotLog.Infof("%s changeMoveMode: %v", r.BaseData.Name, pb.MoveMode_FLY_MODE)
			//更新数据
			r.MoveData.IsFlyMode = true
			r.MoveData.LastChangeMoveModeTime = time.Now()
			r.MoveData.RollBacked = false
		} else {
			return false
		}
	}
	return true
}

func (a *Action) ChangeImage(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//换装频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	mapInfo, ok := robot.MapInfoMap.Load(r.BaseData.SpaceId)
	if !ok {
		logger.RobotLog.Errorf("space %s ,MapInfo not exist", r.BaseData.SpaceId)
		return false
	}
	changeImageMsg := &pb.ChangeImageReq{}
	changeImageMsg.Image = util.GenerateImage(mapInfo.(*model.MapInfo).Type)
	changeImageMsgData := util.PbMarshal(changeImageMsg)
	r.Client.SendRequest(config.ChangeImage, changeImageMsgData)
	logger.RobotLog.Infof("%s changeImage: %v", r.BaseData.Name, changeImageMsg.Image)
	//更新数据
	r.MoveData.LastChangeImageTime = time.Now()
	return true
}

func (a *Action) PrivateChat(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//聊天频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	var pid string
	chatContent := util.RandMessage()

	r.ChatDataLock.Lock()
	defer r.ChatDataLock.Unlock()
	chatSlice := r.ChatData.PrivateChatSlice
	if len(chatSlice) > 0 {
		pid = chatSlice[0]
		newChatSlice := chatSlice[1:]
		r.ChatData.PrivateChatSlice = newChatSlice
	} else {
		////从场景周围随机选一人
		//r.SpaceDataLock.RLock()
		//r.SpaceDataLock.RUnlock()
		//playerMap := r.SpaceData.PlayerMap
		//if len(playerMap) > 0 {
		//
		//	for playerPid, _ := range playerMap {
		//		pid = playerPid
		//		break
		//	}
		//} else {
		//	return false
		//}
		return false
	}
	privateMsg := &pb.SendPrivateMessageReq{Pid: pid, Content: chatContent, SpeakType: 0}
	privateMsgData := util.PbMarshal(privateMsg)
	r.Client.SendRequest(config.SendPrivateMessage, privateMsgData)
	r.ChatData.LastPrivateChatTime = time.Now()
	logger.RobotLog.Infof("%s sendPrivateMsg to %s", pid, chatContent)
	return true
}

func (a *Action) NearChat(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//聊天频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	chatContent := util.RandMessage()
	nearMsg := &pb.SendNearMessageReq{Content: chatContent, SpeakType: 0}
	nearMsgData := util.PbMarshal(nearMsg)
	r.Client.SendRequest(config.SendNearMessage, nearMsgData)
	r.ChatData.LastNearChatTime = time.Now()
	logger.RobotLog.Infof("%s sendNearMsg: %s", r.BaseData.Name, chatContent)
	return true
}

func (a *Action) GlobalChat(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//聊天频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	chatContent := util.RandMessage()
	nearMsg := &pb.SendGlobalMessageReq{Content: chatContent, SpeakType: 0}
	nearMsgData := util.PbMarshal(nearMsg)
	r.Client.SendRequest(config.SendGlobalMessage, nearMsgData)
	r.ChatData.LastGlobalChatTime = time.Now()
	logger.RobotLog.Infof("%s sendGlobalMsg: %s", r.BaseData.Name, chatContent)
	return true
}

func (a *Action) SendCard(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//发送名片频率限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	//未创建名片先创建
	if !r.CardHolderData.Created {
		createCardReq := &model.CreateCardRequest{
			Name:     r.BaseData.Name,
			Company:  "VerseTech",
			Title:    "QA",
			AreaCode: "86",
			Phone:    strconv.FormatInt(int64(config.PhoneNum+config.Config.UidPre+r.Order), 10),
			WeChatId: "WeChat" + strconv.Itoa(r.Order),
			Email:    strconv.Itoa(r.Order) + "@163.com",
		}
		reqParamJson, _ := json.Marshal(createCardReq)
		util.HttpPost(config.EndPoints+"/cards", r.Token, bytes.NewBuffer(reqParamJson))
		//更新数据
		r.CardHolderData.Created = true
		r.CardHolderData.LastSendCardTime = time.Now()
		logger.RobotLog.Infof("%s create card ", r.BaseData.Name)
		return true
	}
	//查询我创建的名片
	if r.CardHolderData.CardId == "" {
		url := config.EndPoints + "/cards"
		data := util.HttpGet(url, r.Token)
		if data != nil {
			getMyCardResponse := &model.GetMyCardResponse{}
			json.Unmarshal(data, getMyCardResponse)
			//更细数据
			r.CardHolderData.CardId = getMyCardResponse.Data.List[0].Id
			r.CardHolderData.LastSendCardTime = time.Now()
			logger.RobotLog.Infof("%s get my card, cardId=%s", r.BaseData.Name, r.CardHolderData.CardId)
			return true
		}
	}
	//向视野内玩家随机发名片
	var inViewPlayer []string
	r.SpaceData.PlayerMap.Range(func(key, value interface{}) bool {
		inViewPlayer = append(inViewPlayer, key.(string))
		return true
	})
	if len(inViewPlayer) > 0 {
		receiverId := inViewPlayer[rand.Intn(len(inViewPlayer))]
		sendCardReq := &model.SendCardRequest{
			Receiver:  receiverId,
			CardId:    r.CardHolderData.CardId,
			EventId:   r.BaseData.EventId,
			EventName: r.BaseData.EventId,
		}
		reqParamJson, _ := json.Marshal(sendCardReq)
		data := util.HttpPost(config.EndPoints+"/cards/send", r.Token, bytes.NewBuffer(reqParamJson))
		//更细数据
		r.CardHolderData.LastSendCardTime = time.Now()
		if data != nil {
			defaultResponse := &model.DefaultResponse{}
			json.Unmarshal(data, defaultResponse)
			logger.RobotLog.Infof("%s send card to %s,code=%d,message=%s", r.BaseData.Name, receiverId, defaultResponse.Code, defaultResponse.Message)
		}
	}
	return true
}

func (a *Action) Follow(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//跟随限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	//指定跟随
	if config.Config.FollowPlayer != "" {
		if r.MoveData.Following {
			return true
		}
		startFollowMsg := &pb.StartFollowReq{FollowPid: config.Config.FollowPlayer}
		startFollowData := util.PbMarshal(startFollowMsg)
		r.Client.SendRequest(config.StartFollow, startFollowData)
		//更新数据
		r.MoveData.Following = true
		r.MoveData.FollowPid = config.Config.FollowPlayer
		r.MoveData.LastFollowTime = time.Now()
		logger.RobotLog.Infof("%s startFollow %s!", r.BaseData.Name, config.Config.FollowPlayer)
	} else {
		//随机跟随
		if r.MoveData.Following {
			r.Client.SendRequest(config.CancelFollow, nil)
			//更新数据
			r.MoveData.Following = false
			r.MoveData.LastFollowTime = time.Now()
			logger.RobotLog.Infof("%s cancelFollow!", r.BaseData.Name)
		} else {
			//视野内选择玩家跟随
			var inViewPlayer []string
			r.SpaceData.PlayerMap.Range(func(key, value interface{}) bool {
				inViewPlayer = append(inViewPlayer, key.(string))
				return true
			})
			if len(inViewPlayer) <= 0 {
				return false
			}
			randFollowPid := inViewPlayer[rand.Intn(len(inViewPlayer))]
			startFollowMsg := &pb.StartFollowReq{FollowPid: randFollowPid}
			startFollowData := util.PbMarshal(startFollowMsg)
			r.Client.SendRequest(config.StartFollow, startFollowData)
			//更新数据
			r.MoveData.Following = true
			r.MoveData.FollowPid = randFollowPid
			r.MoveData.LastFollowTime = time.Now()
			logger.RobotLog.Infof("%s startFollow %s!", r.BaseData.Name, randFollowPid)
		}
	}
	return true
}

func (a *Action) ExitSpace(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Born {
		logger.RobotLog.Debugf("%s wait born", r.BaseData.Name)
		return false
	}
	//退出限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	//60%的概率传送退出
	rand.Seed(time.Now().Unix())
	exitPercent := rand.Intn(10)
	if exitPercent > 5 {
		return false
	}
	//清空当前场景数据
	//r.Connected = false
	r.Joined = false
	r.Born = false
	r.SpaceData = &robot.SpaceModel{
		MapGridPlayerMap: &sync.Map{},
	}
	//发送退出消息
	exitMsg := &pb.ExitSpaceReq{SpaceId: r.BaseData.SpaceId, Pid: r.BaseData.Uuid}
	exitMsgData := util.PbMarshal(exitMsg)
	r.Client.SendRequest(config.ExitSpace, exitMsgData)
	logger.RobotLog.Infof("%s exitSpace!", r.BaseData.Name)
	//设置退出时间
	reconnectDelay := time.Duration(rand.Intn(10)) * time.Second
	r.LastExitTime = time.Now().Add(reconnectDelay)
	return true
}

func (a *Action) CloseConnection(r *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	if !r.Connected {
		//logger.RobotLog.Errorf("%s lost connection", r.BaseData.Name)
		return false
	}
	//断开连接时间限制
	result := verifyCondition(r, childSlice)
	if !result {
		return false
	}
	//清空当前场景数据
	r.Connected = false
	r.Joined = false
	r.Born = false
	r.SpaceData = &robot.SpaceModel{
		MapGridPlayerMap: &sync.Map{},
	}
	r.Client.Disconnect()
	//设置断开连接的时间
	r.LastCloseConnTime = time.Now()
	logger.RobotLog.Infof("%s close connection: %d", r.BaseData.Name, time.Now().Unix())
	return true
}

/*-----------------------------------------------internal method-------------------------------------------------------*/
//验证所有条件节点
func verifyCondition(robot *robot.Robot, conditions []*behaviortree.NodeBase) bool {
	if conditions != nil && len(conditions) > 0 {
		childSize := len(conditions)
		//Condition都为true则执行
		for i := 0; i < childSize; i++ {
			childNode := conditions[i]
			if childNode.Name != "Condition" {
				continue
			}
			result := childNode.Run(robot)
			if !result {
				return false
			}
		}
	}
	return true
}
