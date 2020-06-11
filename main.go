package main

import (
	//缓存IO

	"fmt"
	"io/ioutil" //io 工具包
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type RunPath struct {
	devicePath1 string
	devicePath2 string
	devicePath3 string
	devicePath4 string
	devicePath5 string
	savePath    string
}

type InitPara struct {
	path       RunPath
	deviceType uint8
	logLevel   uint8
	timeOut    uint32
}

func max(num1 int16, num2 int16) int16 {
	/* 声明局部变量 */
	var result int16

	if num1 > num2 {
		result = num1
	} else {
		result = num2
	}
	return result
}

func listFile(myfolder string) {
	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			listFile(myfolder + "\\" + file.Name())
		} else {
			fmt.Println(myfolder + "\\" + file.Name())
		}
	}
}

func checkEmptyDir(myDir string) bool {
	var ret bool = false

	dir, _ := ioutil.ReadDir(myDir)
	if len(dir) == 0 {
		fmt.Println(myDir + " is empty dir!")
		ret = true
	} else {
		fmt.Println(myDir + " is not empty dir!")
	}

	return ret
}

func testDelFile() {
	originalFile, err := os.Open("zyh")
	if err != nil {
		log.Fatal(err)
	}
	originalFile.Close()

	err = os.Remove("zyh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("all is ok")
}

func logTest() {
	// 定义一个文件
	fileName := "ll.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	debugLog := log.New(logFile, "[Debug]", log.LstdFlags)
	debugLog.Println("A debug message here")
	//配置一个日志格式的前缀
	debugLog.SetPrefix("[Info]")
	debugLog.Println("A Info Message here ")
	//配置log的Flag参数
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	debugLog.Println("A different prefix")
}

//go 自带log 库配置
func logInit1() {
	file := "./" + "message" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		LOG.Error(err)
	}
	log.SetOutput(logFile)   // 将文件设置为log输出的文件
	log.SetPrefix("[DFMTP]") //设置log前缀
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}


func mappingTimeOut(timeOut uint32)  uint32{
	var time uint32 = 10000//默认1天过期

	if timeOut <= 168 && timeOut > 0{
		time = (timeOut % 24)*100 + (timeOut/24)*10000
	}else{
		LOG.Error("TimeOut error, set deafult value to run App!")
	}
	return time
}

/**
* 从程序目录中的config.txt获取配置信息
 */
func getRunParaInit(appRunPara *InitPara) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		LOG.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// 典型读取操作，默认分区可以使用空字符串表示
	appRunPara.deviceType = uint8(cfg.Section("paths").Key("number").MustUint(0))
	appRunPara.timeOut = uint32(cfg.Section("paths").Key("timeout").MustUint(0))
	appRunPara.logLevel = uint8(cfg.Section("paths").Key("loglevel").MustUint(0))
	//SM1281的目录地址
	appRunPara.path.devicePath1 = cfg.Section("paths").Key("path1").String()
	appRunPara.path.devicePath2 = cfg.Section("paths").Key("path2").String()
	appRunPara.path.devicePath3 = cfg.Section("paths").Key("path3").String()
	appRunPara.path.devicePath4 = cfg.Section("paths").Key("path4").String()
	appRunPara.path.devicePath5 = cfg.Section("paths").Key("path5").String()
	//转存sm1281到本地的目录
	appRunPara.path.savePath = cfg.Section("paths").Key("runpath").String()

}

/**
* 打印程序启动后的运行信息
 */
func logInitInfo(appRunPara InitPara) {
	LOG.Infof("Welcome DFTMP Demo!!!!")
	LOG.Infof("device Type:%d", appRunPara.deviceType)
	LOG.Infof("time Out:%d", appRunPara.timeOut)
	LOG.Infof("Log level:%d", appRunPara.logLevel)
	LOG.Infof("devicePath1:%s", appRunPara.path.devicePath1)
	LOG.Infof("devicePath2:%s", appRunPara.path.devicePath2)
	LOG.Infof("devicePath3:%s", appRunPara.path.devicePath3)
	LOG.Infof("devicePath4:%s", appRunPara.path.devicePath4)
	LOG.Infof("devicePath5:%s", appRunPara.path.devicePath5)
	LOG.Infof("savePath:%s", appRunPara.path.savePath)
}

/**
* 检查配置文件是否有错误
* 无问题，true
* 有问题，false
 */
func chkParaValid(appRunPara InitPara) bool {
	var bRet bool = true

	if appRunPara.deviceType > 5 && appRunPara.deviceType < 1 {
		bRet = false
	}

	if appRunPara.logLevel > 5 && appRunPara.logLevel < 0 {
		bRet = false
	}

	if strings.Compare(appRunPara.path.savePath, "") == 0 {
		bRet = false
	}

	return bRet
}

/**
* 启动主程序中多线程任务，根据webdav server数量来创建不同的线程来处理任务
 */
func mainAppRun(appRunPara InitPara) {
	var sExpTime1, sExpTime2, sExpTime3, sExpTime4, sExpTime5 ExpTime

	uTimeOut :=  mappingTimeOut(uint32(appRunPara.timeOut))

	bRet := chkParaValid(appRunPara)

	//配置文件错误，退出程序
	if bRet == false {
		LOG.Error("Configuration error!!!")
		return
	}

	//devicetype为sm1281数量，有几个sm1281处理线程
	switch appRunPara.deviceType {
	case 1:
		{
			_, sExpTime1 = runParaRead("runPara1")
			procSM1281RawData(appRunPara.path.devicePath1+"\\", 1, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime1)
		}
	case 2:
		{
			_, sExpTime1 = runParaRead("runPara1")
			_, sExpTime2 = runParaRead("runPara2")
			go procSM1281RawData(appRunPara.path.devicePath1+"\\", 1, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime1)
			procSM1281RawData(appRunPara.path.devicePath2+"\\", 2, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime2)
		}
	case 3:
		{
			_, sExpTime1 = runParaRead("runPara1")
			_, sExpTime2 = runParaRead("runPara2")
			_, sExpTime3 = runParaRead("runPara3")
			go procSM1281RawData(appRunPara.path.devicePath1+"\\", 1, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime1)
			go procSM1281RawData(appRunPara.path.devicePath2+"\\", 2, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime2)
			procSM1281RawData(appRunPara.path.devicePath3+"\\", 3, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime3)
		}
	case 4:
		{
			_, sExpTime1 = runParaRead("runPara1")
			_, sExpTime2 = runParaRead("runPara2")
			_, sExpTime3 = runParaRead("runPara3")
			_, sExpTime4 = runParaRead("runPara4")
			go procSM1281RawData(appRunPara.path.devicePath1+"\\", 1, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime1)
			go procSM1281RawData(appRunPara.path.devicePath2+"\\", 2, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime2)
			go procSM1281RawData(appRunPara.path.devicePath3+"\\", 3, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime3)
			procSM1281RawData(appRunPara.path.devicePath4+"\\", 4, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime4)
		}
	case 5:
		{
			_, sExpTime1 = runParaRead("runPara1")
			_, sExpTime2 = runParaRead("runPara2")
			_, sExpTime3 = runParaRead("runPara3")
			_, sExpTime4 = runParaRead("runPara4")
			_, sExpTime5 = runParaRead("runPara5")
			go procSM1281RawData(appRunPara.path.devicePath5+"\\", 5, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime5)
			go procSM1281RawData(appRunPara.path.devicePath4+"\\", 4, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime4)
			go procSM1281RawData(appRunPara.path.devicePath3+"\\", 3, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime3)
			go procSM1281RawData(appRunPara.path.devicePath2+"\\", 2, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime2)
			procSM1281RawData(appRunPara.path.devicePath1+"\\", 1, uTimeOut, appRunPara.path.savePath+"\\", &sExpTime1)
		}
	default:
		{
			LOG.Error("Device type error, App starts failure!!!")
		}
	}
}


/**
* 主程序，包括log初始化，运行参数初始化，程序启动打印，主程序流程
 */
func main() {

	//运行变量初始化
	var appRunPara InitPara

	//初始化全局运行参数
	getRunParaInit(&appRunPara)

	//log配置初始化
	logInit(appRunPara.logLevel)
	//打印初始时间
	logInitInfo(appRunPara)
	//主流程程序
	mainAppRun(appRunPara)

}
