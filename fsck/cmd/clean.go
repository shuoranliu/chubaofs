package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/chubaofs/chubaofs/proto"
	"github.com/chubaofs/chubaofs/sdk/meta"
	"github.com/chubaofs/chubaofs/util/log"
	"github.com/chubaofs/chubaofs/util/ump"
)

var gMetaWrapper *meta.MetaWrapper

func Clean(opt string) error {
	defer log.LogFlush()

	if MasterAddr == "" || VolName == "" {
		flag.Usage()
		return fmt.Errorf("Lack of parameters: master(%v) vol(%v)", MasterAddr, VolName)
	}

	ump.InitUmp("fsck", "")

	_, err := log.InitLog("fscklog", "fsck", log.InfoLevel, nil)
	if err != nil {
		return fmt.Errorf("Init log failed: %v", err)
	}

	masters := strings.Split(MasterAddr, meta.HostsSeparator)
	var metaConfig = &meta.MetaConfig{
		Volume:  VolName,
		Masters: masters,
	}

	gMetaWrapper, err = meta.NewMetaWrapper(metaConfig)
	if err != nil {
		return fmt.Errorf("NewMetaWrapper failed: %v", err)
	}

	switch opt {
	case "inode":
		err = cleanInodes()
		if err != nil {
			return fmt.Errorf("Clean inodes failed: %v", err)
		}
	case "dentry":
		err = cleanDentries()
		if err != nil {
			return fmt.Errorf("Clean dentries failed: %v", err)
		}
	case "evict":
		err = evictInodes()
		if err != nil {
			return fmt.Errorf("Evict inodes failed: %v", err)
		}
	default:
	}

	return nil
}

func evictInodes() error {
	mps, err := getMetaPartitions(MasterAddr, VolName)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, mp := range mps {
		cmdline := fmt.Sprintf("http://%s:9092/getAllInodes?pid=%d", strings.Split(mp.LeaderAddr, ":")[0], mp.PartitionID)
		wg.Add(1)
		go evictOnTime(&wg, cmdline)
	}

	wg.Wait()
	return nil
}

func evictOnTime(wg *sync.WaitGroup, cmdline string) {
	defer wg.Done()

	client := &http.Client{Timeout: 0}
	resp, err := client.Get(cmdline)
	if err != nil {
		log.LogErrorf("Get request failed: %v %v", cmdline, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.LogErrorf("Invalid status code: %v", resp.StatusCode)
		return
	}

	log.LogWritef("Dealing with meta partition: %v", cmdline)

	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		inode := &Inode{}
		err = dec.Decode(inode)
		if err != nil {
			log.LogErrorf("Decode inode failed: %v", err)
			return
		}
		doEvictInode(inode)
	}
	log.LogWritef("Done! Dealing with meta partition: %v", cmdline)
}
func cleanInodes() error {
	filePath := fmt.Sprintf("_export_%s/%s", VolName, obsoleteInodeDumpFileName)

	fp, err := os.Open(filePath)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(fp)
	for dec.More() {
		inode := &Inode{}
		if err = dec.Decode(inode); err != nil {
			return err
		}
		doEvictInode(inode)
	}

	return nil
}

func doEvictInode(inode *Inode) error {
	if inode.NLink != 0 || time.Since(time.Unix(inode.ModifyTime, 0)) < 24*time.Hour || !proto.IsRegular(inode.Type) {
		return nil
	}
	err := gMetaWrapper.Evict(inode.Inode)
	if err != nil {
		if err != syscall.ENOENT {
			return err
		}
	}
	log.LogWritef("%v", inode)
	return nil
}

func doUnlinkInode(inode *Inode) error {
	/*
	 * Do clean inode with the following exceptions:
	 * 1. nlink == 0, might be a temorary inode
	 * 2. size == 0 && ctime is close to current time, might be in the process of file creation
	 */
	//	if inode.NLink == 0 ||
	//		(inode.Size == 0 &&
	//			time.Unix(inode.CreateTime, 0).Add(24*time.Hour).After(time.Now())) {
	//		return nil
	//	}
	//
	//	err := gMetaWrapper.Unlink_ll(inode.Inode)
	//	if err != nil {
	//		if err != syscall.ENOENT {
	//			return err
	//		}
	//		err = nil
	//	}
	//
	//	err = gMetaWrapper.Evict(inode.Inode)
	//	if err != nil {
	//		if err != syscall.ENOENT {
	//			return err
	//		}
	//	}

	return nil
}

func cleanDentries() error {
	//filePath := fmt.Sprintf("_export_%s/%s", VolName, obsoleteDentryDumpFileName)
	// TODO: send request to meta node directly with pino, name and ino.
	return nil
}
