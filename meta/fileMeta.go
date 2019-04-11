package meta

//FileMeta：文件元信息结构
type FileMeta struct {
	FileSha1   string //文件唯一标识
	FileName   string //文件名
	FileSize   int64  //大小
	Path       string //路径
	UpdateTime string //上传时间
}

var fileMetaMap map[string]*FileMeta

//初始化
func init() {
	fileMetaMap = make(map[string]*FileMeta)
}

//新增\更新文件元信息
func UpdateFileMeta(meta *FileMeta) {
	fileMetaMap[meta.FileSha1] = meta
}

//删除元文件信息
func DeleteFileMetaMap(meta *FileMeta)  {
	delete(fileMetaMap, meta.FileSha1)
}

//根据FileShal获取FileMeta
func GetFileMeta(fileSha1 string) *FileMeta {
	return fileMetaMap[fileSha1]
}

//获取fileMetas切片
func GetFileMetaSplic() (fileMetaSplic []*FileMeta) {
	fileMetaSplic = make([]*FileMeta, 0)
	for _, val := range fileMetaMap {
		if val == nil {
			continue
		}
		fileMetaSplic = append(fileMetaSplic, val)
	}
	return
}
