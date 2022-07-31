package fileinfo

import (
	"os"
	"sort"
)

//SortByModTime sort by ModTime
type SortByModTime struct {
	fullFileInfos []os.FileInfo
}

func (sortByModTime *SortByModTime) Len() int {
	return len(sortByModTime.fullFileInfos)
}
func (sortByModTime *SortByModTime) Swap(i, j int) {
	sortByModTime.fullFileInfos[i], sortByModTime.fullFileInfos[j] = sortByModTime.fullFileInfos[j], sortByModTime.fullFileInfos[i]
}
func (sortByModTime *SortByModTime) Less(i, j int) bool {
	return sortByModTime.fullFileInfos[i].ModTime().After(sortByModTime.fullFileInfos[j].ModTime())
}

// SortByModifyTime sort file by modify time
func SortByModifyTime(fullFileInfos []os.FileInfo) []os.FileInfo {
	modTime := &SortByModTime{
		fullFileInfos: fullFileInfos,
	}
	sort.Sort(modTime)
	return modTime.fullFileInfos
}