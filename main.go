package main

import (
	"concurrent-test/behaviortree"
	"concurrent-test/behaviortree/node"
	"concurrent-test/config"
	"concurrent-test/logger"
	"concurrent-test/model"
	"concurrent-test/robot"
	"concurrent-test/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var wsAddress string

//var endPoints string

func main() {
	//pprof
	util.StartHTTPDebug()
	//解析机器人配置
	parseRobotConfig()
	//set logger
	logger.InitRobotLogger(config.Config.LogLevel)
	//行为树初始化
	behaviortree.InitBehaviorTreeXmlConfig()
	registerNode()
	//初始化配置
	initConfig()
	//加载地图数据
	loadSpaceMapInfo()
	//加载3D场景avatar信息
	load3DMapAvatar()

	for i := 1; i <= config.Config.RobotNum; i++ {
		robotClient := robot.New(i, wsAddress, &tls.Config{InsecureSkipVerify: true})
		if robotClient != nil {
			robotClient.Start()
			randSleep := rand.Intn(1000) //1秒内随机睡眠
			//randSleep := 1234
			time.Sleep(time.Duration(randSleep) * time.Millisecond) //用1234这样不规则的延时创建机器人，可以避免两个或以上的机器人在移动节点执行时，因时间戳相同random.seed(time.now().unix())生成相同的寻路终点，由此导致的机器人寻路跟随现象
		} else {
			i--
		}
		//timer := time.NewTimer(2 * time.Second)
		//logger.Log.Infof("before %v", time.Now())
		//<-timer.C
		//logger.Log.Infof("after %v", time.Now())
	}
	logger.RobotLog.Errorf("%d 个机器人添加完毕", config.Config.RobotNum)

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-sg:
		logger.RobotLog.Infof("shut down!!!")
	}
}

func parseRobotConfig() {

	file, err := os.Open("./robot_config.yaml")
	if err != nil {
		panic(fmt.Errorf("robot config err: %v", err.Error()))
	}
	defer file.Close()
	robotConfigData, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	config.Config = &config.RobotConfig{}
	err = yaml.Unmarshal(robotConfigData, config.Config)

	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	//设置endPoints
	setEnvEndpoints()
}
func setEnvEndpoints() {
	switch config.Config.Env {
	case "dev":
		wsAddress = config.Dev
		config.EndPoints = config.EndpointsDev
	case "beta":
		wsAddress = config.Beta
		config.EndPoints = config.EndpointsBeta
	case "release":
		wsAddress = config.Release
		config.EndPoints = config.EndpointsRelease
	case "us":
		wsAddress = config.US
		config.EndPoints = config.EndPointsUS
	case "local":
		wsAddress = config.Local
		config.EndPoints = config.EndpointsLocal
	default:
		panic("env config error")
	}
}

func registerNode() {
	behaviortree.TypeMap["Action"] = reflect.TypeOf(node.Action{})
	behaviortree.TypeMap["Composite"] = reflect.TypeOf(node.Composite{})
	behaviortree.TypeMap["Condition"] = reflect.TypeOf(node.Condition{})
}

func initConfig() {
	if config.Config.EventId == "" {
		panic("eventId need config ")
		return
	}
	if config.Config.SpaceId == "" {
		url := config.EndPoints + "/event/internal/info/" + config.Config.EventId
		data := util.HttpGet(url, "")
		if data != nil {
			spaceIdRes := &model.IdsResponse{}
			json.Unmarshal(data, spaceIdRes)
			robot.SpaceIdSlice = spaceIdRes.Data.SpaceIds
			logger.RobotLog.Infof("event %s have:%v", config.Config.EventId, robot.SpaceIdSlice)
		}
	} else {
		robot.SpaceIdSlice = append(robot.SpaceIdSlice, config.Config.SpaceId)
	}
}

//初始化加载地图数据
func loadSpaceMapInfo() {
	if len(robot.SpaceIdSlice) > 0 {
		spaceNum := len(robot.SpaceIdSlice)
		for i := 0; i < spaceNum; i++ {
			spaceId := robot.SpaceIdSlice[i]
			initMap(spaceId)
		}
	}
}

func initMap(spaceId string) {
	url := config.EndPoints + "/space/expand/" + spaceId
	data := util.HttpGet(url, "")
	spaceMapInfo := &model.MapInfoResponse{}
	json.Unmarshal(data, spaceMapInfo)
	if spaceMapInfo.Data == nil {
		logger.RobotLog.Error("出错咯-_- 请检查配置文件中的eventId、spaceId和environment是否正确^_^ ")
		return
	}
	mapData := spaceMapInfo.Data
	robot.MapInfoMap.Store(spaceId, mapData)
	//将阻挡信息缓存为map加快查询
	mapData.ObstaclesMap = map[int32]*model.PointXY{}
	size := len(mapData.Obstacles)
	for i := 0; i < size; i++ {
		obstacle := mapData.Obstacles[i]
		mapData.ObstaclesMap[obstacle] = &model.PointXY{X: obstacle % mapData.Width, Y: obstacle / mapData.Width}
	}
	logger.RobotLog.Infof("Init space=%s map info", spaceId)
	/*------------------------------------------------------------*/
	if mapData.Type == config.Map2D {
		printSpaceMap(mapData)
	}
}

/*
	------------------------------打印地图--------------------------
*/
func printSpaceMap(mapData *model.MapInfo) {
	w := int(mapData.Width)
	h := int(mapData.High)
	b := mapData.Obstacles
	bMap := map[int]interface{}{}
	for i := 0; i < len(b); i++ {
		bMap[int(b[i])] = nil
	}
	logger.RobotLog.Infof("SpaceMap:W->%d, H->%d, obstacles->%v", w, h, mapData.Obstacles)
	logger.RobotLog.Infof("SpaceMap:W->%d, H->%d, birthPlace->%v", w, h, mapData.BirthPlace)
	/*
		前景 背景 颜色
		30  40  黑色
		31  41  红色
		32  42  绿色
		33  43  黄色
		34  44  蓝色
		35  45  紫红色
		36  46  青蓝色
		37  47  白色
	*/
	//fmt.Printf("%c[0;41;37m%s%c[0m\n", 0x1B, "testPrintColor", 0x1B)
	for i := 0; i < h; i++ { //高度
		for j := 0; j < w; j++ { //宽度
			order := i*w + j
			if _, ok := bMap[order]; ok {
				// 阻挡点打印设置背景红色
				fmt.Printf("%c[0;41;37m%s%c[0m", 0x1B, "****", 0x1B)
			} else {
				//10000以内打印对齐
				if order < 10 {
					fmt.Print("000")
				}
				if 10 <= order && order < 100 {
					fmt.Print("00")
				}
				if 100 <= order && order < 1000 {
					fmt.Print("0")
				}
				fmt.Print(order)
			}
			fmt.Print("  ")
		}
		fmt.Println()
	}
}

// 初始化加载3D场景avatar
func load3DMapAvatar() {
	url := config.EndPoints + "/sdk/outfit3d"
	data := util.HttpGet(url, "")
	response := &model.Avatar3DResponse{}
	json.Unmarshal(data, response)
	if response.Data == nil {
		logger.RobotLog.Error("获取3D场景avatar信息出错T_T,找开发者对线吧")
		return
	}
	for _, avatar3D := range response.Data.List {
		var avatarBuilder strings.Builder
		avatarBuilder.WriteString(avatar3D.Id)
		avatarBuilder.WriteString("||")
		avatarBuilder.WriteString(avatar3D.Bundle)
		config.Avatar3D = append(config.Avatar3D, avatarBuilder.String())
	}
}
