package model

type Ids struct {
	SpaceIds []string
}

// IdsResponse 获取event下的所以space返回
type IdsResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    *Ids   `json:"data"`
}
type PointXY struct {
	X int32
	Y int32
}

type PointXYZ struct {
	X int32
	Y int32
	Z int32
}

type MapInfo struct {
	BirthPlace   []*MapCoordinates
	Width        int32
	High         int32
	Obstacles    []int32
	Privates     []*MapPrivate
	ObstaclesMap map[int32]*PointXY
	Type         int32
}
type MapCoordinates struct {
	X      int32
	Y      int32
	Z      int32
	Weight int32
}
type MapPrivate struct {
	order []int32
}

// MapInfoResponse 获取space的地图数据返回
type MapInfoResponse struct {
	Code    int32    `json:"code"`
	Message string   `json:"message"`
	Data    *MapInfo `json:"data"`
}

type Avatar3D struct {
	Id      string
	Bundle  string
	ShowUrl string
}

type AvatarInfo struct {
	Total int32
	List  []*Avatar3D
}

type Avatar3DResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    *AvatarInfo `json:"data"`
}

// SendCodeRequest 发送验证码
type SendCodeRequest struct {
	Phone    string `json:"phone"`
	SendType int32  `json:"sendType"`
	AreaCode string `json:"areaCode"`
}

// SignInRequest 登陆请求
type SignInRequest struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	AreaCode string `json:"areaCode"`
}

type SignInInfo struct {
	AccessToken              string `json:"accessToken"`
	ExpiresIn                int32  `json:"expiresIn"`
	RefreshToken             string `json:"refreshToken"`
	Scope                    string `json:"scope"`
	EnterUsernamePwdBootPage int32  `json:"enterUsernamePwdBootPage"`
}

// SignInResponse 登陆返回
type SignInResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    *SignInInfo `json:"data"`
}

type UserInfo struct {
	PlatformNo string `json:"platformNo"`
}

// InfoResponse 获取用户信息返回
type InfoResponse struct {
	Code    int32     `json:"code"`
	Message string    `json:"message"`
	Data    *UserInfo `json:"data"`
}

// CreateCardRequest 创建名片请求
type CreateCardRequest struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Title    string `json:"title"`
	AreaCode string `json:"areaCode"`
	Phone    string `json:"phone"`
	WeChatId string `json:"code"`
	Email    string `json:"email"`
}

type CardInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Title       string `json:"title"`
	AreaCode    string `json:"areaCode"`
	Phone       string `json:"phone"`
	WeChatId    string `json:"code"`
	Email       string `json:"email"`
	GmtCreated  string `json:"gmtCreated"`
	GmtModified string `json:"gmtModified"`
}

type MyCardInfo struct {
	Total int32       `json:"total"`
	List  []*CardInfo `json:"list"`
	//PageSize          int32       `json:"pageSize"`
	//PageNum           int32       `json:"pageNum"`
	//Size              int32       `json:"size"`
	//StartRow          int32       `json:"startRow"`
	//EndRow            int32       `json:"endRow"`
	//Pages             int32       `json:"pages"`
	//PrePage           int32       `json:"prePage"`
	//NextPage          int32       `json:"nextPage"`
	//IsFirstPage       bool        `json:"isFirstPage"`
	//IsLastPage        bool        `json:"isLastPage"`
	//HasPreviousPage   bool        `json:"hasPreviousPage"`
	//HasNextPage       bool        `json:"hasNextPage"`
	//NavigatePageNums  []int32     `json:"navigatePageNums"`
	//NavigateFirstPage int32       `json:"navigateFirstPage"`
	//NavigateLastPage  int32       `json:"navigateLastPage"`
}

// GetMyCardResponse 获取我的名片返回
type GetMyCardResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    *MyCardInfo `json:"data"`
}

// SendCardRequest 发送名片请求
type SendCardRequest struct {
	Receiver  string `json:"receiver"`
	CardId    string `json:"cardId"`
	EventId   string `json:"eventId"`
	EventName string `json:"eventName"`
}

type DefaultResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
