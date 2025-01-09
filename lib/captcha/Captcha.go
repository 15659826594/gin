package captcha

import (
	"bytes"
	"fmt"
	"gin/utils"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var __DIR__ string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	__DIR__ = utils.Pathinfo(filename, 1)["dirname"]
}

type captchaOption struct {
	seKey string
	// 验证码加密密钥
	codeSet string
	// 验证码字符集合
	expire int
	// 验证码过期时间（s）
	useZh bool
	// 使用中文验证码
	zhSet string
	// 中文验证码字符串
	useImgBg bool
	// 使用背景图片
	fontSize int
	// 验证码字体大小(px)
	useCurve bool
	// 是否画混淆曲线
	useNoise bool
	// 是否添加杂点
	imageH int
	// 验证码图片高度
	imageW int
	// 验证码图片宽度
	length int
	// 验证码位数
	fontttf string
	// 验证码字体，不设置随机获取
	bg [3]int
	// 背景颜色
	reset bool
	// 验证成功后是否重置
	useArithmetic bool //是否使用算术验证码
}

type Captcha struct {
	config captchaOption
	im     *image.RGBA
	color  image.Image
}

/*Entry
 * 输出验证码并把验证码的值保存的session中
 * 验证码保存到session的格式为： array('verify_code' => '验证码值', 'verify_time' => '验证码创建时间');
 * @access public
 * @param string $id 要生成验证码的标识
 * @return \think\Response
 */
func (that *Captcha) Entry(id string) []byte {
	// 图片宽(px)
	if that.config.imageW == 0 {
		that.config.imageW = int(float64(that.config.length)*float64(that.config.fontSize)*1.5 + float64(that.config.length)*float64(that.config.fontSize)/2)
	}
	// 图片高(px)
	if that.config.imageH == 0 {
		that.config.imageH = int(float64(that.config.fontSize) * 2.5)
	}
	// 建立一幅 that.config.imageW x that.config.imageH 的图像
	that.im = image.NewRGBA(image.Rect(0, 0, that.config.imageW, that.config.imageH))
	bgColor := imagecolorallocate(that.config.bg[0], that.config.bg[1], that.config.bg[2])
	// 设置背景
	draw.Draw(that.im, that.im.Bounds(), bgColor, image.Point{}, draw.Src)

	// 验证码字体随机颜色
	that.color = imagecolorallocate(utils.MtRand(1, 150), utils.MtRand(1, 150), utils.MtRand(1, 150))

	if that.config.useImgBg {
		that.Background()
	}

	if that.config.useNoise {
		// 绘杂点
		that.WriteNoise()
	}

	if that.config.useCurve {
		// 绘干扰线
		that.WriteCurve()
	}

	var ttf string
	if that.config.useZh {
		ttf = "zhttfs"
	} else {
		ttf = "ttfs"
	}

	ttfPath := fmt.Sprintf("%s/assets/%s/", __DIR__, ttf)

	// 绘验证码
	var code []any
	var codeNx float64
	var codeStr string
	if that.config.useArithmetic {
		that.config.fontttf = ttfPath + "6.ttf"
		fontBytes, _ := os.ReadFile(that.config.fontttf)
		freeFont, _ := freetype.ParseFont(fontBytes)

		code = []any{strconv.Itoa(utils.MtRand(1, 9)), "+", strconv.Itoa(utils.MtRand(1, 9))}
		for _, char := range code {
			codeNx += float64(that.config.fontSize) * 1.5
			imagettftext(that.im, float64(that.config.fontSize), 0, int(codeNx), int(float64(that.config.fontSize)*1.6), that.color, freeFont, char.(string))
		}
		codeStr = that.Authcode(code[0].(string) + code[2].(string))
	} else {
		// 验证码使用随机字体
		if that.config.fontttf == "" {
			var ttfs []any
			_ = filepath.Walk(ttfPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					ttfs = append(ttfs, utils.Pathinfo(path, 2)["basename"])
				}
				return nil
			})
			that.config.fontttf = utils.ArrayRand(ttfs)[0].(string)
		}
		that.config.fontttf = ttfPath + that.config.fontttf

		// 读取字体
		fontBytes, _ := os.ReadFile(that.config.fontttf)
		freeFont, _ := freetype.ParseFont(fontBytes)

		var codeArr []string
		if that.config.useZh {
			// 中文验证码
			for i := 0; i < that.config.length; i++ {
				index := utils.MtRand(0, utils.MbStrlen(that.config.zhSet)-1)
				codei := string([]rune(that.config.zhSet)[index : index+1])
				imagettftext(that.im, float64(that.config.fontSize), float64(utils.MtRand(-40, 40)), int(float64(that.config.fontSize)*(float64(i)+1)*1.5), that.config.fontSize+utils.MtRand(10, 20), that.color, freeFont, codei)
			}
		} else {
			for i, strlen := 0, len(that.config.codeSet)-1; i < that.config.length; i++ {
				codei := string(that.config.codeSet[utils.MtRand(0, strlen)])
				code = append(code, codei)
				codeNx += float64(utils.MtRand(int(float64(that.config.fontSize)*1.2), int(float64(that.config.fontSize)*1.6)))
				imagettftext(that.im, float64(that.config.fontSize), float64(utils.MtRand(-40, 40)), int(codeNx), int(float64(that.config.fontSize)*1.6), that.color, freeFont, codei)
			}
		}
		codeStr = that.Authcode(strings.ToUpper(utils.Implode("", codeArr)))
	}

	// 保存验证码
	key := that.Authcode(that.config.seKey)

	secode := map[string]any{
		"verify_code": codeStr,
		"verify_time": time.Now().Unix(),
	}

	fmt.Println(key, secode)

	var buf bytes.Buffer
	_ = png.Encode(&buf, that.im)

	return buf.Bytes()
}

/*WriteNoise
 * 画杂点
 * 往图片上写不同颜色的字母或数字
 */
func (that *Captcha) WriteNoise() {
	codeSet := "2345678abcdefhijkmnpqrstuvwxyz"

	for i := 0; i < 10; i++ {
		noiseColor := imagecolorallocate(utils.MtRand(150, 225), utils.MtRand(150, 225), utils.MtRand(150, 225))
		for j := 0; j < 5; j++ {
			imagestring(that.im, 14, utils.MtRand(-10, that.config.imageW), utils.MtRand(-10, that.config.imageH), string(codeSet[utils.MtRand(0, len(codeSet)-1)]), noiseColor)
		}
	}
}

/*WriteCurve
 * 画一条由两条连在一起构成的随机正弦函数曲线作干扰线(你可以改成更帅的曲线函数)
 *
 *      高中的数学公式咋都忘了涅，写出来
 *        正弦型函数解析式：y=Asin(ωx+φ)+b
 *      各常数值对函数图像的影响：
 *        A：决定峰值（即纵向拉伸压缩的倍数）
 *        b：表示波形在Y轴的位置关系或纵向移动距离（上加下减）
 *        φ：决定波形与X轴位置关系或横向移动距离（左加右减）
 *        ω：决定周期（最小正周期T=2π/∣ω∣）
 *
 */
func (that *Captcha) WriteCurve() {
	var py float64
	var px float64
	// 曲线前部分
	A := utils.MtRand(1, int(float64(that.config.imageH)/2))                                       // 振幅
	b := utils.MtRand(-(int(float64(that.config.imageH) / 4)), int(float64(that.config.imageH)/4)) // Y轴方向偏移量
	f := utils.MtRand(-(int(float64(that.config.imageH) / 4)), int(float64(that.config.imageH)/4)) // X轴方向偏移量
	T := utils.MtRand(that.config.imageH, int(float64(that.config.imageW)/2))                      // 周期

	w := math.Pi * 2 / float64(T)

	var px1 float64                                                                                        // 曲线横坐标起始位置
	px2 := float64(utils.MtRand(int(float64(that.config.imageW)/2), int(float64(that.config.imageW)*0.8))) // 曲线横坐标结束位置

	for px = px1; px <= px2; px = px + 1 {
		if w != 0 {
			py = float64(A)*math.Sin(w*px+float64(f)) + float64(b) + float64(that.config.imageH)/2 // y = Asin(ωx+φ) + b
			i := int(float64(that.config.fontSize) / 5)
			for i > 0 {
				imagesetpixel(that.im, int(px+float64(i)), int(py+float64(i)), that.color) // 这里(while)循环画像素点比imagettftext和imagestring用字体大小一次画出（不用这while循环）性能要好很多
				i--
			}
		}
	}

	// 曲线后部分
	A = utils.MtRand(1, that.config.imageH/2)                         // 振幅
	f = utils.MtRand(-(that.config.imageH / 4), that.config.imageH/4) // X轴方向偏移量
	T = utils.MtRand(that.config.imageH, that.config.imageW*2)        // 周期
	w = math.Pi * 2 / float64(T)
	b = int(py - float64(A)*math.Sin(w*px+float64(f)) - float64(that.config.imageH)/2)
	px1 = px2
	px2 = float64(that.config.imageW)

	for px = px1; px <= px2; px = px + 1 {
		if w != 0 {
			py = float64(A)*math.Sin(w*px+float64(f)) + float64(b) + float64(that.config.imageH)/2 // y = Asin(ωx+φ) + b
			i := int(float64(that.config.fontSize) / 5)
			for i > 0 {
				imagesetpixel(that.im, int(px+float64(i)), int(py+float64(i)), that.color)
				i--
			}
		}
	}
}

/*Background
 * 绘制背景图片
 * 注：如果验证码输出图片比较大，将占用比较多的系统资源
 */
func (that *Captcha) Background() {
	path := __DIR__ + "/assets/bgs/"

	var bgs []any
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			bgs = append(bgs, utils.Pathinfo(path, 2)["basename"])
		}
		return nil
	})

	gb := path + utils.ArrayRand(bgs)[0].(string)

	bgFile, _ := os.Open(gb)

	defer func(bgFile *os.File) {
		_ = bgFile.Close()
	}(bgFile)

	bgImage, _, _ := image.Decode(bgFile)

	//缩放图像
	width, height := getimagesize(that.im)
	srcImage := Resize(uint(width), uint(height), bgImage, Lanczos3)

	draw.Draw(that.im, that.im.Bounds(), srcImage, image.Point{}, draw.Src)
}

//Authcode
/* 加密验证码 */
func (that *Captcha) Authcode(str string) string {
	key := utils.Substr(utils.Md5(that.config.seKey), 5, 8)
	str1 := utils.Substr(utils.Md5(str), 8, 10)
	return utils.Md5(key + str1)
}

func NewCaptcha(config map[string]any) *Captcha {
	captchaOption := captchaOption{
		seKey:         "ThinkPHP.CN",
		codeSet:       "2345678abcdefhijkmnpqrstuvwxyzABCDEFGHJKLMNPQRTUVWXY",
		expire:        1800,
		useZh:         false,
		zhSet:         "们以我到他会作时要动国产的一是工就年阶义发成部民可出能方进在了不和有大这主中人上为来分生对于学下级地个用同行面说种过命度革而多子后自社加小机也经力线本电高量长党得实家定深法表着水理化争现所二起政三好十战无农使性前等反体合斗路图把结第里正新开论之物从当两些还天资事队批点育重其思与间内去因件日利相由压员气业代全组数果期导平各基或月毛然如应形想制心样干都向变关问比展那它最及外没看治提五解系林者米群头意只明四道马认次文通但条较克又公孔领军流入接席位情运器并飞原油放立题质指建区验活众很教决特此常石强极土少已根共直团统式转别造切九你取西持总料连任志观调七么山程百报更见必真保热委手改管处己将修支识病象几先老光专什六型具示复安带每东增则完风回南广劳轮科北打积车计给节做务被整联步类集号列温装即毫知轴研单色坚据速防史拉世设达尔场织历花受求传口断况采精金界品判参层止边清至万确究书术状厂须离再目海交权且儿青才证低越际八试规斯近注办布门铁需走议县兵固除般引齿千胜细影济白格效置推空配刀叶率述今选养德话查差半敌始片施响收华觉备名红续均药标记难存测士身紧液派准斤角降维板许破述技消底床田势端感往神便贺村构照容非搞亚磨族火段算适讲按值美态黄易彪服早班麦削信排台声该击素张密害侯草何树肥继右属市严径螺检左页抗苏显苦英快称坏移约巴材省黑武培著河帝仅针怎植京助升王眼她抓含苗副杂普谈围食射源例致酸旧却充足短划剂宣环落首尺波承粉践府鱼随考刻靠够满夫失包住促枝局菌杆周护岩师举曲春元超负砂封换太模贫减阳扬江析亩木言球朝医校古呢稻宋听唯输滑站另卫字鼓刚写刘微略范供阿块某功套友限项余倒卷创律雨让骨远帮初皮播优占死毒圈伟季训控激找叫云互跟裂粮粒母练塞钢顶策双留误础吸阻故寸盾晚丝女散焊功株亲院冷彻弹错散商视艺灭版烈零室轻血倍缺厘泵察绝富城冲喷壤简否柱李望盘磁雄似困巩益洲脱投送奴侧润盖挥距触星松送获兴独官混纪依未突架宽冬章湿偏纹吃执阀矿寨责熟稳夺硬价努翻奇甲预职评读背协损棉侵灰虽矛厚罗泥辟告卵箱掌氧恩爱停曾溶营终纲孟钱待尽俄缩沙退陈讨奋械载胞幼哪剥迫旋征槽倒握担仍呀鲜吧卡粗介钻逐弱脚怕盐末阴丰雾冠丙街莱贝辐肠付吉渗瑞惊顿挤秒悬姆烂森糖圣凹陶词迟蚕亿矩康遵牧遭幅园腔订香肉弟屋敏恢忘编印蜂急拿扩伤飞露核缘游振操央伍域甚迅辉异序免纸夜乡久隶缸夹念兰映沟乙吗儒杀汽磷艰晶插埃燃欢铁补咱芽永瓦倾阵碳演威附牙芽永瓦斜灌欧献顺猪洋腐请透司危括脉宜笑若尾束壮暴企菜穗楚汉愈绿拖牛份染既秋遍锻玉夏疗尖殖井费州访吹荣铜沿替滚客召旱悟刺脑措贯藏敢令隙炉壳硫煤迎铸粘探临薄旬善福纵择礼愿伏残雷延烟句纯渐耕跑泽慢栽鲁赤繁境潮横掉锥希池败船假亮谓托伙哲怀割摆贡呈劲财仪沉炼麻罪祖息车穿货销齐鼠抽画饲龙库守筑房歌寒喜哥洗蚀废纳腹乎录镜妇恶脂庄擦险赞钟摇典柄辩竹谷卖乱虚桥奥伯赶垂途额壁网截野遗静谋弄挂课镇妄盛耐援扎虑键归符庆聚绕摩忙舞遇索顾胶羊湖钉仁音迹碎伸灯避泛亡答勇频皇柳哈揭甘诺概宪浓岛袭谁洪谢炮浇斑讯懂灵蛋闭孩释乳巨徒私银伊景坦累匀霉杜乐勒隔弯绩招绍胡呼痛峰零柴簧午跳居尚丁秦稍追梁折耗碱殊岗挖氏刃剧堆赫荷胸衡勤膜篇登驻案刊秧缓凸役剪川雪链渔啦脸户洛孢勃盟买杨宗焦赛旗滤硅炭股坐蒸凝竟陷枪黎救冒暗洞犯筒您宋弧爆谬涂味津臂障褐陆啊健尊豆拔莫抵桑坡缝警挑污冰柬嘴啥饭塑寄赵喊垫丹渡耳刨虎笔稀昆浪萨茶滴浅拥穴覆伦娘吨浸袖珠雌妈紫戏塔锤震岁貌洁剖牢锋疑霸闪埔猛诉刷狠忽灾闹乔唐漏闻沈熔氯荒茎男凡抢像浆旁玻亦忠唱蒙予纷捕锁尤乘乌智淡允叛畜俘摸锈扫毕璃宝芯爷鉴秘净蒋钙肩腾枯抛轨堂拌爸循诱祝励肯酒绳穷塘燥泡袋朗喂铝软渠颗惯贸粪综墙趋彼届墨碍启逆卸航衣孙龄岭骗休借",
		useImgBg:      false,
		fontSize:      25,
		useCurve:      true,
		useNoise:      true,
		imageH:        0,
		imageW:        0,
		length:        5,
		fontttf:       "",
		bg:            [3]int{243, 251, 254},
		reset:         true,
		useArithmetic: false,
	}
	for key, value := range config {
		switch key {
		case "seKey":
			captchaOption.seKey = value.(string)
		case "codeSet":
			captchaOption.codeSet = value.(string)
		case "expire":
			captchaOption.expire = value.(int)
		case "useZh":
			captchaOption.useZh = value.(bool)
		case "zhSet":
			captchaOption.zhSet = value.(string)
		case "useImgBg":
			captchaOption.useImgBg = value.(bool)
		case "fontSize":
			captchaOption.fontSize = value.(int)
		case "useCurve":
			captchaOption.useCurve = value.(bool)
		case "useNoise":
			captchaOption.useNoise = value.(bool)
		case "imageH":
			captchaOption.imageH = value.(int)
		case "imageW":
			captchaOption.imageW = value.(int)
		case "length":
			captchaOption.length = value.(int)
		case "fontttf":
			captchaOption.fontttf = value.(string)
		case "bg":
			captchaOption.bg = value.([3]int)
		case "reset":
			captchaOption.reset = value.(bool)
		case "useArithmetic":
			captchaOption.useArithmetic = value.(bool)
		}
	}

	captcha := &Captcha{
		config: captchaOption,
	}
	return captcha
}
