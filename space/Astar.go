package space

import (
	"concurrent-test/logger"
	"concurrent-test/model"
	"concurrent-test/pb"
	"container/heap"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const MoveCost = 1
const MapBorder = 1    //地图边缘（通过这个参数限制随机点尽量靠地图中心的程度，最小为1，最外一圈不能行走）
const RandomRange = 10 //随机点范围(必须大于上面的MapBorder参数，否则在坐标为0时无法随机)

//PointHeap 使用堆存放open，加快查找最小F的效率,需实现heap.Interface接口内的方法(Push,Pop)以及sort.Interface的(Len,Less,Swap)
type PointHeap struct {
	elemSlice []*AStarPoint
}

func (p *PointHeap) Push(i interface{}) {
	p.elemSlice = append(p.elemSlice, i.(*AStarPoint))

}
func (p *PointHeap) Pop() interface{} {
	length := len(p.elemSlice)
	popV := p.elemSlice[length-1]
	p.elemSlice = p.elemSlice[:length-1]
	return popV
}
func (p *PointHeap) Len() int {
	return len(p.elemSlice)
}
func (p *PointHeap) Less(i, j int) bool {
	return p.elemSlice[i].GValue+p.elemSlice[i].HValue < p.elemSlice[j].GValue+p.elemSlice[j].HValue
}
func (p *PointHeap) Swap(i, j int) {
	p.elemSlice[i], p.elemSlice[j] = p.elemSlice[j], p.elemSlice[i]
}

type AStarPoint struct {
	XY          *model.PointXY
	GValue      int
	HValue      int
	FatherPoint *AStarPoint
}

/*
calGValue G表示从起点到当前节点的计算成本
*/
func calGValue(point *AStarPoint) {
	if point.FatherPoint != nil {
		if math.Abs(float64(point.FatherPoint.XY.X-point.XY.X)) == 1 {
			point.GValue += MoveCost
		}
		if math.Abs(float64(point.FatherPoint.XY.Y-point.XY.Y)) == 1 {
			point.GValue += MoveCost
		}
	}
}

/*
calHValue H表示从当前即诶单到终点的估算成本
*/
func calHValue(point *AStarPoint, end *model.PointXY) {
	point.HValue = int(math.Abs(float64(end.X-point.XY.X)) + math.Abs(float64(end.Y-point.XY.Y)))
}

/*
AStarSearchRoad AStar寻路 F = G+ H
返回从终点到起点的坐标切片
*/
func AStarSearchRoad(x int32, y int32, end *model.PointXY, mapGridPlayerMap *sync.Map, mapInfo *model.MapInfo) []*model.PointXY {
	if !canMove(x, y, mapInfo) {
		logger.RobotLog.Errorf("current point : %d,%d can not move!!!", x, y)
		return nil
	}
	start := &AStarPoint{
		XY:          &model.PointXY{X: x, Y: y},
		FatherPoint: nil,
		GValue:      0,
		HValue:      0,
	}
	if end == nil {
		end = randomEndpoint(x, y, mapInfo)
	}
	logger.RobotLog.Debugf("current point:%d,%d, random next point:%v", x, y, end)

	if end == nil {
		logger.RobotLog.Debugf("Random AStarEndPoint return nil,wait for next AStar move")
		return nil
	}

	openHeap := &PointHeap{
		elemSlice: []*AStarPoint{},
	}
	heap.Push(openHeap, start)

	openMap := map[int]*AStarPoint{}
	openMap[pointXYAsMapIndex(start.XY, mapInfo.Width)] = start

	closeMap := map[int]*AStarPoint{}
	//将阻挡点放入closeMap
	obstaclesMap := mapInfo.ObstaclesMap
	if len(obstaclesMap) > 0 {
		for index, obstaclePoint := range obstaclesMap {
			closeMap[int(index)] = &AStarPoint{XY: obstaclePoint}
		}
	}
	//将地图内玩家站位点加入closeMap
	mapGridPlayerMap.Range(func(key, value interface{}) bool {
		playerNum := value.(int)
		if playerNum > 0 {
			index := key.(int32)
			//已经在closeMap中的忽略
			if _, ok := closeMap[int(index)]; ok {
				return true
			} else {
				closeMap[int(index)] = &AStarPoint{XY: &model.PointXY{X: index % mapInfo.Width, Y: index / mapInfo.Width}}
			}
		}
		return true
	})
	//递归寻路
	targetPoint := search(openHeap, openMap, closeMap, end, mapInfo)

	if targetPoint != nil {
		//反向添加点形成路径
		p := targetPoint
		road := []*model.PointXY{p.XY}
		for p.FatherPoint != nil {
			road = append(road, p.FatherPoint.XY)
			p = p.FatherPoint
		}
		//将起点从路径中移除
		road = road[:len(road)-1]
		return road
	}
	return nil
}

//随机寻路终点坐标
func randomEndpoint(x int32, y int32, mapInfo *model.MapInfo) *model.PointXY {
	//限定随机范围
	minX := x - RandomRange
	if minX < MapBorder {
		minX = MapBorder
	}
	maxX := x + RandomRange
	if maxX > mapInfo.Width-MapBorder {
		maxX = mapInfo.Width - MapBorder
	}
	minY := y - RandomRange
	if minY < MapBorder {
		minY = MapBorder
	}
	maxY := y + RandomRange
	if maxY > mapInfo.High-MapBorder {
		maxY = mapInfo.High - MapBorder
	}
	//logger.Log.Infof("%d,%d random:[%d %d], [%d,%d]", x, y, minX, maxX, minY, maxY)
	//通过边界阈值保证尽量少的随机到边缘不可行走区域
	rand.Seed(time.Now().Unix())
	randomX := int32(rand.Intn(int(maxX-minX))) + minX
	randomY := int32(rand.Intn(int(maxY-minY))) + minY
	if canMove(randomX, randomY, mapInfo) {
		return &model.PointXY{X: randomX, Y: randomY}
	}
	if x == randomX && y == randomY {
		return nil
	}
	return nil
}

// CalculateFollowPosition 计算跟随的终点坐标
func CalculateFollowPosition(x, y int32, dir pb.DirectionType, mapInfo *model.MapInfo, mapGridPlayer *sync.Map) *model.PointXY {
	//目标位置为队长身后的位置
	var targetX, targetY int32
	switch dir {
	case pb.DirectionType_LEFT: //反向<- *
		targetX, targetY = x+1, y
	case pb.DirectionType_RIGHT:
		targetX, targetY = x-1, y
	case pb.DirectionType_UP:
		targetX, targetY = x, y+1
	case pb.DirectionType_DOWN:
		targetX, targetY = x, y-1
	}
	//range=0,选择距离target距离为range的点，如果目标点全部被占或不可行走，则range++(<视野大小15)，循环
	for i := 0; i <= 15; i++ {
		var targetArr []*model.PointXY
		//从左往右找
		if dir == pb.DirectionType_RIGHT {
			for startX := int(targetX) - i; startX <= int(targetX)+i; startX++ {
				distanceY := i - int(math.Abs(float64(int(targetX)-startX)))
				//同一x在距离目标点相等的位置有两个y值
				targetArr = append(targetArr, &model.PointXY{X: int32(startX), Y: targetY - int32(distanceY)})
				targetArr = append(targetArr, &model.PointXY{X: int32(startX), Y: targetY + int32(distanceY)})
			}
			for j := 0; j < len(targetArr); j++ {
				point := targetArr[j]
				index := point.Y*mapInfo.Width + point.X
				if _, ok := mapInfo.ObstaclesMap[index]; ok {
					continue
				}
				if playerNum, ok := mapGridPlayer.Load(index); ok {
					if playerNum.(int) > 0 {
						continue
					}
				}
				return point
			}
		}
		//从右往左找
		if dir == pb.DirectionType_LEFT {
			for startX := int(targetX) + i; startX >= int(targetX)-i; startX-- {
				distanceY := i - int(math.Abs(float64(int(targetX)-startX)))
				//同一x在距离目标点相等的位置有两个y值
				targetArr = append(targetArr, &model.PointXY{X: int32(startX), Y: targetY - int32(distanceY)})
				targetArr = append(targetArr, &model.PointXY{X: int32(startX), Y: targetY + int32(distanceY)})
			}
			for j := 0; j < len(targetArr); j++ {
				point := targetArr[j]
				index := point.Y*mapInfo.Width + point.X
				if _, ok := mapInfo.ObstaclesMap[index]; ok {
					continue
				}
				if playerNum, ok := mapGridPlayer.Load(index); ok {
					if playerNum.(int) > 0 {
						continue
					}
				}
				return point
			}
		}
		//从上往下找
		if dir == pb.DirectionType_DOWN {
			for startY := int(targetY) - i; startY <= int(targetY)+i; startY++ {
				distanceX := i - int(math.Abs(float64(int(targetY)-startY)))
				//同一y在距离目标点相等的位置有两个x值
				targetArr = append(targetArr, &model.PointXY{X: targetX - int32(distanceX), Y: int32(startY)})
				targetArr = append(targetArr, &model.PointXY{X: targetX + int32(distanceX), Y: int32(startY)})
			}
			for j := 0; j < len(targetArr); j++ {
				point := targetArr[j]
				index := point.Y*mapInfo.Width + point.X
				if _, ok := mapInfo.ObstaclesMap[index]; ok {
					continue
				}
				if playerNum, ok := mapGridPlayer.Load(index); ok {
					if playerNum.(int) > 0 {
						continue
					}
				}
				return point
			}
		}
		//从下往上找
		if dir == pb.DirectionType_UP {
			for startY := int(targetY) + i; startY >= int(targetY)-i; startY-- {
				distanceX := i - int(math.Abs(float64(int(targetY)-startY)))
				//同一y在距离目标点相等的位置有两个x值
				targetArr = append(targetArr, &model.PointXY{X: targetX - int32(distanceX), Y: int32(startY)})
				targetArr = append(targetArr, &model.PointXY{X: targetX + int32(distanceX), Y: int32(startY)})
			}
			for j := 0; j < len(targetArr); j++ {
				point := targetArr[j]
				index := point.Y*mapInfo.Width + point.X
				if _, ok := mapInfo.ObstaclesMap[index]; ok {
					continue
				}
				if playerNum, ok := mapGridPlayer.Load(index); ok {
					if playerNum.(int) > 0 {
						continue
					}
				}
				return point
			}
		}
	}
	return nil
}

func search(open *PointHeap, openMap map[int]*AStarPoint, closeMap map[int]*AStarPoint, targetXY *model.PointXY, mapInfo *model.MapInfo) *AStarPoint {
	if open.Len() <= 0 {
		return nil
	}
	//在openHeap中弹出F最小点
	minFPoint := heap.Pop(open).(*AStarPoint)
	//从open中移除，加入closeMap
	minFPointIndex := pointXYAsMapIndex(minFPoint.XY, mapInfo.Width)
	delete(openMap, minFPointIndex)
	closeMap[minFPointIndex] = minFPoint
	//获取当前格子相邻的下一步格子切片（只有上下左右四个方向）
	nextPoint := getForwardAStarPoint(minFPoint, closeMap, mapInfo)
	if len(nextPoint) <= 0 {
		//没有下一步路径
		return nil
	}
	for i := 0; i < len(nextPoint); i++ {
		np := nextPoint[i]
		index := pointXYAsMapIndex(np.XY, mapInfo.Width)
		//如果获取的相邻格子不在openMap中，则将它加入openMap
		if _, ok := openMap[index]; !ok {
			//设置父节点
			np.FatherPoint = minFPoint
			//计算G,H
			calGValue(np)
			calHValue(np, targetXY)
			//加入open
			heap.Push(open, np)
			openMap[index] = np
		} else {
			//如果已经在openMap中,用G值衡量从当前minFPoint到该点是否更好
			npG := np.GValue
			//如果G值更小，设置minFPoint为父，重新计算G
			if minFPoint.GValue+1 < npG {
				np.FatherPoint = minFPoint
				calGValue(np)
				//G值变化导致openHeap重新排序
				sort.Sort(open)
			}
		}
		//如果是终点则返回
		if index == pointXYAsMapIndex(targetXY, mapInfo.Width) {
			return np
		}
	}
	//递归
	return search(open, openMap, closeMap, targetXY, mapInfo)
}

func canMove(x int32, y int32, mapInfo *model.MapInfo) bool {
	maxX := mapInfo.Width
	maxY := mapInfo.High
	//边界
	if x < MapBorder || x > maxX-MapBorder || y < MapBorder || y > maxY-MapBorder {
		return false
	}
	//判断是否在阻挡区
	obstacles := mapInfo.Obstacles
	if obstacles == nil {
		return true
	}
	if _, ok := mapInfo.ObstaclesMap[y*maxX+x]; ok {
		return false
	}
	return true
}

/*
*将坐标转换为地图索引
 */
func pointXYAsMapIndex(xy *model.PointXY, mapWidth int32) int {
	//return strconv.Itoa(int(xy.X)) + "," + strconv.Itoa(int(xy.Y))
	return int(xy.Y)*int(mapWidth) + int(xy.X)
}

/*
getForwardAStarPoint 获取当前点可能的的下一步移动点
*/
func getForwardAStarPoint(point *AStarPoint, closeMap map[int]*AStarPoint, mapInfo *model.MapInfo) []*AStarPoint {
	var nextPoint []*AStarPoint
	//如果是阻挡点或已在closeMap中的则忽略（当前阻挡点已在初始化close时加入，即只需判断是否在closeMap中）
	x := point.XY.X
	y := point.XY.Y
	//坐标轴正向：X向右  Y向下
	//上
	if y-1 >= MapBorder {
		p := &AStarPoint{XY: &model.PointXY{X: x, Y: y - 1}}
		index := pointXYAsMapIndex(p.XY, mapInfo.Width)
		if _, ok := closeMap[index]; !ok {
			nextPoint = append(nextPoint, p)
		}
	}
	//下
	if y+1 <= mapInfo.High-MapBorder {
		p := &AStarPoint{XY: &model.PointXY{X: x, Y: y + 1}}
		index := pointXYAsMapIndex(p.XY, mapInfo.Width)
		if _, ok := closeMap[index]; !ok {
			nextPoint = append(nextPoint, p)
		}
	}
	//左
	if x-1 >= MapBorder {
		p := &AStarPoint{XY: &model.PointXY{X: x - 1, Y: y}}
		index := pointXYAsMapIndex(p.XY, mapInfo.Width)
		if _, ok := closeMap[index]; !ok {
			nextPoint = append(nextPoint, p)
		}
	}
	//右
	if x+1 <= mapInfo.Width-MapBorder {
		p := &AStarPoint{XY: &model.PointXY{X: x + 1, Y: y}}
		index := pointXYAsMapIndex(p.XY, mapInfo.Width)
		if _, ok := closeMap[index]; !ok {
			nextPoint = append(nextPoint, p)
		}
	}
	return nextPoint
}
