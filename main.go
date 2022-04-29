// Copyright 2022 Billy Lynch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

type modConfig struct {
	Modules []*Module `json:"modules"`
}

type Module struct {
	URL      string `json:"url,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Path     string `json:"path,omitempty"`
}

func main() {
	configPath := "getsum.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	b, _ := io.ReadAll(f)
	cfg := new(modConfig)
	if err := yaml.Unmarshal(b, cfg); err != nil {
		panic(err)
	}
	f.Close()

	for _, m := range cfg.Modules {
		resp, _ := http.Get(m.URL)

		u, _ := url.Parse(m.URL)
		fmt.Printf("fetched %q\n", m.URL)

		path := m.Path
		if path == "" {
			path = filepath.Base(u.Path)
		}
		fmt.Println("\toutput:", path)

		f, _ := os.Create(path)
		defer resp.Body.Close()
		defer f.Close()
		h := sha256.New()

		w := io.MultiWriter(f, h)
		io.Copy(w, resp.Body)

		checksum := hex.EncodeToString(h.Sum(nil))
		fmt.Println("\tchecksum:", checksum)

		if m.Checksum == "" {
			m.Checksum = checksum
		} else {
			if checksum != m.Checksum {
				fmt.Printf("\tchecksums do not match! want %s, got %s\n", m.Checksum, checksum)
				os.Exit(1)
			}
		}
	}

	f, _ = os.Create(f.Name())
	b, _ = yaml.Marshal(cfg)
	f.Write(b)
}
