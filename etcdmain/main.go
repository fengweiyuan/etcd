// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdmain

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-systemd/daemon"
	systemdutil "github.com/coreos/go-systemd/util"
	"go.uber.org/zap"
)

/**
 * 启动一个etcd实例
 */
func Main() {
	// 检查系统的架构是否可以满足要求
	checkSupportArch()

	// 如果启动时在命令后有跟随参数，那么...
	if len(os.Args) > 1 {
		cmd := os.Args[1]                                                               // 第一个参数是一个命令
		if covArgs := os.Getenv("ETCDCOV_ARGS"); len(covArgs) > 0 {                // 如果有定义ETCDCOV_ARGS环境变量，那么
			args := strings.Split(os.Getenv("ETCDCOV_ARGS"), "\xe7\xcd")[1:]  // 不懂\xe7这些是什么
			rootCmd.SetArgs(args)
			cmd = "grpc-proxy"
		}
		switch cmd {
		case "gateway", "grpc-proxy":
			if err := rootCmd.Execute(); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}

	// 如果没有跟随参数，那么直接启动。
	startEtcdOrProxyV2()
}

func notifySystemd(lg *zap.Logger) {
	if !systemdutil.IsRunningSystemd() {
		return
	}

	if lg != nil {
		lg.Info("host was booted with systemd, sends READY=1 message to init daemon")
	}
	sent, err := daemon.SdNotify(false, "READY=1")
	if err != nil {
		if lg != nil {
			lg.Error("failed to notify systemd for readiness", zap.Error(err))
		} else {
			plog.Errorf("failed to notify systemd for readiness: %v", err)
		}
	}

	if !sent {
		if lg != nil {
			lg.Warn("forgot to set Type=notify in systemd service file?")
		} else {
			plog.Errorf("forgot to set Type=notify in systemd service file?")
		}
	}
}
