package domain

type FileInfo struct {
	//- compulsory fields
	UserIP         string `json:"userip" bson:"userip"`
	FileName       string `json:"fileName" bson:"fileName"`
	FileSize       int64  `json:"fileSize" bson:"fileSize"`
	FileType       string `json:"fileType" bson:"fileType"`
	FileBucketPath string `json:"filePath" bson:"filePath"`
}
