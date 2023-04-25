package behaviortree

import (
	"concurrent-test/logger"
	"reflect"
)

var TypeMap = make(map[string]reflect.Type)

//type Node interface {}

type NodeBase struct {
	Name           string      `xml:"Name,attr"`
	Method         string      `xml:"Method,attr"`
	ChildNodeSlice []*NodeBase `xml:"Node"`

	MethodValue reflect.Value
}

func (n *NodeBase) Init() {
	logger.RobotLog.Debugf("%s Init", n.Name)
	t := TypeMap[n.Name]
	if t != nil {
		//reflect.New()返回一个指针的反射对象，而不是直接返回目标对象的反射对象
		realNode := reflect.New(t)
		n.MethodValue = realNode.MethodByName(n.Method)
	}
	if n.ChildNodeSlice == nil {
		return
	}
	childSize := len(n.ChildNodeSlice)
	for i := 0; i < childSize; i++ {
		n.ChildNodeSlice[i].Init()
	}
}

func (n *NodeBase) Run(robot interface{}) bool {
	//logger.Log.Errorf("method point:%d,robot point:%d", &runMethod, robot)
	if !n.MethodValue.IsValid() {
		logger.RobotLog.Errorf("Method -%s- not exist in %s", n.Method, n.Name)
		return false
	}
	if robot == nil {
		return false
	}
	callValue := []reflect.Value{reflect.ValueOf(robot), reflect.ValueOf(n.ChildNodeSlice)}
	ret := n.MethodValue.Call(callValue)
	return ret[0].Bool()
}
