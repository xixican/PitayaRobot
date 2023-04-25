package space

import (
	"concurrent-test/model"
	"github.com/topfreegames/pitaya/logger"
	"math"
	"math/rand"
)

const (
	FrameWalkDistance = 100 // 3D移动速度为 5000/s, 按50fps, 100/frame
)

// RLineSearchRoad 随机角度直线移动
func RLineSearchRoad(x, y, z, angle int32, mapInfo *model.MapInfo) []*model.PointXYZ {
	var road []*model.PointXYZ
	// 转换为弧度
	radian := math.Pi / 180 * float64(angle)
	// 坐标轴方向的帧移动 （左手坐标系下，规定了客户端z轴的负方向，也即游戏端负y方向为0度，故x值变化为-sin,y值变化为-cos）
	frameXDistance := int32(-FrameWalkDistance * math.Sin(radian))
	frameYDistance := int32(-FrameWalkDistance * math.Cos(radian))
	logger.Log.Infof("x,y,z=(%d,%d,%d), angle=%d, radian=%f, frameXDistance=%d, frameYDistance=%d", x, y, z, angle, radian, frameXDistance, frameYDistance)
	// 以地图长宽中的小者作为单次移动的最大距离限制
	maxDistance := mapInfo.High
	if mapInfo.Width < mapInfo.High {
		maxDistance = mapInfo.High
	}
	randomDistance := rand.Intn(int(maxDistance))
	// 理论帧数
	frameCount := randomDistance / FrameWalkDistance

	for i := 0; i < frameCount; i++ {
		x += frameXDistance
		y += frameYDistance
		if x < 0 || x > mapInfo.Width || y < 0 || y > mapInfo.High {
			break
		}
		road = append(road, &model.PointXYZ{
			X: x,
			Y: y,
			Z: z,
		})
	}
	return road
}
