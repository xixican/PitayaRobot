package robot

import (
	"bytes"
	"concurrent-test/behaviortree"
	"concurrent-test/config"
	"concurrent-test/logger"
	"concurrent-test/model"
	"concurrent-test/myclient"
	"concurrent-test/pb"
	"concurrent-test/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/protos"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SpaceIdSlice = make([]string, 0, 1)

// MapInfoMap space地图信息map
var MapInfoMap = &sync.Map{}

type Robot struct {
	SendConnect       bool
	Connected         bool
	Joined            bool
	Born              bool //标记是否出生到地图
	LastConnectTime   time.Time
	LastExitTime      time.Time
	LastCloseConnTime time.Time

	BaseData       *BaseModel
	BaseDataLock   *sync.RWMutex
	MoveData       *MoveModel
	MoveDataLock   *sync.RWMutex
	SpaceData      *SpaceModel
	SpaceDataLock  *sync.RWMutex
	ChatData       *ChatModel
	ChatDataLock   *sync.RWMutex
	CardHolderData *CardHolderModel

	BehaviorTree *behaviortree.NodeBase
	Client       *myclient.Client

	Order int
	Token string
}

func New(order int, addr string, tlsConfig *tls.Config) *Robot {
	uid := int64(config.RobotPid + config.Config.UidPre + order)
	cli := myclient.New(logrus.InfoLevel, 100*time.Millisecond)
	var err error
	if strings.Compare(addr, config.Local) == 0 {
		err = cli.ConnectToWS(addr, fmt.Sprintf("/%s", config.Config.EventId), fmt.Sprintf("token=%d", uid))
	} else {
		err = cli.ConnectToWS(addr, fmt.Sprintf("/%s", config.Config.EventId), fmt.Sprintf("token=%d", uid), tlsConfig)
	}
	if err != nil {
		logger.RobotLog.Error(err.Error())
		return nil
	}
	robotTree := behaviortree.NewTree()
	if robotTree == nil {
		logger.RobotLog.Error("create robot behavior tree failed")
		return nil
	}
	robot := &Robot{
		BaseData: &BaseModel{
			Uuid:    fmt.Sprintf("%d", uid),
			Name:    util.GenerateName(config.Config.Scene),
			EventId: config.Config.EventId,
		},
		SpaceData: &SpaceModel{
			MapGridPlayerMap: &sync.Map{},
		},
		MoveData:       &MoveModel{},
		ChatData:       &ChatModel{},
		CardHolderData: &CardHolderModel{},

		BaseDataLock:  &sync.RWMutex{},
		SpaceDataLock: &sync.RWMutex{},
		MoveDataLock:  &sync.RWMutex{},
		ChatDataLock:  &sync.RWMutex{},

		BehaviorTree: robotTree,
		Client:       cli,
	}
	robot.Order = order
	if config.Config.Env != "release" {
		//robot.SignIn()
	}
	return robot
}

// SignIn 通过order构造不会重复的phoneNum登陆,返回token
func (r *Robot) SignIn() {
	var response []byte
	phone := strconv.FormatInt(int64(config.PhoneNum+config.Config.UidPre+r.Order), 10)
	//发送验证码
	sendCodeReq := &model.SendCodeRequest{
		Phone:    phone,
		SendType: config.SendType,
		AreaCode: config.AreaCode,
	}
	sendCodeReqJson, err := json.Marshal(sendCodeReq)
	if err != nil {
		logger.RobotLog.Error("SignIn request error: ", err.Error())
		return
	}
	sendCodeUrl := config.EndPoints + "/sms/send_code"
	util.HttpPost(sendCodeUrl, "", bytes.NewBuffer(sendCodeReqJson))
	//登陆
	signInReq := &model.SignInRequest{
		Phone:    phone,
		Code:     config.AuthCode,
		AreaCode: config.AreaCode,
	}
	signInReqJson, err := json.Marshal(signInReq)
	if err != nil {
		logger.RobotLog.Error("SignIn request error: ", err.Error())
		return
	}
	signInUrl := config.EndPoints + "/signin/phone"
	response = util.HttpPost(signInUrl, "", bytes.NewBuffer(signInReqJson))
	if response == nil {
		logger.RobotLog.Error("SignIn response error: ", err.Error())
		return
	}
	signInResponse := &model.SignInResponse{}
	err = json.Unmarshal(response, signInResponse)
	if err != nil {
		logger.RobotLog.Error("SignIn response json.Unmarshal error: ", err.Error())
		return
	}
	//通过signInResponse返回设置token数据
	r.Token = "Bearer " + signInResponse.Data.AccessToken
	logger.RobotLog.Infof("%s SignIn response code=%d, message=%s", r.BaseData.Name, signInResponse.Code, signInResponse.Message)

	//设置名字和avatar
	visualReq := []byte(util.BuildVisualJson(r.BaseData.Name))
	visualUrl := config.EndPoints + "/v2/visual"
	response = util.HttpPut(visualUrl, r.Token, bytes.NewBuffer(visualReq))
	if response == nil {
		logger.RobotLog.Error("visual response error: ", err.Error())
		return
	}
	visualResponse := &model.DefaultResponse{}
	err = json.Unmarshal(response, visualResponse)
	if err != nil {
		logger.RobotLog.Error("visual response json.Unmarshal error: ", err.Error())
		return
	}
	logger.RobotLog.Infof("%s visual response code=%d, message=%s", r.BaseData.Name, visualResponse.Code, visualResponse.Message)

	//查询生成的平台号并设置
	infoUrl := config.EndPoints + "/info"
	response = util.HttpGet(infoUrl, r.Token)
	if response == nil {
		logger.RobotLog.Error("info response error: ", err.Error())
		return
	}
	infoResponse := &model.InfoResponse{}
	err = json.Unmarshal(response, infoResponse)
	if err != nil {
		logger.RobotLog.Error("info response json.Unmarshal error: ", err.Error())
		return
	}
	r.BaseData.Uuid = infoResponse.Data.PlatformNo
	logger.RobotLog.Infof("%s info response code=%d, message=%s", r.BaseData.Name, infoResponse.Code, infoResponse.Message)
}

func (r *Robot) Start() {
	go r.behaviorTreeRun()
	go r.handleReceiveMsg()
}

func (r *Robot) behaviorTreeRun() {
	//ticker := time.NewTicker(40 * time.Millisecond)
	ticker := time.NewTicker(20 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			r.BehaviorTree.Run(r)
		}
	}
}

/*
handleReceiveMsg 处理收到服务器返回的消息，只在此处修改的数据不需要加锁
*/
func (r *Robot) handleReceiveMsg() {
	for {
		select {
		case receiveMsg := <-r.Client.IncomingMsgChan:
			route := receiveMsg.Route
			data := receiveMsg.Data
			switch route {
			case "":
				result := &protos.Error{}
				err := util.PbUnmarshal(data, result)
				if err == nil {
					if result.Code == "PIT-200" {
						r.Connected = true
					}
				}
			case config.SendMyPositionRoute:
				position := &pb.SendMyPosition{}
				err := util.PbUnmarshal(data, position)
				if err == nil {
					r.UpdateRobotMoveData(func(moveData *MoveModel) {
						moveData.X = position.X
						moveData.Y = position.Y
						moveData.Z = position.Z
						moveData.Angle = position.Angle
					})
					r.Born = true
					logger.RobotLog.Infof("%s born at : (%d,%d,%d) ", r.BaseData.Name, position.X, position.Y, position.Z)
				}
			case config.NewPositionRoute:
				position := &pb.Position{}
				err := util.PbUnmarshal(data, position)
				if err == nil {
					r.SpaceData.PlayerMap.Store(position.Pid, position)
					r.addMapGridPlayer(position.X, position.Y)
					//跟随的队长进视野删除不在视野期间的位置缓存
					if r.MoveData.FollowPid == position.Pid {
						r.SpaceData.CaptainPosition = nil
					}
				}
			case config.MissPositionRoute:
				misPosition := &pb.Pid{}
				err := util.PbUnmarshal(data, misPosition)
				if err == nil {
					position, loaded := r.SpaceData.PlayerMap.LoadAndDelete(misPosition.Pid)
					if loaded {
						viewPosition := position.(*pb.Position)
						//删除地图格子站位
						r.removeMapGridPlayer(viewPosition.X, viewPosition.Y)
					}
				}
			case config.RollbackRoute:
				rollBack := &pb.Rollback{}
				err := util.PbUnmarshal(data, rollBack)
				if err == nil {
					r.UpdateRobotMoveData(func(moveData *MoveModel) {
						logger.RobotLog.Debugf("%s 被回退，从%d,%d回退到到%d,%d", r.BaseData.Name, moveData.X, moveData.Y, rollBack.X, rollBack.Y)
						moveData.X = rollBack.X
						moveData.Y = rollBack.Y
						moveData.LastStopTime = time.Now()
						moveData.MoveRoad = nil
						moveData.RollBacked = true
					})
					//发送停止移动动作
					stopAnimation := &pb.ChangeAnimationReq{A: pb.AnimationType_STAND_F}
					stopAnimationData := util.PbMarshal(stopAnimation)
					r.Client.SendNotify(config.ChangeAnimation, stopAnimationData)
				}
			case config.SendPrivateMessageRoute:
				userMessage := &pb.UserMessage{}
				err := util.PbUnmarshal(data, userMessage)
				if err == nil {
					r.UpdateRobotChatData(func(ChatData *ChatModel) {
						ChatData.PrivateChatSlice = append(ChatData.PrivateChatSlice, userMessage.Pid)
					})
				}
			case config.OnMoveRoute:
				moveMsg := &pb.MovePosition{}
				err := util.PbUnmarshal(data, moveMsg)
				if err == nil {
					position, ok := r.SpaceData.PlayerMap.Load(moveMsg.Pid)
					if ok {
						viewPosition := position.(*pb.Position)
						//设置朝向
						if moveMsg.X < viewPosition.X {
							viewPosition.D = pb.DirectionType_LEFT
						}
						if moveMsg.X > viewPosition.X {
							viewPosition.D = pb.DirectionType_RIGHT
						}
						if moveMsg.Y < viewPosition.Y {
							viewPosition.D = pb.DirectionType_UP
						}
						if moveMsg.Y > viewPosition.Y {
							viewPosition.D = pb.DirectionType_DOWN
						}
						//更新地图站位信息
						//从原站位删除
						r.removeMapGridPlayer(viewPosition.X, viewPosition.Y)
						//新站位增加
						r.addMapGridPlayer(moveMsg.X, moveMsg.Y)
						//更新坐标
						viewPosition.X = moveMsg.X
						viewPosition.Y = moveMsg.Y
						viewPosition.Z = moveMsg.Z
						//更新地图玩家信息
						r.SpaceData.PlayerMap.Store(viewPosition.Pid, viewPosition)
					}
				}
			case config.OnDirectionRoute:
				directionMsg := &pb.Direction{}
				err := util.PbUnmarshal(data, directionMsg)
				if err == nil {
					position, ok := r.SpaceData.PlayerMap.Load(directionMsg.Pid)
					if ok {
						viewPosition := position.(*pb.Position)
						viewPosition.D = directionMsg.D
						//更新数据
						r.SpaceData.PlayerMap.Store(viewPosition.Pid, viewPosition)
					}
				}
			case config.CaptainPositionUpdateRoute:
				captainPositionMsg := &pb.CaptainPositionUpdate{}
				err := util.PbUnmarshal(data, captainPositionMsg)
				if err == nil {
					r.SpaceData.CaptainPosition = captainPositionMsg
				}
			}
		}
	}
}

//robot 停止移动封装(停止移动是多个操作的前提) func

/*--------------------------------------------robot data-------------------------------------------------------------*/

/*UpdateRobotBaseData 更新机器人基础数据*/
func (r *Robot) UpdateRobotBaseData(f func(baseData *BaseModel)) {
	r.BaseDataLock.Lock()
	defer r.BaseDataLock.Unlock()
	f(r.BaseData)
}

func (r *Robot) GetRobotBaseData() *BaseModel {
	r.BaseDataLock.RLock()
	defer r.BaseDataLock.Unlock()
	return r.BaseData
}

/*UpdateSpaceData 更新robot所在space数据*/
func (r *Robot) UpdateSpaceData(f func(spaceData *SpaceModel)) {
	r.SpaceDataLock.Lock()
	defer r.SpaceDataLock.Unlock()
	f(r.SpaceData)
}
func (r *Robot) GetSpaceData() *SpaceModel {
	r.SpaceDataLock.Lock()
	defer r.SpaceDataLock.Unlock()
	return r.SpaceData
}

func (r *Robot) UpdateRobotMoveData(f func(baseData *MoveModel)) {
	r.MoveDataLock.Lock()
	defer r.MoveDataLock.Unlock()
	f(r.MoveData)
}

func (r *Robot) GetRobotMoveData() *MoveModel {
	r.BaseDataLock.RLock()
	defer r.BaseDataLock.Unlock()
	return r.MoveData
}

func (r *Robot) UpdateRobotChatData(f func(ChatData *ChatModel)) {
	r.ChatDataLock.Lock()
	defer r.ChatDataLock.Unlock()
	f(r.ChatData)
}

func (r *Robot) addMapGridPlayer(x, y int32) {
	mapData, ok := MapInfoMap.Load(r.BaseData.SpaceId)
	if !ok {
		return
	}
	mapWidth := mapData.(*model.MapInfo).Width
	playerNum, _ := r.SpaceData.MapGridPlayerMap.LoadOrStore(y*mapWidth+x, 0)
	r.SpaceData.MapGridPlayerMap.Store(y*mapWidth+x, playerNum.(int)+1)
}

func (r *Robot) removeMapGridPlayer(x, y int32) {
	mapData, ok := MapInfoMap.Load(r.BaseData.SpaceId)
	if !ok {
		return
	}
	mapWidth := mapData.(*model.MapInfo).Width
	playerNum, ok := r.SpaceData.MapGridPlayerMap.Load(y*mapWidth + x)
	if !ok {
		return
	}
	r.SpaceData.MapGridPlayerMap.Store(y*mapWidth+x, playerNum.(int)-1)
}

// SpaceModel space数据
type SpaceModel struct {
	PlayerMap        sync.Map                  //视野里的人
	MapGridPlayerMap *sync.Map                 //地图格子对应的玩家 map[index]playerNum
	CaptainPosition  *pb.CaptainPositionUpdate //跟随的队长不在视野内的位置
}

// BaseModel 基础数据
type BaseModel struct {
	Uuid    string
	Name    string
	EventId string
	SpaceId string
}

// MoveModel 移动数据
type MoveModel struct {
	X                       int32
	Y                       int32
	Z                       int32
	Angle                   int32
	Dir                     pb.DirectionType
	RollBacked              bool
	MoveRoad                []*model.PointXY
	MoveRoad3D              []*model.PointXYZ
	LastMoveTime            time.Time
	LastStopTime            time.Time
	Following               bool
	FollowPid               string
	LastFollowTime          time.Time
	NextMoveTime            time.Time
	ShowAnimation           bool
	LastChangeAnimationTime time.Time
	LastChangeDirectionTime time.Time
	IsFlyMode               bool
	LastChangeMoveModeTime  time.Time
	LastChangeImageTime     time.Time
}

// ChatModel 聊天数据
type ChatModel struct {
	LastPrivateChatTime time.Time
	PrivateChatSlice    []string
	LastNearChatTime    time.Time
	LastGlobalChatTime  time.Time
}

// CardHolderModel 名片夹信息
type CardHolderModel struct {
	Created          bool
	CardId           string
	LastSendCardTime time.Time
}
