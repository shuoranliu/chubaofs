// Copyright 2020 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/chubaofs/chubaofs/fsck/cmd"
)

//var (
//	checkOpt string
//	cleanOpt string
//	evictOpt bool
//)

//func init() {
//	flag.StringVar(&checkOpt, "check", "", "check and export obsolete inodes and dentries")
//	flag.StringVar(&cleanOpt, "clean", "", "clean inodes or dentries")
//	flag.BoolVar(&evictOpt, "evict", false, "evict inodes whose nlink is zero")
//
//	flag.Parse()
//}

func main() {
	c := cmd.NewRootCmd()
	if err := c.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %v\n", err)
		os.Exit(1)
	}
	//	if checkOpt != "" {
	//		err := cmd.Check(checkOpt)
	//		if err != nil {
	//			fmt.Println(err)
	//			os.Exit(1)
	//		}
	//	}
	//
	//	if cleanOpt != "" && (cleanOpt == "inode" || cleanOpt == "dentry" || cleanOpt == "evict") {
	//		err := cmd.Clean(cleanOpt)
	//		if err != nil {
	//			fmt.Println(err)
	//			os.Exit(1)
	//		}
	//	}
}
