package handler

type ClassNode struct {
	ClassId   int64  `json:"classId"`
	ClassName string `json:"className"`
}
type Resp struct {
	GradeName     string      `json:"gradeName"`
	Male          int64       `json:"male"`
	MaleClasses   []ClassNode `json:"maleClasses"`
	Female        int64       `json:"female"`
	FemaleClasses []ClassNode `json:"femaleClasses"`
}

type Requ struct {
	GradeId int64 `json:"gradeId"`
}
