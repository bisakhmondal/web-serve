// Package core
/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* 	http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func DurationCast(t int, d time.Duration) time.Duration {
	return time.Duration(t) * d
}

func CombinedWriter() (io.Writer, func() error) {
	f, err := os.OpenFile(filepath.Join(conf.Configuration.Logging.Logdir,
		conf.Configuration.Logging.Logfile),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating log file: %s\n", err)
		os.Exit(1)
	}

	writer := io.Writer(f)
	if conf.Configuration.Logging.Stdout {
		writer = io.MultiWriter(f, os.Stdout)
	}
	return writer, f.Close
}
