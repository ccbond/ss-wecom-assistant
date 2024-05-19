package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type MessageHistory struct {
	UserId   string `json:"UserId"`
	Question string `json:"Question"`
	Answer   string `json:"Answer"`
}

var mu sync.Mutex

func SaveMessage(newHistory *MessageHistory) error {
	mu.Lock()
	defer mu.Unlock()

	filePath := "/data/message_history.json"

	var histories []MessageHistory
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("读取文件时出错:", err)
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &histories)
		if err != nil {
			fmt.Println("解析JSON数据时出错:", err)
			return err
		}
	}

	histories = append(histories, *newHistory)

	// 将文件指针移动到文件开始位置
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("移动文件指针时出错:", err)
		return err
	}

	// 清空文件
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("清空文件时出错:", err)
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(histories)
	if err != nil {
		fmt.Println("写入JSON文件时出错:", err)
		return err
	}

	return nil
}
