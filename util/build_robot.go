package util

import (
	"concurrent-test/config"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func BuildVisualJson(name string) string {
	var builder strings.Builder
	builder.WriteString("{")
	builder.WriteString(`"nickname":"`)
	builder.WriteString(name)
	builder.WriteString(`",`)
	builder.WriteString(config.AvatarJson)
	builder.WriteString("}")
	return builder.String()
}

// GenerateImage 生成场景形象
func GenerateImage(MapType int32) string {
	// 3D形象
	if MapType == config.Map3D {
		randomIndex := rand.Intn(len(config.Avatar3D))
		return config.Avatar3D[randomIndex]
	}
	//2D形象 格式："back/cs/outfit/1/show.png?1-1"
	imageArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 17, 18, 20} //15,16,19套装有问题
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(imageArr), func(i, j int) {
		imageArr[i], imageArr[j] = imageArr[j], imageArr[i]
	})
	//id := rand.Intn(20) + 1 //套装序号1-20
	var builder strings.Builder
	builder.WriteString("back/cs/outfit/")
	builder.WriteString(strconv.Itoa(imageArr[0]))
	builder.WriteString("/show.png?timestamp=")
	builder.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	builder.WriteString("&code=1-")
	builder.WriteString(strconv.Itoa(imageArr[0]))
	return builder.String()
}

var (
	alphabet  = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	firstName = []string{"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋", "沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",
		"陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章", "云", "苏", "潘", "葛", "奚", "范", "彭", "郎", "鲁", "韦", "昌", "马", "苗", "凤", "花", "方", "俞", "任", "袁", "柳",
		"酆", "鲍", "史", "唐", "费", "廉", "岑", "薛", "雷", "贺", "倪", "汤", "滕", "殷", "罗", "毕", "郝", "邬", "安", "常", "乐", "于", "时", "傅", "皮", "卞", "齐", "康", "伍", "余",
		"元", "卜", "顾", "孟", "平", "黄", "和", "穆", "萧", "尹", "姚", "邵", "湛", "汪", "祁", "毛", "禹", "狄", "米", "贝", "明", "臧", "计", "伏", "成", "戴", "谈", "宋", "茅", "庞",
		"熊", "纪", "舒", "屈", "项", "祝", "董", "梁", "杜", "阮", "蓝", "闵", "席", "季", "麻", "强", "贾", "路", "娄", "危", "江", "童", "颜", "郭", "梅", "盛", "林", "刁", "钟", "徐",
		"邱", "骆", "高", "夏", "蔡", "田", "樊", "胡", "凌", "霍", "虞", "万", "支", "柯", "昝", "管", "卢", "莫", "经", "房", "裘", "缪", "干", "解", "应", "宗", "丁", "宣", "贲", "邓",
		//"郁", "单", "杭", "洪", "包", "诸", "左", "石", "崔", "吉", "钮", "龚", "程", "嵇", "邢", "滑", "裴", "陆", "荣", "翁", "荀", "羊", "於", "惠", "甄", "麴", "家", "封", "芮", "羿",
		//"储", "靳", "汲", "邴", "糜", "松", "井", "段", "富", "巫", "乌", "焦", "巴", "弓", "牧", "隗", "山", "谷", "车", "侯", "宓", "蓬", "全", "郗", "班", "仰", "秋", "仲", "伊", "宫",
		//"宁", "仇", "栾", "暴", "甘", "钭", "厉", "戎", "祖", "武", "符", "刘", "景", "詹", "束", "龙", "叶", "幸", "司", "韶", "郜", "黎", "蓟", "薄", "印", "宿", "白", "怀", "蒲", "邰",
		//"从", "鄂", "索", "咸", "籍", "赖", "卓", "蔺", "屠", "蒙", "池", "乔", "阴", "欎", "胥", "能", "苍", "双", "闻", "莘", "党", "翟", "谭", "贡", "劳", "逄", "姬", "申", "扶", "堵",
		//"冉", "宰", "郦", "雍", "舄", "璩", "桑", "桂", "濮", "牛", "寿", "通", "边", "扈", "燕", "冀", "郏", "浦", "尚", "农", "温", "别", "庄", "晏", "柴", "瞿", "阎", "充", "慕", "连",
		//"茹", "习", "宦", "艾", "鱼", "容", "向", "古", "易", "慎", "戈", "廖", "庾", "终", "暨", "居", "衡", "步", "都", "耿", "满", "弘", "匡", "国", "文", "寇", "广", "禄", "阙", "东",
		//"殴", "殳", "沃", "利", "蔚", "越", "夔", "隆", "师", "巩", "厍", "聂", "晁", "勾", "敖", "融", "冷", "訾", "辛", "阚", "那", "简", "饶", "空", "曾", "毋", "沙", "乜", "养", "鞠",
		//"须", "丰", "巢", "关", "蒯", "相", "查", "後", "荆", "红", "游", "竺", "权", "逯", "盖", "益", "桓", "公", "万俟", "司马", "上官", "欧阳", "夏侯", "诸葛", "闻人", "东方",
		//"赫连", "皇甫", "尉迟", "澹台", "公冶", "宗政", "濮阳", "淳于", "单于", "太叔", "申屠", "公孙", "仲孙", "轩辕", "令狐", "钟离", "宇文", "长孙",
		//"慕容", "鲜于", "闾丘", "司徒", "司空", "亓官", "司寇", "仉", "督", "子车", "颛孙", "端木", "巫马", "公西", "漆雕", "乐正", "壤驷", "公良", "拓跋", "夹谷",
		//"宰父", "谷梁", "晋", "楚", "闫", "法", "汝", "鄢", "涂", "钦", "段干", "百里", "东郭", "南门", "呼延", "归", "海", "微生", "岳", "帅", "缑", "亢", "况", "后", "有",
		//"琴", "梁丘", "左丘", "东门", "西门", "商", "牟", "佘", "佴", "伯", "赏", "南宫", "墨", "哈", "谯", "笪", "年", "爱", "阳", "佟", "第五", "言", "福",
	}
	secondName = []string{"涛", "昌", "进", "林", "有", "坚", "和", "彪", "博", "诚", "先", "敬", "震", "振", "壮", "会", "群", "豪", "心", "邦", "承", "乐", "绍", "功", "松", "善", "厚", "庆", "磊", "民", "友", "裕", "河",
		"哲", "江", "超", "浩", "亮", "政", "谦", "亨", "奇", "固", "之", "轮", "翰", "朗", "伯", "宏", "言", "若", "鸣", "朋", "斌", "梁", "栋", "维", "启", "克", "伦", "翔", "旭", "鹏", "泽", "晨", "辰", "士",
		"以", "建", "家", "致", "树", "炎", "德", "行", "时", "泰", "盛", "雄", "琛", "钧", "冠", "策", "腾", "伟", "刚", "勇", "毅", "俊", "峰", "强", "军", "平", "保", "东", "文", "辉", "力", "明", "永", "健",
		"世", "广", "志", "义", "兴", "良", "海", "山", "仁", "波", "宁", "贵", "福", "生", "龙", "元", "全", "国", "胜", "学", "祥", "才", "发", "成", "康", "星", "光", "天", "达", "安", "岩", "中", "茂", "武",
		"新", "利", "清", "飞", "彬", "富", "顺", "信", "子", "杰", "楠", "榕", "风", "航", "弘嘉", "琼", "桂", "娣", "叶", "璧", "璐", "娅", "琦", "晶", "妍", "茜", "秋", "珊", "莎", "锦", "黛", "青", "倩", "婷",
		"姣", "婉", "娴", "瑾", "颖", "露", "瑶", "怡", "婵", "雁", "蓓", "纨", "仪", "荷", "丹", "蓉", "眉", "君", "琴", "蕊", "薇", "菁", "梦", "岚", "苑", "婕", "馨", "瑗", "琰", "韵", "融", "园", "艺", "咏",
		"卿", "聪", "澜", "纯", "毓", "悦", "昭", "冰", "爽", "琬", "茗", "羽", "希", "宁", "欣", "飘", "育", "滢", "馥", "筠", "柔", "竹", "霭", "凝", "晓", "欢", "霄", "枫", "芸", "菲", "寒", "伊", "亚", "宜",
		"可", "姬", "舒", "影", "荔", "枝", "思", "丽", "秀", "娟", "英", "华", "慧", "巧", "美", "娜", "静", "淑", "惠", "珠", "翠", "雅", "芝", "玉", "萍", "红", "娥", "玲", "芬", "芳", "燕", "彩", "春", "菊",
		"勤", "珍", "贞", "莉", "兰", "凤", "洁", "梅", "琳", "素", "云", "莲", "真", "环", "雪", "荣", "爱", "妹", "霞", "香", "月", "莺", "媛", "艳", "瑞", "凡", "佳"}
	collegeName = []string{"南京理工大学", "南京师范大学", "兰州大学", "西南交通大学", "北京科技大学", "华东理工大学", "武汉理工大学", "北京交通大学", "华中师范大学", "河海大学", "西南大学", "暨南大学", "江南大学", "北京工业大学", "西安电子科技大学", "合肥工业大学", "" +
		"东北师范大学", "南昌大学", "南京农业大学", "哈尔滨工程大学", "陕西师范大学", "中国矿业大学", "华中农业大学", "浙江工业大学", "华南师范大学", "中国海洋大学", "华北电力大学", "燕山大学", "江苏大学", "东华大学", "北京邮电大学", "云南大学", "福州大学", "" +
		"中央财经大学", "扬州大学", "宁波大学", "西北大学", "首都医科大学", "北京化工大学", "湖南师范大学", "福建师范大学", "浙江师范大学", "西北农林科技大学", "上海师范大学", "广西大学", "华南农业大学", "北京协和医学院", "河南大学", "安徽大学", "深圳大学", "" +
		"首都师范大学", "上海财经大学", "中国地质大学（武汉）", "湘潭大学", "中国石油大学（华东）", "杭州电子科技大学", "中国政法大学", "西安建筑科技大学", "太原理工大学", "河北大学", "中央民族大学", "北京林业大学", "长安大学", "广东工业大学", "昆明理工大学", "" +
		"南京工业大学", "青岛大学", "中国传媒大学", "山西大学", "南京医科大学", "天津医科大学", "对外经济贸易大学", "南方医科大学", "山东师范大学", "北京语言大学", "南京邮电大学", "中央美术学院", "浙江理工大学", "中央音乐学院", "中南财经政法大学", "华侨大学", "" +
		"广州大学", "西南财经大学", "福建农林大学", "济南大学", "东北农业大学", "中国医科大学", "中国石油大学（北京）", "山东科技大学", "哈尔滨医科大学", "天津师范大学", "河北工业大学", "中国美术学院", "中国音乐学院", "上海中医药大学", "东北财经大学", "南京信息工程大学", "" +
		"中国矿业大学（北京）", "第二军医大学", "第四军医大学", "东北林业大学", "黑龙江大学", "贵州大学", "上海理工大学", "新疆大学", "大连海事大学", "中国地质大学（北京）", "浙江工商大学", "西安理工大学", "北京中医药大学", "内蒙古大学", "中北大学", "辽宁大学", "" +
		"长沙理工大学", "安徽师范大学", "西北师范大学", "云南师范大学", "北京体育大学", "海南大学", "温州医科大学", "广西师范大学", "中国药科大学", "江西师范大学", "重庆邮电大学", "江苏师范大学", "河南科技大学", "河南师范大学", "中国计量大学", "江西财经大学", "" +
		"西南石油大学", "武汉科技大学", "哈尔滨理工大学", "杭州师范大学", "长春理工大学", "天津工业大学", "河南理工大学", "上海体育学院", "上海外国语大学", "四川农业大学", "山东理工大学", "河北师范大学", "南京中医药大学", "三峡大学", "四川师范大学", "山东农业大学", "" +
		"重庆医科大学", "南通大学", "石河子大学", "北京外国语大学", "南京林业大学", "南京艺术学院", "广东外语外贸大学", "青岛科技大学", "湖北大学", "汕头大学", "上海海事大学", "成都理工大学", "温州大学", "中南民族大学", "西南科技大学", "兰州理工大学", "天津科技大学", "" +
		"上海音乐学", "华东政法大学", "曲阜师范大学", "重庆交通大学", "陕西科技大学", "兰州交通大学", "湖南科技大学", "辽宁师范大学", "湖南农业大学", "南方科技大学", "沈阳工业大学", "安徽工业大学", "河北农业大学", "延边大学", "广州中医药大学", "河南工业大学"}
	englishName = []string{"Wioletta", "maartje", "Germain", "ryun", "ayden", "kaiden", "Dorine", "mouna", "Brook", "hershel", "Lygia", "Landly", "Rosita", "Benita", "jelle", "Dollie", "iole", "Sydelle", "Iouliana", "viviyn",
		"Agam", "Green", "Clyde", "dominik", "Daton", "Varghese", "Jalynn", "Manaia", "Tayte", "Punites", "delphina", "pakwa", "Wilmar", "Kamron", "natalia", "Alion", "Daphne", "kiarr", "jezmine", "Tresa", "Bila", "Auriar", "Tziyon",
		"ramelle", "marius", "Giuliana", "jaylon", "Maddog", "Heartha", "Nancee", "Tayler", "Shiralee", "Josip", "kae", "ilu", "manel", "Dragan", "brande", "Viktorija", "Boyd", "Tamarind", "tottie", "Wincent", "Roslyn", "maritssa",
		"Arisha", "Mikkel", "Mchumba", "silvester", "Sarika", "Ingram", "Sharpay", "elmira", "Foster", "Anita", "Maryon", "ivet", "Mikan", "Trevan", "daxton"}
)

func randomName(builder *strings.Builder) *strings.Builder {
	builder.WriteString(firstName[rand.Intn(len(firstName))])
	builder.WriteString(secondName[rand.Intn(len(secondName))])
	thirdRate := rand.Intn(10) //[0,9]
	if thirdRate > 2 {
		builder.WriteString(secondName[rand.Intn(len(secondName))])
	}
	return builder
}

func GenerateName(nameType int) string {
	builder := new(strings.Builder)
	switch nameType {
	case 1: // 默认
		builder.WriteString("Vlander")
		builder.WriteString(strconv.Itoa(rand.Intn(9000) + 1000))
	case 2: // 真实名字
		randomName(builder)
	case 3: // 校园招聘
		builder.WriteString(collegeName[rand.Intn(len(collegeName))])
		builder.WriteString("-")
		randomName(builder)
	case 4: // 英文名字
		builder.WriteString(englishName[rand.Intn(len(englishName))])
	}

	//alphabetLen := len(alphabet)
	////三种类型 1a1a1a 11a11a 1aa1aa
	//randomType := rand.Intn(3)
	//switch randomType {
	//case 0:
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//case 1:
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//case 2:
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(strconv.Itoa(rand.Intn(10)))
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//	builder.WriteString(alphabet[rand.Intn(alphabetLen)])
	//}

	return builder.String()
}
