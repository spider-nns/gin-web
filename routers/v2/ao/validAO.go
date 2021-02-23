package ao

type Login struct {
	UserName string `form:"user" json:"userName"  xml:"userName" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"-"`
}