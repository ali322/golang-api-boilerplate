package dto

type QueryUser struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

type UpdateUser struct {
	Email  string `binding:"omitempty,lt=200,email"`
	Avatar string `binding:"omitempty,url"`
	Memo   string `binding:"omitempty"`
}

type RegisterUser struct {
	Username       string `binding:"required,lt=100"`
	Password       string `binding:"required,lt=200"`
	Repeatpassword string `binding:"required,lt=200,eqfield=Password" json:"repeat_password"`
	Email          string `binding:"lt=200,email"`
}

type LoginUser struct {
	UsernameOrEmail string `binding:"required,lt=100" json:"username_or_email"`
	Password        string `binding:"required,lt=200"`
}

type ChangePassword struct {
	OldPassword    string `binding:"required,lt=100" json:"old_password"`
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}

type ResetPassword struct {
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}
