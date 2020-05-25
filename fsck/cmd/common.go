package cmd

import (
	"encoding/json"
	"flag"
)

var (
	MasterAddr string
	VolName    string
	InodesFile string
	DensFile   string
)

func init() {
	flag.StringVar(&MasterAddr, "master", "", "master address")
	flag.StringVar(&VolName, "vol", "", "volume name")
	flag.StringVar(&InodesFile, "inodes", "", "inode list file")
	flag.StringVar(&DensFile, "dens", "", "dentry list file")
}

var (
	inodeDumpFileName          string = "inode.dump"
	dentryDumpFileName         string = "dentry.dump"
	inodeUpdateDumpFileName    string = "inode.dump.update"
	obsoleteInodeDumpFileName  string = "inode.dump.obsolete"
	obsoleteDentryDumpFileName string = "dentry.dump.obsolete"
)

type Inode struct {
	Inode      uint64
	Type       uint32
	Size       uint64
	CreateTime int64
	AccessTime int64
	ModifyTime int64
	NLink      uint32

	Dens  []*Dentry
	Valid bool
}

func (i *Inode) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

type Dentry struct {
	ParentId uint64
	Name     string
	Inode    uint64
	Type     uint32

	Valid bool
}

func (d *Dentry) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(data)
}
