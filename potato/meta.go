package potato

import (
	//"errors"
	"strconv"
	"strings"

	"github.com/couchbase/moss"
)

func MetaGet(prefix string, key string) ([]byte, error) {
	metakey := strings.Join([]string{prefix, key}, ":")
	Logger.Debug("Meta Get:", metakey)
	data, err := cm_kget(metakey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MetaSet(prefix string, key string, data []byte) error {
	metakey := strings.Join([]string{prefix, key}, ":")
	Logger.Debug("Meta Set:", metakey)
	if err := cm_set(metakey, data); err != nil {
		return err
	}
	return nil
}

func MetaDelete(prefix string, key string) error {
	metakey := strings.Join([]string{prefix, key}, ":")
	Logger.Debug("Meta Delete:", metakey)
	if err := cm_del(metakey); err != nil {
		return err
	}
	return nil
}

func MetaMultiDelete(keys []string) error {
	Logger.Debug("Meta Multi-Delete: ", len(keys))
	if err := cm_mdel(keys); err != nil {
		return err
	}
	return nil
}

func MetaSyncCount() (res int) {
	res = 0
	ssm, err := CMETA.Snapshot()
	if err != nil {
		res = -1
	}

	iter, err := ssm.StartIterator([]byte("sync:"), nil, moss.IteratorOptions{})
	if err != nil || iter == nil {
		return -1
	}

	err = iter.Next()
	if err != nil {
		res = -1
	} else {
		res++
	}

	ssm.Close()
	return res
}

func MetaSyncList() (listHtml string) {
	slaves := CFG.Replication.Slaves
	fileKeys := []string{}
	if len(slaves) > 0 {

		ssm, err := CMETA.Snapshot()
		if err != nil {
			Logger.Error("expected ssm ok")
		}

		iter, err := ssm.StartIterator([]byte("sync:"), nil, moss.IteratorOptions{})
		if err != nil || iter == nil {
			Logger.Error("expected iter")
		}

		for i := 0; i < 1000; i++ {
			k, _, err := iter.Current()
			if err != nil {
				continue
			}
			if k != nil {
				k_raw := string(k)
				fileKeys = append(fileKeys, k_raw)
			}

			err = iter.Next()
			if err != nil {
				break
			}
		}

		ssm.Close()

	}

	fileKeys_len := len(fileKeys)
	if fileKeys_len == 0 {
		Logger.Debug("No Entities Replication Needed.")
		return ""
	}
	listHtml = ""
	for k, v := range fileKeys {
		listHtml = strings.Join([]string{v, " : ", strconv.Itoa(k), "<br/>", listHtml}, "")
	}

	return listHtml
}

func MetaList() (listHtml string) {
	fileKeys := []string{}

	ssm, err := CMETA.Snapshot()
	if err != nil {
		Logger.Error("expected ssm ok")
	}

	iter, err := ssm.StartIterator(nil, nil, moss.IteratorOptions{})
	if err != nil || iter == nil {
		Logger.Error("expected iter")
	}

	for i := 0; i < 500; i++ {
		k, _, err := iter.Current()
		if err != nil {
			continue
		}
		if k != nil {
			fileKeys = append(fileKeys, string(k))
		}

		err = iter.Next()
		if err != nil {
			break
		}
	}

	ssm.Close()

	fileKeys_len := len(fileKeys)
	if fileKeys_len == 0 {
		Logger.Debug("No Entities Replication Needed.")
		return ""
	}
	listHtml = ""
	for k, v := range fileKeys {
		listHtml = strings.Join([]string{v, " : ", strconv.Itoa(k), "<br/>", listHtml}, "")
	}

	return listHtml
}
