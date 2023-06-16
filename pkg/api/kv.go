package api

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (h *LeafHandler) HandleKVMsg(m *nats.Msg) {
	key := m.Subject[9:]
	var res Response

	switch m.Header.Get("Ops") {
	case "get":
		var item KV

		err := h.DB.View(func(txn *badger.Txn) error {
			val, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			value, err := val.ValueCopy(nil)
			if err != nil {
				return err
			}

			item = KV{Key: key, Value: string(value)}
			return nil
		})

		if err != nil {
			res.Error = err.Error()
			h.Logger.Error("Error getting value", zap.Error(err))
		} else {
			res.Data = item
			h.Logger.Info("Successfully got value", zap.Any("item", item))
		}

	case "set":
		err := h.DB.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(key), m.Data)
			return err
		})

		if err != nil {
			res.Error = err.Error()
			h.Logger.Error("Error setting value", zap.Error(err))
		} else {
			res.Data = "OK"
			h.Logger.Info("Successfully set value", zap.String("key", key))
		}

	case "delete":
		err := h.DB.Update(func(txn *badger.Txn) error {
			err := txn.Delete([]byte(key))
			return err
		})

		if err != nil {
			res.Error = err.Error()
			h.Logger.Error("Error deleting value", zap.Error(err))
		} else {
			res.Data = "OK"
			h.Logger.Info("Successfully deleted value", zap.String("key", key))
		}

	default:
		res.Error = "Unknown operation: " + m.Header.Get("Ops")
		h.Logger.Error("Unknown operation", zap.String("Ops", m.Header.Get("Ops")))
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		h.Logger.Error("Error encoding response to JSON", zap.Error(err))
		m.RespondMsg(&nats.Msg{Subject: m.Reply, Data: []byte("Error encoding response to JSON")})
		return
	}

	m.RespondMsg(&nats.Msg{Subject: m.Reply, Data: jsonRes})
}
