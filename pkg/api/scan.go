package api

import (
	"bytes"
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (h *LeafHandler) HandleScanMsg(m *nats.Msg) {
	startKey := m.Header.Get("Start-Key")
	endKey := m.Header.Get("End-Key")
	var res Response

	items := []KV{}

	err := h.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if (startKey == "" || bytes.Compare(k, []byte(startKey)) >= 0) && (endKey == "" || bytes.Compare(k, []byte(endKey)) <= 0) {
				err := item.Value(func(v []byte) error {
					items = append(items, KV{Key: string(k), Value: string(v)})
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		res.Error = err.Error()
		h.Logger.Error("Error scanning values", zap.Error(err))
	} else {
		res.Data = items
		h.Logger.Info("Successfully scanned values", zap.Any("items", items))
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		h.Logger.Error("Error encoding response to JSON", zap.Error(err))
		m.RespondMsg(&nats.Msg{Subject: m.Reply, Data: []byte(err.Error())})
		return
	}

	m.RespondMsg(&nats.Msg{Subject: m.Reply, Data: jsonRes})
}
