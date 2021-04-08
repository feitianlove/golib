package tailf

import (
	"github.com/hpcloud/tail"
	"os"
	"time"
)

type TailObj struct {
	tail     *tail.Tail
	fileName string
	topic    string
}
type TextMsg struct {
	Msg string
}
type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
}

var (
	tailObjMgr *TailObjMgr
)

func InitTail() {
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, 10),
	}
}
func CreateTailFInstance(fileName []string) error {
	for _, fileItem := range fileName {
		//判断文件是否存在
		if _, err := os.Stat(fileItem); err != nil {
			return err
		}
		// 创建监控实例
		t, err := tail.TailFile(fileItem, tail.Config{
			ReOpen: true,
			Follow: true,
			//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll:      true,
		})
		if err != nil {
			return err
		}
		obj := &TailObj{
			tail:     t,
			fileName: fileItem}
		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)
		go ReadFormTailFInstance(obj)

	}
	return nil
}
func ReadFormTailFInstance(obj *TailObj) {
	for true {
		select {
		case line, ok := <-obj.tail.Lines:
			if !ok {
				time.Sleep(100 * time.Microsecond)
			}
			msg := &TextMsg{
				Msg: line.Text,
			}
			tailObjMgr.msgChan <- msg
		}
	}
}
func GetOneLine() *TextMsg {
	msg := <-tailObjMgr.msgChan
	return msg
}
