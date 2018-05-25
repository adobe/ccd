//
//  Copyright 2017 Adobe.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package main

import (
	"flag"
	"log"

	"io"
	"os"

	"github.com/adobe/ccd"
	"github.com/adobe/ccd/libs/config"
	"github.com/pkg/errors"
)

func main() {
	args := parseArgs()
	if err := run(args); err != nil {
		f := "ERR: %s"
		if args.Debug {
			f = "ERR: %+v"
		}
		log.Fatalf(f, err)
	}
}

func run(args *Args) error {
	l, err := setupLogging(args)
	if err != nil {
		return err
	}
	defer l.Close()

	cfg, err := config.Load(args.ConfigFile)
	if err != nil {
		return err
	}

	return ccd.Run(l.Logger, cfg)
}

type Args struct {
	Debug      bool
	ConfigFile string
	LogFile    string
}

func parseArgs() *Args {
	args := new(Args)

	flag.StringVar(&args.ConfigFile, "config", "", "path to the config file")
	flag.StringVar(&args.LogFile, "log", "", "path to the log file (defaults to stdout)")
	flag.BoolVar(&args.Debug, "debug", false, "enable debug information")
	flag.Parse()

	return args
}

type Logger struct {
	*log.Logger

	logFile io.WriteCloser
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func setupLogging(args *Args) (*Logger, error) {
	l := new(Logger)
	l.Logger = log.New(os.Stdout, "", log.LstdFlags)

	if args.LogFile != "" {
		fh, err := os.Open(args.LogFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open logfile")
		}
		l.logFile = fh
		l.Logger.SetOutput(fh)
	}

	return l, nil
}
