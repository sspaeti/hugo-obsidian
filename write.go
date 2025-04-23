package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func write(links []Link, contentIndex ContentIndex, toIndex bool, out string, root string) error {
	index := index(links)
	resStruct := struct {
		Index Index  `json:"index"`
		Links []Link `json:"links"`
	}{
		Index: index,
		Links: links,
	}
	marshalledIndex, mErr := json.MarshalIndent(&resStruct, "", "  ")
	if mErr != nil {
		return mErr
	}

	writeErr := ioutil.WriteFile(path.Join(out, "linkIndex.json"), marshalledIndex, 0644)
	if writeErr != nil {
		return writeErr
	}

	// check whether to index content
	if toIndex {
		marshalledContentIndex, mcErr := json.MarshalIndent(&contentIndex, "", "  ")
		if mcErr != nil {
			return mcErr
		}

		writeErr = ioutil.WriteFile(path.Join(out, "contentIndex.json"), marshalledContentIndex, 0644)
		if writeErr != nil {
			return writeErr
		}

		// write linkmap
		writeErr = writeLinkMap(&contentIndex, root)
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}

func writeLinkMap(contentIndex *ContentIndex, root string) error {
	fp := path.Join(root, "static", "linkmap")
	file, err := os.OpenFile(fp, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)
	for path := range *contentIndex {
		if path == "/" {
			_, _ = datawriter.WriteString("/index.html /\n")
		} else {
			_, _ = datawriter.WriteString(path + "/index.{html} " + path + "/\n")
		}
	}
	datawriter.Flush()
	file.Close()

	return nil
}

// constructs index from links
func index(links []Link) (index Index) {
	// Pre-allocate maps with reasonable capacity to avoid resizing
	linkMap := make(map[string][]Link, 2000)
	backlinkMap := make(map[string][]Link, 2000)
	
	// Use a more efficient approach to building the maps
	for _, l := range links {
		// Regular links - store by source
		linkMap[l.Source] = append(linkMap[l.Source], l)
		
		// Backlinks - store by target (only if internal)
		// This creates a reverse index from target to sources
		backlinkMap[l.Target] = append(backlinkMap[l.Target], l)
	}
	
	index.Links = linkMap
	index.Backlinks = backlinkMap
	return index
}
