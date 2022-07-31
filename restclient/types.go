/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package restclient

// PatchType these are constants to support HTTP PATCH utilized by
// both the client and server that didn't make sense for a whole package to be
// dedicated to.
type PatchType string

const (
	JSONPatchType           PatchType = "application/json-patch+json"
	MergePatchType          PatchType = "application/merge-patch+json"
	StrategicMergePatchType PatchType = "application/strategic-merge-patch+json"
)

//Object 用来抽象rest接口的不同类型的数据结构
type Object interface {
	//GetObjectKind return what kind with this object
	GetObjectKind() string
}

//ObjectImpl 嵌入此结构，实现Object接口
type ObjectImpl struct {
}

//GetObjectKind get object kind
func (bo *ObjectImpl) GetObjectKind() string {
	return ""
}
