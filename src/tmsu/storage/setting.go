/*
Copyright 2011-2015 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package storage

import (
	"tmsu/entities"
)

var defaultSettings = map[string]string{
	"autoCreateTags":                "yes",
	"autoCreateValues":              "yes",
	"fileFingerprintAlgorithm":      "dynamic:SHA256",
	"directoryFingerprintAlgorithm": "sumSizes",
}

// The complete set of settings.
func (storage *Storage) Settings() (entities.Settings, error) {
	settings, err := storage.Db.Settings()
	if err != nil {
		return nil, err
	}

	// enrich with defaults
	for name, value := range defaultSettings {
		if !settings.ContainsName(name) {
			settings = append(settings, &entities.Setting{name, value})
		}
	}

	return settings, nil
}

func (storage *Storage) Setting(name string) (*entities.Setting, error) {
	setting, err := storage.Db.Setting(name)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		value, ok := defaultSettings[name]
		if !ok {
			return nil, nil
		}

		setting = &entities.Setting{name, value}
	}

	return setting, nil
}

func (storage *Storage) UpdateSetting(name, value string) (*entities.Setting, error) {
	return storage.Db.UpdateSetting(name, value)
}
