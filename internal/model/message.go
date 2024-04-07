// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package model

import (
	"encoding/xml"
	"time"
)

type Image struct {
	MediaId string `xml:"MediaId"`
}

type Voice struct {
	MediaId string `xml:"MediaId"`
}

type Video struct {
	MediaId string `xml:"MediaId"`
}

type Music struct {
	Title        string `xml:"Title"`
	Description  string `xml:"Description"`
	MusicURL     string `xml:"MusicURL"`
	HQMusicUrl   string `xml:"HQMusicUrl"`
	ThumbMediaId string `xml:"ThumbMediaId"`
}

type Article struct {
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	PicUrl      string `xml:"PicUrl"`
	Url         string `xml:"Url"`
}

type Articles struct {
	Item []Article `xml:"item"`
}

// This `TextMessage` is the struct which passively replies to the user.
type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	// Event        string   `xml:"Event,attr,omitempty"`
	Content string `xml:"Content"`
	// Image        Image    `xml:"Image,omitempty"`
	// Voice        Voice    `xml:"Voice,omitempty"`
	// Video        Video    `xml:"Video,omitempty"`
	// Music        Music    `xml:"Music,omitempty"`
	// Articles     Articles `xml:"Articles,omitempty"`
	// Recognition  string   `xml:"Recognition,attr,omitempty"`
	MsgId     int64  `xml:"MsgId"`
	MsgDataId string `xml:"MsgDataId"`
	Idx       string `xml:"Idx"`
}

func NewMsg(data []byte) *any {
	var msg any
	if err := xml.Unmarshal(data, &msg); err != nil {
		return nil
	}
	return &msg
}

type MessageType int

const (
	TextType MessageType = iota + 1
	ImageType
	VoiceType
	MusicType
	VideoType
	ArticlesType
)

func (msg *Message) GenTextData(s string) []byte {
	data := Message{
		ToUserName:   msg.ToUserName,
		FromUserName: msg.FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}

	bs, err := xml.Marshal(&data)
	if err != nil {
		return bs
	}

	return nil
}

type SpecifiedMessage struct {
	ID      uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Key     string `gorm:"column:key"`
	Message string `gorm:"column:message"`
}
