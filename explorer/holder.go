// Copyright (c) 2023-2024 The UXUY Developer Team
// License:
// MIT License

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE

package explorer

import (
	"github.com/uxuycom/indexer/dcache"
	"github.com/uxuycom/indexer/xylog"
	"time"
)

func (e *Explorer) CountHolder() {
	defer func() {
		e.cancel()
	}()

	_ = e.updateInsHolder()

	t := time.NewTicker(30 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := e.updateInsHolder(); err != nil {
				xylog.Logger.Errorf("failed to obtain the current block height. chain:%s err=%s", e.config.Chain.ChainName, err)
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *Explorer) updateInsHolder() error {
	idx := 0
	start := uint32(0)
	limit := 10000
	xylog.Logger.Infof("updateInsHolder load inscription-stats data start...")
	for {
		items, err := e.db.GetInscriptionStatsByIdLimit(e.config.Chain.ChainName, uint64(start), limit)
		if err != nil {
			xylog.Logger.Fatalf("failed inscription-stats  data. err:%v", err)
		}
		idx++
		xylog.Logger.Infof("load inscription-stats ret, items[%d], idx:%d", len(items), idx)

		if len(items) <= 0 {
			break
		}

		for _, v := range items {
			holder, err := e.updateHolders(v.Protocol, v.Tick)
			if err == nil {
				v.Holders = uint64(holder)
				params := make(map[string]interface{})
				params["holders"] = holder
				xylog.Logger.Infof("update holder tick[%s], idx:%d, holder[%d]", v.Tick, idx, holder)
				err := e.db.UpdateInscriptionsStatsBySID(e.db.SqlDB, e.config.Chain.ChainName, v.SID, params)
				if err != nil {
					xylog.Logger.Fatalf("failed updateHolders. err:%v", err)
				}
			} else {
				xylog.Logger.Errorf("update holder error. err[%s]", err)
			}
		}

		//update id index
		start = items[len(items)-1].ID
	}
	return nil
}

func (e *Explorer) updateHolders(protocol, tick string) (int64, error) {
	count, err := e.db.GetHolderNumberByTick(e.config.Chain.ChainName, protocol, tick)
	if err != nil {
		xylog.Logger.Infof("error getting holder. err[%s]", err)
		return 0, err
	}
	xylog.Logger.Infof("GetHolderNumberByTick count[%d]", count)

	ok, _ := e.dCache.InscriptionStats.Get(protocol, tick)
	if ok {
		insStats := &dcache.InsStats{
			Holders: count,
		}
		e.dCache.InscriptionStats.Update(protocol, tick, insStats)
	}

	return count, nil
}
