package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pbnjay/memory"
)

type Unpack struct {
	Href     string `xml:"href,attr"`
	ToDir    string `xml:"toDir,attr"`
	UseModes bool   `xml:"useModes,attr"`
}
type File struct {
	Href   string `xml:"href,attr"`
	ToFile string `xml:"toFile,attr"`
	Mode   string `xml:"mode,attr"`
}
type Resource struct {
	Os        string   `xml:"os,attr"`
	Bits      int      `xml:"bits,attr"`
	MinMemory string   `xml:"minMemory,attr"`
	MaxMemory string   `xml:"maxMemory,attr"`
	Unpack    []Unpack `xml:"unpack"`
	File      []File   `xml:"file"`
}
type Startup struct {
	Os        string `xml:"os,attr"`
	Bits      int    `xml:"bits,attr"`
	MinMemory string `xml:"minMemory,attr"`
	MaxMemory string `xml:"maxMemory,attr"`
	File      string `xml:"file,attr"`
}
type MWSXML struct {
	Resources []Resource `xml:"resources"`
	Startup   []Startup  `xml:"startup"`
}

func parseXml(path string) MWSXML {
	xmlFile, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Error read XML': %v", err.Error()))
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(fmt.Sprintf("Error read XML: %v", err.Error()))
	}

	var desc MWSXML
	xml.Unmarshal(byteValue, &desc)

	return desc
}

func getResources(xml MWSXML) []Resource {
	var result []Resource
	for _, rs := range xml.Resources {
		if matchCurrentComputer(rs.Os, rs.Bits, rs.MinMemory, rs.MaxMemory) {
			result = append(result, rs)
		}
	}
	return result
}

/**
 * Get the first allowed startup section.
 */
func getStartupFile(xml MWSXML) string {
	for _, st := range xml.Startup {
		if matchCurrentComputer(st.Os, st.Bits, st.MinMemory, st.MaxMemory) {
			if st.File == "" {
				panic("XML description error: file is not defined in startup section")
			}
			return st.File
		}
	}
	panic("XML description error: startup section is not defined for this computer")
}

/**
 * Checks if section info corresponds with current computer.
 * Minimum memory - exclusive, maximum memory - inclusive.
 */
func matchCurrentComputer(os string, bits int, min string, max string) bool {
	if os != "" && os != os_name() {
		return false
	}
	if bits != 0 && bits != os_bits() {
		return false
	}
	if min != "" && parseMemorySize(min) >= memory.TotalMemory() {
		return false
	}
	if max != "" && parseMemorySize(max) < memory.TotalMemory() {
		return false
	}

	return true
}

/**
 * Parse memory size like '2g' or '1500m'
 */
func parseMemorySize(sz string) uint64 {
	var mul uint64
	var p string
	if sz[len(sz)-1] == 'm' {
		p = sz[0 : len(sz)-1]
		mul = 1024 * 1024
	} else if sz[len(sz)-1] == 'g' {
		p = sz[0 : len(sz)-1]
		mul = 1024 * 1024 * 1024
	} else {
		mul = 1
		p = sz
	}
	r, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error parse memory size '"+sz+"': %v", err.Error()))
	}
	return r * mul
}
