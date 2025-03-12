/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package layer

type Config struct {
	Version string `yaml:"version"`
	Package struct {
		ID          string `yaml:"id"`
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Kind        string `yaml:"kind"`
		Description string `yaml:"description"`
	} `yaml:"package"`
	Command []string `yaml:"command"`
	Base    string   `yaml:"base"`
	Runtime string   `yaml:"runtime,omitempty"`
	Build   string   `yaml:"build"`
}

const LayerMagic = "<<< deepin linglong layer archive >>>\x00\x00\x00"
const MkfsErofs = "mkfs.erofs"
const DumpErofs = "dump.erofs"
const ErofsFuse = "erofsfuse"
const FuserMount = "fusermount"

type LayerInfo struct {
	Arch          []string `json:"arch"`
	Base          string   `json:"base"`
	Runtime       string   `json:"runtime,omitempty"`
	Channel       string   `json:"channel"`
	Command       []string `json:"command"`
	Description   string   `json:"description"`
	ID            string   `json:"id"`
	Kind          string   `json:"kind"`
	Module        string   `json:"module"`
	Name          string   `json:"name"`
	SchemaVersion string   `json:"schema_version"`
	Size          int64    `json:"size"`
	Version       string   `json:"version"`
}
type LayerInfoHeader struct {
	Info    LayerInfo `json:"info"`
	Version string    `json:"version"`
}
