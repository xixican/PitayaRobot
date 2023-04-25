package node

import (
	"concurrent-test/behaviortree"
	"concurrent-test/robot"
)

type Composite struct{}

/*
Selector 选择节点，有一个返回true则为true
*/
func (c *Composite) Selector(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	//logger.Log.Infof("%s do Selector", robot.BaseData.Name)
	if childSlice != nil {
		childSize := len(childSlice)
		for i := 0; i < childSize; i++ {
			result := childSlice[i].Run(robot)
			if result {
				return true
			}
		}
	}
	return false
}

/*
Sequence 顺序节点，所有返回成功则为ture,否则为false
*/
func (c *Composite) Sequence(robot *robot.Robot, childSlice []*behaviortree.NodeBase) bool {
	//logger.Log.Infof("%s do Sequence", robot.BaseData.Name)
	if childSlice != nil {
		childSize := len(childSlice)
		for i := 0; i < childSize; i++ {
			result := childSlice[i].Run(robot)
			if !result {
				return false
			}
		}
		return true
	}
	return false
}
