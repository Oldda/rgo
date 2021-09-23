package xls

import (
	"errors"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

type XlsxHeader *[]string

type Xlsx struct {
	FileName        string        //文件名称
	SheetName       string        //表名称
	SheetLimit      int           //每个表显示的条数
	DstDir          string        //目标路径
	SrcDir          string        //引用路径
	Suffix          string        //后缀
	Hearder         XlsxHeader    //头部信息
	Data            reflect.Value //导出时数据内容
	Cnt             int           //导入/导出条数
	ExceptFirstLine bool          //导入时是否提出第一行标题
}

func NewXlsx() *Xlsx {
	return &Xlsx{}
}

//设置文件名称
func (x *Xlsx) SetFileName(name string) *Xlsx {
	x.FileName = name
	return x
}

//设置表名称
func (x *Xlsx) SetSheetName(name string) *Xlsx {
	x.SheetName = name
	return x
}

//设置每表显示上限
func (x *Xlsx) SetSheetLimit(limit int) *Xlsx {
	x.SheetLimit = limit
	return x
}

//设置存放目录
func (x *Xlsx) SetSrcDir(src string) *Xlsx {
	x.SrcDir = src
	return x
}

//设置存放目录
func (x *Xlsx) SetDstDir(dst string) *Xlsx {
	x.DstDir = dst
	return x
}

//设置文件后缀
func (x *Xlsx) SetSuffix(suffix string) *Xlsx {
	x.Suffix = suffix
	return x
}

//设置导出内容
func (x *Xlsx) SetData(data interface{}) *Xlsx {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		log.Fatal("wrong params:data, not slice error")
	}
	x.Data = value
	return x
}

//设置是否跳过头部
func (x *Xlsx) SetExceptFirstLine(exept bool) *Xlsx {
	x.ExceptFirstLine = exept
	return x
}

//设置头部信息
func (x *Xlsx) SetHeader(header XlsxHeader) *Xlsx {
	x.Hearder = header
	return x
}

//从excel导入
func (x *Xlsx) Import(receiptSlice interface{}) error {
	if reflect.TypeOf(receiptSlice).Kind() != reflect.Ptr {
		return errors.New("needs a pointer to a slice")
	}
	//读取文件
	srcFile := x.SrcDir + x.FileName
	if srcFile == "" {
		return errors.New("file not setted")
	}
	xlFile, err := xlsx.OpenFile(srcFile)
	if err != nil {
		return err
	}
	//如何将memeber作为泛型操作-解决思路来了
	//利用reflect new初始化一个子元素 reflect.New(dest.Type().Elem()).Interface()
	sliceValue := reflect.Indirect(reflect.ValueOf(receiptSlice))
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a pointer to a slice")
	}
	sliceElementType := sliceValue.Type().Elem()
	if sliceElementType.Kind() != reflect.Struct {
		return errors.New("needs a struct")
	}
	pv := reflect.New(sliceElementType).Interface()
	//获取数据
	x.Cnt = 0
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 && x.ExceptFirstLine {
				//第一行的标题剔除
				continue
			}
			if err := row.ReadStruct(pv); err != nil {
				return err
			}
			//放到切片内
			sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(pv).Elem()))
			x.Cnt++
		}
	}
	return nil
}

/**
导出到excel
参数说明:
1. fileName 定义导出的文件名称
2. sheetName 定义表名称，默认是sheet1,2,3...
3. sheetLimit 定义表内数量上限，超过上限则分表处理，默认1000
4. dstDir 服务器存储绝对路径 如：/home/oldda/saves/
功能描述:
**/
func (x *Xlsx) Export() error {
	var file *xlsx.File
	var sheetName string
	var chanLock bool
	var limit int
	cnt := make(chan int)

	file = xlsx.NewFile()

	//默认每个表格1000条，不包含头部内容
	if x.SheetLimit == 0 {
		x.SheetLimit = 1000
	}
	//总数据量
	totalData := x.Data.Len()
	//计算总表数
	totalSheet := math.Ceil(float64(totalData) / float64(x.SheetLimit))
	//填充每个表格的内容
	for k := 1; k <= int(totalSheet); k++ {

		if x.SheetName == "" {
			sheetName = "Sheet" + strconv.Itoa(k)
		} else {
			sheetName = x.SheetName + strconv.Itoa(k)
		}

		if k*x.SheetLimit >= totalData {
			limit = totalData
		} else {
			limit = k * x.SheetLimit
		}

		if k == int(totalSheet) {
			chanLock = true
		}
		go x.wirte(file, cnt, sheetName, x.Data.Slice((k-1)*x.SheetLimit, limit), chanLock)
		x.Cnt += <-cnt
	}
	return nil
}

//写入excel文件
func (x *Xlsx) wirte(file *xlsx.File, cnt chan int, sheetName string, data reflect.Value, chanLock bool) {
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error
	var total int

	sheet, err = file.AddSheet(sheetName)
	if err != nil {
		//记录入Log
		log.Println(err)
	}

	//添加头部
	if x.Hearder != nil {
		sheet.AddRow().WriteSlice(x.Hearder, -1)
	}

	//添加行
	var rowStruct reflect.Value
	for k := 0; k < data.Len(); k++ {
		row = sheet.AddRow()

		if data.Index(k).Kind() != reflect.Ptr {
			rowStruct = data.Index(k).Addr()
		} else {
			rowStruct = data.Index(k)
		}

		if row.WriteStruct(rowStruct.Interface(), len(*x.Hearder)) > 0 {
			total++
		} else {
			//记录入Log
			log.Println(err)
		}
	}

	if x.FileName == "" {
		x.FileName = "export_" + time.Now().Format("2006-01-02")
	}

	if x.Suffix == "" {
		x.Suffix = ".xlsx"
	}

	//存储成file
	file.Save(x.DstDir + x.FileName + x.Suffix)
	cnt <- total
	if chanLock {
		close(cnt)
	}
}
