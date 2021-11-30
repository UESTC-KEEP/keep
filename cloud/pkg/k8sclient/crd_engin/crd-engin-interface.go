package crd_engin

type CrdEngine interface {
	// CreateCrd 创建keep crd-example
	/*
		输入参数：crdFiles : 所需要进行创建的crd文件的位置
		输出参数：error表示创建crd 的结果
	*/
	CreateCrd(Dir string) error
}
