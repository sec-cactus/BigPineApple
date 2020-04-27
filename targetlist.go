package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func getTargetList(path string) {

	pblackTargetList = &blackTargetList
	pwhiteTargetList = &whiteTargetList

	//打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(path)
	defer f.Close()

	//异常处理 以及确保函数结尾关闭文件流
	if err != nil {
		panic(err)
	}

	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//取得等号右边的value值，判断是否为空
		value, err := strconv.Atoi(strings.TrimSpace(s[index+1:]))
		if err != nil {
			fmt.Println("Read target label error", err.Error())
			continue
		}

		//interval or single ip
		index = strings.Index(key, "-")
		if index < 0 {
			//calc key/8 and key%8
			key8 := ipStringToInt64(key) / 8
			keymod8 := ipStringToInt64(key) % 8

			//single ip
			if value == 0 {
				//black list
				(*pblackTargetList)[key8] = (1 << uint8(keymod8))
				//	fmt.Println("set black list: ", ipStringToInt64(key), key8, (*pblackTargetList)[key8])
				continue
			} else if value == 1 {
				//white list
				(*pwhiteTargetList)[key8] = (1 << uint8(keymod8))
				//	fmt.Println("set white list: ", ipStringToInt64(key), key8, (*pwhiteTargetList)[key8])
				continue
			}
			continue
		} else if index > -1 {
			//interval
			if value == 0 {
				//black list
				pblackTargetList = ipIntervalToTargetList(key, pblackTargetList)
				continue
			} else if value == 1 {
				//white list
				pwhiteTargetList = ipIntervalToTargetList(key, pwhiteTargetList)
				continue
			}
			continue
		}
	}

	fmt.Println("Get targets,lens: ", len(*pwhiteTargetList), len(*pwhiteTargetList))

	return
}
