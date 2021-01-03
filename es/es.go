package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/feitianlove/golib/common/logger"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

type ConfigES struct {
	Address string
	Port    string
}
type ClientES struct {
	Client *elastic.Client
}
type DataStructES struct {
	Index   string
	Type    string
	Id      string
	Data    interface{}
	Context context.Context
}

func NewEsClient(conf *ConfigES) (*ClientES, error) {
	// 这里不用判断空，
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("%s:%s", conf.Address, conf.Port)))
	if err != nil {
		logger.Ctrl.WithFields(logrus.Fields{
			"error": err,
		}).Error("NewEsClient err")
		return nil, err
	}
	return &ClientES{Client: client}, nil
}
func (es *ClientES) CreateRecord(data DataStructES) (*elastic.IndexResponse, error) {
	result, err := es.Client.Index().Index(data.Index).
		Type(data.Type).Id(data.Id).BodyJson(data.Data).Do(data.Context)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"error": err,
			"data":  data,
		}).Error("ES CreateRecord  err")
		return nil, err
	}
	logger.Console.WithFields(logrus.Fields{
		"ID":   result.Id,
		"data": data,
	}).Info("ES CreateRecord success")
	return result, nil
}

func (es *ClientES) DeleteRecord(user, data DataStructES) (*elastic.DeleteResponse, error) {

	res, err := es.Client.Delete().Index("megacorp").
		Type(data.Type).
		Id(data.Id).
		Do(data.Context)
	if err != nil {
		println(err.Error())
		logger.Console.WithFields(logrus.Fields{
			"user":  user,
			"data":  data,
			"error": err,
		}).Error("ES deleteRecord err")
		return nil, err
	}
	logger.Console.WithFields(logrus.Fields{
		"user":  user,
		"error": err,
		"data":  data,
	}).Info("ES deleteRecord success")
	return res, nil
}

//修改
func (es *ClientES) UpdateRecord(user, data DataStructES) (*elastic.UpdateResponse, error) {
	res, err := es.Client.Update().
		Index(data.Index).
		Type(data.Type).
		Id(data.Id).
		Doc(data.Data).
		Do(data.Context)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"user":  user,
			"data":  data,
			"error": err,
		}).Error("ES UpdateRecord err")
	}
	logger.Console.WithFields(logrus.Fields{
		"user":  user,
		"error": err,
		"data":  data,
	}).Info("ES UpdateRecord success")
	return res, nil
}

func (es *ClientES) SearchById(data DataStructES) (*elastic.GetResult, error) {
	result, err := es.Client.Get().
		Index(data.Index).
		Type(data.Type).
		Id(data.Id).
		Do(context.Background())
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"error": err,
			"data":  data,
		}).Error("ES SearchById err")
	}
	logger.Console.WithFields(logrus.Fields{
		"result": result,
	}).Error("ES SearchById success")
	if result != nil && result.Found {
		return result, nil
	}
	return nil, errors.New("search result is nil")
}
