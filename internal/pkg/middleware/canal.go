package middleware

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/withlin/canal-go/client"
	prot "github.com/withlin/canal-go/protocol/entry"
	"log"
	"os"
	"time"
)

type Canal struct {
	cli         *client.SimpleCanalConnector
	address     string
	port        int
	username    string
	password    string
	soTime      int32
	idleTimeout int32
	batchSize   int32
}

func NewCanalClient(conf *viper.Viper) (*client.SimpleCanalConnector, error) {

	connector := &client.SimpleCanalConnector{
		Address:     conf.GetString("data.canal.addr"),
		Port:        conf.GetInt("data.canal.port"),
		UserName:    conf.GetString("data.canal.username"),
		PassWord:    conf.GetString("data.canal.password"),
		SoTime:      conf.GetInt32("data.canal.sotime"),
		IdleTimeOut: conf.GetInt32("data.canal.idletimeout"),
	}
	err := connector.Connect()
	return connector, err
}

func (c *Canal) SyncData(conf *viper.Viper) {
	cli := c.cli
	// subscribe all tables from database.
	err := cli.Subscribe(".*\\..*")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	batchSize := conf.GetInt32("data.canal.ack.batchSize")

	for {

		msg, err := cli.Get(batchSize, nil, nil)
		if err != nil {
			return
		}
		id := msg.Id
		if id == -1 || len(msg.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		optEntry(msg.Entries)
	}
}

func optEntry(entrys []prot.Entry) {
	for _, entry := range entrys {
		if entry.GetEntryType() == prot.EntryType_TRANSACTIONBEGIN ||
			entry.GetEntryType() == prot.EntryType_TRANSACTIONEND {
			continue
		}
		row := new(prot.RowChange)
		err := proto.Unmarshal(entry.GetStoreValue(), row)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Errors:%s", err.Error())
			os.Exit(1)
		}
		if row != nil {
			eventType := row.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
			for _, data := range row.GetRowDatas() {
				if eventType == prot.EventType_DELETE {
					// delete the data
					data.GetBeforeColumns()
				} else if eventType == prot.EventType_INSERT {
					// insert the data
				} else if eventType == prot.EventType_UPDATE {
					// update the data
				} else {
					// do nothing
				}
			}
		}
	}
}
