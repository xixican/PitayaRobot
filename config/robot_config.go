package config

var Config *RobotConfig
var EndPoints string

const (
	SilenceMode   = 0 // 机器人进入场景静止
	NormalMode    = 1 // 机器人正常频率移动停止
	PerpetualMode = 2 //机器人一直移动
	FollowMode    = 3 // 跟随模式
)

const (
	PhoneNum = 100_0000_0000 // 登陆手机号
	SendType = 2             // 2表示手机号登录
	AuthCode = "123456"      // 登陆验证码
	AreaCode = "86"          // 区号
	RobotPid = 4294967296    //机器人pid大于2的32次方，跟玩家id做区分
)

type RobotConfig struct {
	Env          string `yaml:"environment"`
	RobotNum     int    `yaml:"robot_num"`
	UidPre       int    `yaml:"uid_pre"`
	Scene        int    `yaml:"scene"`
	EventId      string `yaml:"event_id"`
	SpaceId      string `yaml:"space_id"`
	Mode         int    `yaml:"mode"`
	FollowPlayer string `yaml:"follow_player"`
	LogLevel     string `yaml:"log_level"`
}

// 节点下次执行随机时间的区间（单位为秒）
const (
	RandomNextMove            = 60
	RandomNextChangeAnimation = 60
	RandomNextChangeImage     = 60
	RandomNextChangeMoveMode  = 60
	RandomNextPrivateChat     = 60
	RandomNextNearChat        = 60
	RandomNextGlobalChat      = 60
)

const (
	Map2D = 1
	Map3D = 2
)

var Avatar3D []string

const AvatarJson = `"avatar":{"code":"2","jsonUrl":"back/cs/skeleton/1/spine/avatar.json","txtUrl":"back/cs/outfit/2/avatar.txt",
"composeUrl":"back/cs/outfit/2/compose.png","showUrl":"back/cs/outfit/2/show.png","skinColorType":1,"skeletonId":"6305c19cb3a32f0b4e2a10f9",
"skeletonCode":"1","units":[{"code":"8","type":4,"images":[{"anchorX":"46.08","anchorY":"7.28","h":101,"w":74,"url":"back/cs/wu/8/right_sleeves_P.png"},
{"anchorX":"35.8","anchorY":"-2.53","h":184,"w":172,"url":"back/cs/wu/8/cloth1_B.png"},{"anchorX":"35.81","anchorY":"-1.53","h":184,"w":172,"url":"back/cs/wu/8/cloth1_F.png"},
{"anchorX":"43.63","anchorY":"0.78","h":185,"w":159,"url":"back/cs/wu/8/cloth1_P.png"},{"anchorX":"35.81","anchorY":"-1.53","h":184,"w":172,"url":"back/cs/wu/8/cloth2_F.png"},
{"anchorX":"45.63","anchorY":"0.78","h":181,"w":159,"url":"back/cs/wu/8/cloth2_P.png"},{"anchorX":"49.23","anchorY":"-4.6","h":105,"w":79,"url":"back/cs/wu/8/left_sleeves_B.png"},
{"anchorX":"50.27","anchorY":"-3.4","h":104,"w":80,"url":"back/cs/wu/8/left_sleeves_F.png"},{"anchorX":"48.82","anchorY":"-5.72","h":105,"w":81,"url":"back/cs/wu/8/left_sleeves_P.png"},
{"anchorX":"49.88","anchorY":"4.67","h":104,"w":81,"url":"back/cs/wu/8/right_sleeves_B.png"},{"anchorX":"49.7","anchorY":"5.35","h":103,"w":80,"url":"back/cs/wu/8/right_sleeves_F.png"}]},
{"code":"9","type":3,"images":[{"anchorX":"17.26","anchorY":"-16.6","h":93,"w":143,"url":"back/cs/wu/9/face_P.png"},{"anchorX":"10.5","anchorY":"0.3","h":92,"w":164,"url":"back/cs/wu/9/face_F.png"}]},
{"code":"10","type":1,"images":[{"anchorX":"55.28","anchorY":"-11.27","h":274,"w":312,"url":"back/cs/wu/10/hair2_P.png"},{"anchorX":"66.1","anchorY":"-0.56","h":272,"w":346,"url":"back/cs/wu/10/hair1_B.png"},
{"anchorX":"65.64","anchorY":"3.94","h":273,"w":345,"url":"back/cs/wu/10/hair1_F.png"},{"anchorX":"42.26","anchorY":"-19.13","h":252,"w":283,"url":"back/cs/wu/10/hair1_P.png"},{"anchorX":"65.64","anchorY":"3.94","h":273,"w":345,"url":"back/cs/wu/10/hair2_F.png"}]},
{"code":"12","type":2,"images":[{"anchorX":"16.92","anchorY":"22.26","h":164,"w":326,"url":"back/cs/wu/12/tire_P.png"},{"anchorX":"-47.65","anchorY":"-0.41","h":190,"w":333,"url":"back/cs/wu/12/tire_B.png"},{"anchorX":"-41.13","anchorY":"-2.35","h":175,"w":333,"url":"back/cs/wu/12/tire_F.png"}]},
{"code":"14","type":5,"images":[{"anchorX":"3.22","anchorY":"8.5","h":45,"w":127,"url":"back/cs/wu/14/trousers_P.png"},{"anchorX":"39.01","anchorY":"2.59","h":128,"w":70,"url":"back/cs/wu/14/left_trousers_B.png"},{"anchorX":"39.01","anchorY":"0.59","h":128,"w":68,"url":"back/cs/wu/14/left_trousers_F.png"},
{"anchorX":"35.61","anchorY":"0.75","h":130,"w":72,"url":"back/cs/wu/14/left_trousers_P.png"},{"anchorX":"41.01","anchorY":"-0.19","h":124,"w":68,"url":"back/cs/wu/14/right_trousers_B.png"},{"anchorX":"39.5","anchorY":"-0.67","h":127,"w":69,"url":"back/cs/wu/14/right_trousers_F.png"},
{"anchorX":"35.81","anchorY":"-2.9","h":130,"w":75,"url":"back/cs/wu/14/right_trousers_P.png"},{"anchorX":"1.5","anchorY":"4.99","h":46,"w":133,"url":"back/cs/wu/14/trousers_B.png"},{"anchorX":"1.5","anchorY":"10.49","h":57,"w":133,"url":"back/cs/wu/14/trousers_F.png"}]}],"skeletonUnits":[{"images":[{"anchorX":"61.53","anchorY":"0.85","h":128,"w":84,"url":"back/cs/skeleton/1/unit/1/right_hand_F.png"},
{"anchorX":"58.81","anchorY":"3.98","h":125,"w":78,"url":"back/cs/skeleton/1/unit/1/right_hand_P.png"},{"anchorX":"41.01","anchorY":"-0.19","h":122,"w":58,"url":"back/cs/skeleton/1/unit/1/right_leg_B.png"},{"anchorX":"41.01","anchorY":"-0.19","h":122,"w":58,"url":"back/cs/skeleton/1/unit/1/right_leg_F.png"},
{"anchorX":"32.85","anchorY":"-1.32","h":134,"w":70,"url":"back/cs/skeleton/1/unit/1/right_leg_P.png"},{"anchorX":"58.78","anchorY":"-3.87","h":118,"w":128,"url":"back/cs/skeleton/1/unit/1/body_B.png"},{"anchorX":"61.28","anchorY":"-3.91","h":123,"w":128,"url":"back/cs/skeleton/1/unit/1/body_F.png"},
{"anchorX":"67.63","anchorY":"-0.72","h":127,"w":118,"url":"back/cs/skeleton/1/unit/1/body_P.png"},{"anchorX":"1.5","anchorY":"15.49","h":51,"w":125,"url":"back/cs/skeleton/1/unit/1/crotch_B.png"},{"anchorX":"1","anchorY":"14.49","h":53,"w":126,"url":"back/cs/skeleton/1/unit/1/crotch_F.png"},
{"anchorX":"0.72","anchorY":"22","h":68,"w":120,"url":"back/cs/skeleton/1/unit/1/crotch_P.png"},{"anchorX":"109.8","anchorY":"-2.61","h":214,"w":264,"url":"back/cs/skeleton/1/unit/1/head_B.png"},{"anchorX":"109.3","anchorY":"-2.6","h":215,"w":264,"url":"back/cs/skeleton/1/unit/1/head_F.png"},
{"anchorX":"106.6","anchorY":"-15","h":214,"w":241,"url":"back/cs/skeleton/1/unit/1/head_P.png"},{"anchorX":"62.54","anchorY":"0.42","h":127,"w":82,"url":"back/cs/skeleton/1/unit/1/left_hand_B.png"},{"anchorX":"62.8","anchorY":"-0.01","h":127,"w":83,"url":"back/cs/skeleton/1/unit/1/left_hand_F.png"},
{"anchorX":"60.84","anchorY":"-1.74","h":126,"w":83,"url":"back/cs/skeleton/1/unit/1/left_hand_P.png"},{"anchorX":"39.01","anchorY":"3.09","h":126,"w":61,"url":"back/cs/skeleton/1/unit/1/left_leg_B.png"},{"anchorX":"39.01","anchorY":"2.59","h":126,"w":62,"url":"back/cs/skeleton/1/unit/1/left_leg_F.png"},
{"anchorX":"33.71","anchorY":"-4.28","h":130,"w":66,"url":"back/cs/skeleton/1/unit/1/left_leg_P.png"},{"anchorX":"61.53","anchorY":"0.85","h":128,"w":84,"url":"back/cs/skeleton/1/unit/1/right_hand_B.png"},{"anchorX":"2.28","anchorY":"-3.5","h":30,"w":163,"url":"back/cs/skeleton/1/unit/shadow.png"}]}]}`
