// Copyright (C) 2019-2022 Chrystian Huot <chrystian.huot@saubeo.solutions>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type Options struct {
	AfsSystems                  string `json:"afsSystems"`
	AutoPopulate                bool   `json:"autoPopulate"`
	DimmerDelay                 uint   `json:"dimmerDelay"`
	DisableAudioConversion      bool   `json:"disableAudioConversion"`
	DisableDuplicateDetection   bool   `json:"disableDuplicateDetection"`
	DuplicateDetectionTimeFrame uint   `json:"duplicateDetectionTimeFrame"`
	KeypadBeeps                 string `json:"keypadBeeps"`
	MaxClients                  uint   `json:"maxClients"`
	PruneDays                   uint   `json:"pruneDays"`
	SearchPatchedTalkgroups     bool   `json:"searchPatchedTalkgroups"`
	ShowListenersCount          bool   `json:"showListenersCount"`
	SortTalkgroups              bool   `json:"sortTalkgroups"`
	TagsToggle                  bool   `json:"tagsToggle"`
	adminPassword               string
	adminPasswordNeedChange     bool
	mutex                       sync.Mutex
	secret                      string
}

func NewOptions() *Options {
	return &Options{
		mutex: sync.Mutex{},
	}
}

func (options *Options) FromMap(m map[string]interface{}) {
	options.mutex.Lock()
	defer options.mutex.Unlock()

	switch v := m["autoPopulate"].(type) {
	case bool:
		options.AutoPopulate = v
	default:
		options.AutoPopulate = defaults.options.autoPopulate
	}

	switch v := m["dimmerDelay"].(type) {
	case float64:
		options.DimmerDelay = uint(v)
	default:
		options.DimmerDelay = defaults.options.dimmerDelay
	}

	switch v := m["disableAudioConversion"].(type) {
	case bool:
		options.DisableAudioConversion = v
	default:
		options.DisableAudioConversion = defaults.options.disableAudioConversion
	}

	switch v := m["disableDuplicateDetection"].(type) {
	case bool:
		options.DisableDuplicateDetection = v
	default:
		options.DisableDuplicateDetection = defaults.options.disableDuplicateDetection
	}

	switch v := m["duplicateDetectionTimeFrame"].(type) {
	case float64:
		options.DuplicateDetectionTimeFrame = uint(v)
	default:
		options.DuplicateDetectionTimeFrame = defaults.options.duplicateDetectionTimeFrame
	}

	switch v := m["afsSystems"].(type) {
	case string:
		options.AfsSystems = v
	}

	switch v := m["keypadBeeps"].(type) {
	case string:
		options.KeypadBeeps = v
	default:
		options.KeypadBeeps = defaults.options.keypadBeeps
	}

	switch v := m["maxClients"].(type) {
	case float64:
		options.MaxClients = uint(v)
	default:
		options.MaxClients = defaults.options.maxClients
	}

	switch v := m["pruneDays"].(type) {
	case float64:
		options.PruneDays = uint(v)
	default:
		options.PruneDays = defaults.options.pruneDays
	}

	switch v := m["searchPatchedTalkgroups"].(type) {
	case bool:
		options.SearchPatchedTalkgroups = v
	default:
		options.SearchPatchedTalkgroups = defaults.options.searchPatchedTalkgroups
	}

	switch v := m["showListenersCount"].(type) {
	case bool:
		options.ShowListenersCount = v
	default:
		options.ShowListenersCount = defaults.options.showListenersCount
	}

	switch v := m["sortTalkgroups"].(type) {
	case bool:
		options.SortTalkgroups = v
	default:
		options.SortTalkgroups = defaults.options.sortTalkgroups
	}

	switch v := m["tagsToggle"].(type) {
	case bool:
		options.TagsToggle = v
	default:
		options.TagsToggle = defaults.options.tagsToggle
	}
}

func (options *Options) Read(db *Database) error {
	var (
		defaultPassword []byte
		err             error
		f               interface{}
	)

	options.mutex.Lock()
	defer options.mutex.Unlock()

	defaultPassword, _ = bcrypt.GenerateFromPassword([]byte(defaults.adminPassword), bcrypt.DefaultCost)

	options.adminPassword = string(defaultPassword)
	options.adminPasswordNeedChange = defaults.adminPasswordNeedChange
	options.AutoPopulate = defaults.options.autoPopulate
	options.DimmerDelay = defaults.options.dimmerDelay
	options.DisableAudioConversion = defaults.options.disableAudioConversion
	options.DisableDuplicateDetection = defaults.options.disableDuplicateDetection
	options.DuplicateDetectionTimeFrame = defaults.options.duplicateDetectionTimeFrame
	options.KeypadBeeps = defaults.options.keypadBeeps
	options.MaxClients = defaults.options.maxClients
	options.PruneDays = defaults.options.pruneDays
	options.SearchPatchedTalkgroups = defaults.options.searchPatchedTalkgroups
	options.ShowListenersCount = defaults.options.showListenersCount
	options.SortTalkgroups = defaults.options.sortTalkgroups
	options.TagsToggle = defaults.options.tagsToggle

	err = db.Sql.QueryRow("select `val` from `rdioScannerConfigs` where `key` = 'adminPassword'").Scan(&f)
	if err == nil {
		switch v := f.(type) {
		case []uint8:
			var f string
			if err = json.Unmarshal(v, &f); err == nil {
				options.adminPassword = f
			}
		case string:
			var f string
			if err = json.Unmarshal([]byte(v), &f); err == nil {
				options.adminPassword = f
			}
		}
	}

	err = db.Sql.QueryRow("select `val` from `rdioScannerConfigs` where `key` = 'adminPasswordNeedChange'").Scan(&f)
	if err == nil {
		switch v := f.(type) {
		case []uint8:
			var f bool
			if err = json.Unmarshal(v, &f); err == nil {
				options.adminPasswordNeedChange = f
			}
		case string:
			var f bool
			if err = json.Unmarshal([]byte(v), &f); err == nil {
				options.adminPasswordNeedChange = f
			}
		}
	}

	err = db.Sql.QueryRow("select `val` from `rdioScannerConfigs` where `key` = 'options'").Scan(&f)
	if err == nil {
		switch v := f.(type) {
		case string:
			if err = json.Unmarshal([]byte(v), &f); err == nil {
				switch v := f.(type) {
				case map[string]interface{}:
					switch v := v["afsSystems"].(type) {
					case string:
						options.AfsSystems = v
					}

					switch v := v["autoPopulate"].(type) {
					case bool:
						options.AutoPopulate = v
					}

					switch v := v["dimmerDelay"].(type) {
					case float64:
						options.DimmerDelay = uint(v)
					}

					switch v := v["disableAudioConversion"].(type) {
					case bool:
						options.DisableAudioConversion = v
					}

					switch v := v["disableDuplicateDetection"].(type) {
					case bool:
						options.DisableDuplicateDetection = v
					}

					switch v := v["duplicateDetectionTimeFrame"].(type) {
					case float64:
						options.DuplicateDetectionTimeFrame = uint(v)
					}

					switch v := v["keypadBeeps"].(type) {
					case string:
						options.KeypadBeeps = v
					}

					switch v := v["maxClients"].(type) {
					case float64:
						options.MaxClients = uint(v)
					}

					switch v := v["pruneDays"].(type) {
					case float64:
						options.PruneDays = uint(v)
					}

					switch v := v["searchPatchedTalkgroups"].(type) {
					case bool:
						options.SearchPatchedTalkgroups = v
					}

					switch v := v["showListenersCount"].(type) {
					case bool:
						options.ShowListenersCount = v
					}

					switch v := v["sortTalkgroups"].(type) {
					case bool:
						options.SortTalkgroups = v
					}

					switch v := v["tagsToggle"].(type) {
					case bool:
						options.TagsToggle = v
					}
				}
			}
		}
	}

	err = db.Sql.QueryRow("select `val` from `rdioScannerConfigs` where `key` = 'secret'").Scan(&f)
	if err == nil {
		switch v := f.(type) {
		case string:
			var f string
			if err = json.Unmarshal([]byte(v), &f); err == nil {
				options.secret = f
			}
		}
	}

	return nil
}

func (options *Options) Write(db *Database) error {
	var (
		b   []byte
		err error
		i   int64
		res sql.Result
	)

	options.mutex.Lock()
	defer options.mutex.Unlock()

	formatError := func(err error) error {
		return fmt.Errorf("options.write: %v", err)
	}

	if b, err = json.Marshal(options.adminPassword); err != nil {
		return formatError(err)
	}

	if res, err = db.Sql.Exec("update `rdioScannerConfigs` set `val` = ? where `key` = 'adminPassword'", string(b)); err != nil {
		return formatError(err)
	}

	if i, err = res.RowsAffected(); err == nil && i == 0 {
		db.Sql.Exec("insert into `rdioScannerConfigs` (`key`, `val`) values (?, ?)", "adminPassword", string(b))
	}

	if b, err = json.Marshal(options.adminPasswordNeedChange); err != nil {
		return formatError(err)
	}

	if res, err = db.Sql.Exec("update `rdioScannerConfigs` set `val` = ? where `key` = 'adminPasswordNeedChange'", string(b)); err != nil {
		return formatError(err)
	}

	if i, err = res.RowsAffected(); err == nil && i == 0 {
		db.Sql.Exec("insert into `rdioScannerConfigs` (`key`, `val`) values (?, ?)", "adminPasswordNeedChange", string(b))
	}

	if b, err = json.Marshal(map[string]interface{}{
		"afsSystems":                  options.AfsSystems,
		"autoPopulate":                options.AutoPopulate,
		"dimmerDelay":                 options.DimmerDelay,
		"disableAudioConversion":      options.DisableAudioConversion,
		"disableDuplicateDetection":   options.DisableDuplicateDetection,
		"duplicateDetectionTimeFrame": options.DuplicateDetectionTimeFrame,
		"keypadBeeps":                 options.KeypadBeeps,
		"maxClients":                  options.MaxClients,
		"pruneDays":                   options.PruneDays,
		"searchPatchedTalkgroups":     options.SearchPatchedTalkgroups,
		"showListenersCount":          options.ShowListenersCount,
		"sortTalkgroups":              options.SortTalkgroups,
		"tagsToggle":                  options.TagsToggle,
	}); err != nil {
		return formatError(err)
	}

	if res, err = db.Sql.Exec("update `rdioScannerConfigs` set `val` = ? where `key` = 'options'", string(b)); err != nil {
		return formatError(err)
	}

	if i, err = res.RowsAffected(); err == nil && i == 0 {
		db.Sql.Exec("insert into `rdioScannerConfigs` (`key`, `val`) values (?, ?)", "options", string(b))
	}

	return nil
}
