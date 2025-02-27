package models

type PostSettings struct {
	Meaning string `json:"meaning"` //ФИО
	Marking string `json:"marking"`
}
type Paymentstudent struct {
	ID      uint `json:"id"`
	Payment int  `json:"payment"`
}
type PrintLecture struct {
	Lecture_Person_id int    `json:"Lecture_Person_id"`
	Lecture           string `json:"Lecture"`
}
type PostModuls struct {
	Student_id   uint   `json:"Student_id"`
	Module       string `json:"Module"`
	Lock_lecture bool   `json:"lock_lecture"`
}
type PostLecture struct {
	Number  int    `json:"number"`
	Lecture string `json:"lecture"`
}
type PostStudent struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Lecture int    `json:"lecture"`
}
type Resultstruct struct {
	ID           int
	Name_Student string
	Payment      int
	Lecture      string
	Theory       int
	Practice     int
	Tasks        int
}
type PostTelbot struct {
	ModuleAllToggle bool            `json:"moduleAllToggle"`
	Students        []Table_student `json:"students"`
}
