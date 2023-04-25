package behaviortree

import (
	"concurrent-test/logger"
	"encoding/xml"
	"io/ioutil"
	"os"
)

var treeData []byte

type Config struct {
	XMLName  xml.Name
	RootNode *NodeBase `xml:"Node"`
}

func InitBehaviorTreeXmlConfig() {
	configFile, err := os.Open("./robot_behavior_tree.xml")
	if err != nil {
		logger.RobotLog.Errorf(err.Error())
		return
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.RobotLog.Errorf(err.Error())
		return
	}

	treeData = data
}

func NewTree() *NodeBase {
	tree := &Config{}
	err := xml.Unmarshal(treeData, tree)
	if err != nil {
		logger.RobotLog.Errorf(err.Error())
		return nil
	}
	tree.RootNode.Init()
	return tree.RootNode
}
