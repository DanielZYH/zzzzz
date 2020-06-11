package main

import (
	//缓存IO

	"bufio"
	"io"
	"io/ioutil" //io 工具包
	"os"
	"strconv" //Package strconv implements conversions to and from string representations of basic data types.
	"strings"
	"time"
)

const maxDirSize uint16 = 300
const maxFileNameSize uint16 = 45
const WAITBUFFTIME uint8 = 10
const testSrcFileName string = "20191202_075807_demo_VIB1_VIB2.wav"
const LASTIMEMAXSIZE uint8 = 24 * 7

type ExpTime struct {
	idxCurr    uint8 //新数据存储idx
	idxLast    uint8 //比较过期时间idx
	length     uint8 //数组长度
	expTimeAry [LASTIMEMAXSIZE]string
}

//const testDestFileName string = "20191202_075807_demo_VIB1_VIB2_BACK.wav"
/**
 * 任务执行时相关变量保存在本地文件，程序重启后读取后再使用，保证上次的运行状态
 */
func runParaWrite(expTime *ExpTime, devNo uint8) {
	//var tmpData ExpTime
	var indx uint8 = 0

	fileObj, err := os.OpenFile("runPara"+strconv.Itoa(int(devNo)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		LOG.Errorf("Failed to open the file---[%s]", err.Error())
		//os.Exit(2)
	}
	//defer fileObj.Close()

	tmpDataInx := strconv.Itoa(int(expTime.idxLast)) + "\n"
	if _, err = fileObj.WriteString(tmpDataInx); err == nil {
		//log.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", tmpDataInx)
	}

	tmpDataInx = strconv.Itoa(int(expTime.idxCurr)) + "\n"
	if _, err = fileObj.WriteString(tmpDataInx); err == nil {
		//log.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", tmpDataInx)
	}

	tmpDataInx = strconv.Itoa(int(expTime.length)) + "\n"
	if _, err = fileObj.WriteString(tmpDataInx); err == nil {
		//log.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", tmpDataInx)
	}

	for ; indx < expTime.length; indx++ {
		flag := strings.Contains(expTime.expTimeAry[indx], "\n")
		if flag == false {
			expTime.expTimeAry[indx] = expTime.expTimeAry[indx] + "\n"
		}

		if _, err = fileObj.WriteString(expTime.expTimeAry[indx]); err == nil {
			//	log.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", expTime.expTimeAry[indx])
		}
	}
	fileObj.Close()
	logZyh := strconv.Itoa(int(devNo))
	LOG.Infof("Device[%d]:"+"Running parameters saved!!!!", logZyh)
}

/**
 * 任务执行时相关变量保存在本地文件，程序重启后读取后再使用，保证上次的运行状态
 */
func runParaRead(fileName string) (error, ExpTime) {
	var localExpTime ExpTime

	f, err := os.Open(fileName)
	if err != nil {
		return err, localExpTime
	}

	buf := bufio.NewReader(f)
	//the index of Last pos
	line, err := buf.ReadString('\n')
	line = strings.TrimSpace(line)
	tmpLine, _ := strconv.Atoi(line)
	localExpTime.idxLast = uint8(tmpLine)

	//the index of current pos
	line, err = buf.ReadString('\n')
	line = strings.TrimSpace(line)
	tmpLine, _ = strconv.Atoi(line)
	localExpTime.idxCurr = uint8(tmpLine)

	//the index of current pos
	line, err = buf.ReadString('\n')
	line = strings.TrimSpace(line)
	tmpLine, _ = strconv.Atoi(line)
	localExpTime.length = uint8(tmpLine)

	var aryIndx uint8
	for aryIndx = 0; aryIndx < localExpTime.length; aryIndx++ {
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		localExpTime.expTimeAry[aryIndx] = line

		if err != nil {
			if err == io.EOF {
				return nil, localExpTime
			}
			return err, localExpTime
		}

	}
	return nil, localExpTime
}

/**
 * 下载指定的文件到本地
 */
func procRawDataDL(srcFileName, fileDir, storeDir string, severNo uint8) {
	originalFile, err := os.Open(fileDir + srcFileName)
	if err != nil {
		LOG.Error(err)
	}
	//defer originalFile.Close()

	vibFlagIndx := strings.Index(srcFileName, "VIB")
	destDir := storeDir + srcFileName[0:8] + "\\" + srcFileName[9:11] + "\\" + srcFileName[11:13] + "\\"
	destFileName := destDir + srcFileName[0:vibFlagIndx] + "_" + strconv.Itoa(int(severNo)) + ".wav"

	os.MkdirAll(destDir, 0777)

	//创建新文件
	newFile, err := os.Create(destFileName)
	if err != nil {
		LOG.Error(err)
	}
	//defer newFile.Close()

	//文件复制
	bytes, err := io.Copy(newFile, originalFile)
	if err != nil {
		LOG.Error(err)
	}
	LOG.Infof("Device%d:copy file from[%s] to [%s]---[%dbytes]", severNo, fileDir+srcFileName, destFileName, bytes)

	err = newFile.Sync()
	if err != nil {
		LOG.Error(err)
	}

	originalFile.Close()
	newFile.Close()

	//删除SM1281原文件
	err = os.Remove(fileDir + srcFileName)
	if err != nil {
		LOG.Error(err)
	}

}

func getDirName(targetDirName string) (string, string, string) {
	var strRet string = ""
	var strTopDir string = ""
	var strSecDir string = ""

	if targetDirName != "" {
		timeAIndx := strings.Index(targetDirName, "_")

		strTopDir = targetDirName[:timeAIndx]
		strSecDir = targetDirName[timeAIndx+1 : timeAIndx+3]
		tmpStrA := strTopDir + "\\" + strSecDir
		LOG.Infof("get dir name: [%s]\n", tmpStrA)
		strRet = tmpStrA

	} else {
		LOG.Warningf("get dir name: [%s]\n", targetDirName)
	}

	return strRet, strTopDir, strSecDir
}

/**
 * 删除过期的目录，以小时为单元
 * 如果当天目录已经删空，将当天目录也删除
 */
func delOutDateDir(pExpTime *ExpTime, saveDir, currTime string, timeOut uint32, devNo uint8) {

	delDirName, delTopDir, delSecDir := getDirName(pExpTime.expTimeAry[pExpTime.idxLast])

	//删除过期的目录
	delDir(saveDir + delDirName)
	//23时目录删除后，将当天整个目录删除
	if strings.Compare(delSecDir, "23") == 0 {
		delDir(saveDir + delTopDir)
	} else {

	}
	//*lastTime = currTime
	LOG.Debugf("Device%d:Delete directory[%s] currTime[%s], Expire time arry size[%d] curr_pos[%d] last_pos[%d]", devNo, (saveDir + delDirName), currTime, pExpTime.length, pExpTime.idxCurr, pExpTime.idxLast)
}

/**
 * 检查本地raw data目录是否超过7天，删除数据以小时数单位
 * 如果是，删除第一天的对应的小时的数据
 */
func procRawDirOverSize(webDavDir, saveDir string, pExpTime *ExpTime, currTime string, timeOut uint32, devNo uint8) {
	//strA := "20191202_075807_demo_VIB1_VIB2.wav"
	//strB := "20191902_015807_demo_VIB1_VIB2.wav"
	var flag bool = false

	var lastTime string = pExpTime.expTimeAry[pExpTime.idxCurr] //获取数据最后一个存储过期日期
	var cmpTime string = pExpTime.expTimeAry[pExpTime.idxLast]  //获取过期日期中第一个

	timeAIndx := strings.Index(lastTime, "_") //过期时间数组当前下标
	tmpStrAIns := lastTime[:timeAIndx] + lastTime[timeAIndx+1:timeAIndx+3]
	//tmpStrA := lastTime[:timeAIndx] + lastTime[timeAIndx+1:timeAIndx+5]

	timeBIndx := strings.Index(currTime, "_") //当前需要比较的时间
	tmpStrBIns := currTime[:timeBIndx] + currTime[timeBIndx+1:timeBIndx+3]
	tmpStrB := currTime[:timeBIndx] + currTime[timeBIndx+1:timeBIndx+5]

	timeCIndx := strings.Index(cmpTime, "_") //需要比较时间下标
	/////tmpStrCIns := cmpTime[:timeCIndx] + cmpTime[timeCIndx+1:timeCIndx+3]
	tmpStrC := cmpTime[:timeCIndx] + cmpTime[timeCIndx+1:timeCIndx+5]

	//比较年月日时，插入新的元素到过期时间数组中
	if strings.Compare(tmpStrAIns, tmpStrBIns) != 0 {

		pExpTime.idxCurr = (pExpTime.idxCurr + 1) % LASTIMEMAXSIZE
		pExpTime.expTimeAry[pExpTime.idxCurr] = currTime
		//过期时间数组是循环数组
		if pExpTime.length < LASTIMEMAXSIZE {
			pExpTime.length++
		}
		flag = true
		//runParaWrite(pExpTime, devNo)
		LOG.Infof("Device%d:Insert a new expire time[%d]=[%s], index[%d]\n", devNo, pExpTime.idxCurr, pExpTime.expTimeAry[pExpTime.idxCurr], pExpTime.length)
	}

	timeGapC, _ := strconv.Atoi(tmpStrC) //过期比较时间
	timeGapB, _ := strconv.Atoi(tmpStrB) //当前时间
	timeGap := uint32(timeGapB - timeGapC)

	/**
	 * 检查本地raw data目录是否超过7天，删除数据以小时数单位
	 * 如果是，删除第一天的对应的小时的数据
	 */
	if timeGap >= timeOut {

		delOutDateDir(pExpTime, saveDir, currTime, timeOut, devNo)
		pExpTime.idxLast = (pExpTime.idxLast + 1) % LASTIMEMAXSIZE

		LOG.Infof("Device%d:Time expired[%d][%s][%s][%v]", devNo, pExpTime.idxLast, tmpStrC, tmpStrB, timeGap)
		//保存修改的数据
		//runParaWrite(pExpTime, devNo)
		flag = true
	} else {
		LOG.Debugf("Device%d:No expired[%d] [%s] [%s] [%v]", devNo, pExpTime.idxLast, tmpStrC, tmpStrB, timeGap)
	}

	if flag == true {
		//保存修改的数据
		runParaWrite(pExpTime, devNo)
	}
}

/**
 * 删除，指定的目录下所有内容
 */
func delDir(dir string) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		LOG.Infof("The Directory[%s] does not exist", dir)
	} else {
		LOG.Infof("The Directory[%s] has been deleted", dir)
		os.RemoveAll(dir)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/**
 * 遍历一个文件夹  返回文件夹下所有文件及文件个数
 */
func findDir(dir string, num uint8) ([]string, uint16) {
	fileinfo, err := ioutil.ReadDir(dir)
	var fileIdnx uint16 = 0
	gDirArray := make([]string, len(fileinfo))

	if err != nil {
		LOG.Errorf("%s", err.Error())
		//panic(err)
	} else {
		//var fileIdnx uint16 = 0
		//gDirArray := make([]string, len(fileinfo))

		// 遍历这个文件夹
		for _, fi := range fileinfo {

			// 重复输出制表符，模拟层级结构
			//print(strings.Repeat("\t", num))

			// 判断是不是目录
			if fi.IsDir() != true {
				//println(`文件：`, fi.Name())
				if strings.Contains(fi.Name(), ".wav") {
					gDirArray[fileIdnx] = fi.Name()
					// DirArry[fileIdnx] = fi.Name()
					fileIdnx++
				}

			}
		}
		LOG.Debugf("Found %d file from device %d", fileIdnx, num)
	}

	return gDirArray, fileIdnx
}

/**
* 下载SM1281到本地，并删除SM1281缓冲中数据
* @webDavDir SM1281服务器目录
* @devNo SM1281设备序号
* @timeOut 数据超时设置
 */
func procSM1281RawData(webDavDir string, devNo uint8, timeOut uint32, saveDir string, expTime *ExpTime) {
	var indx uint16 = 0
	var lastTimeIdx uint8 = 0
	var sExpTime *ExpTime = expTime

	for true {
		//获取当前sm1281缓冲中数据文件
		gDirArray, count := findDir(webDavDir, devNo)
		indx = 0

		if count != 0 {
			//初始化, 数组第一个为空
			if sExpTime.expTimeAry[sExpTime.idxCurr] == "" && sExpTime.idxCurr == 0 {
				//初始化为第一个文件的时间
				sExpTime.expTimeAry[sExpTime.idxCurr] = gDirArray[lastTimeIdx]
				sExpTime.idxLast = 0
				sExpTime.length++
				runParaWrite(sExpTime, devNo)
				LOG.Debugf("Device-%d:initial expire time array [content(%s), length(%d)]", devNo, sExpTime.expTimeAry[sExpTime.idxCurr], sExpTime.length)
			}

			for indx = 0; indx < count; indx++ {
				//查看是否有过期文件及目录
				procRawDirOverSize(webDavDir, saveDir, sExpTime, gDirArray[indx], timeOut, devNo)
				//println("File name：", gDirArray[indx], "webDAVDir: ", webDavDir, "devNo", devNo)
				//从sm1281下载raw data
				procRawDataDL(gDirArray[indx], webDavDir, saveDir, devNo)
			}
			//log.Printf("None file to move, wait for next round! device(%d)", devNo)
		} else {
			//如果SM1281缓冲中无新数据，等待10s
			LOG.Infof("device[%d]:No file to move, wait for next round!", devNo)
			time.Sleep(time.Second * 10)
		}

	}

}
