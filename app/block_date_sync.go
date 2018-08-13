package app

import (
	"log"
	"time"
)

func StartSyncingBlockDate(context *Context) error {
	max, err := context.TableBlockDate.Max(context.DB)

	if err != nil {
		return err
	}

	log.Println("Start syncing from: ", max)
	go syncBlockDate(context, max)

	return nil
}

func syncBlockDate(context *Context, currentBlock int64) {

	for {
		result, err := context.Client.Eth.GetBlockByNumber(currentBlock, false)

		if err != nil {
			log.Println(err.Error())
			continue
		}

		bd := &BlockDate{
			Block: currentBlock,
			Date:  result.Timestamp,
		}

		err = context.TableBlockDate.Put(bd, context.DB)

		if err != nil {
			log.Println(err.Error())
			continue
		}

		time.Sleep(100 * time.Millisecond)

		log.Println("Block: ", currentBlock)
		currentBlock++
	}
}
