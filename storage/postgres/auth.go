package postgres

import (
	"database/sql"
	pb "my_module/generated/auth_service"
	"time"
)

type UserRepo struct{
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo{
	return &UserRepo{DB: db}
}

func (u *UserRepo)RegisterUser(user *pb.RegisterRequest)(*pb.RegisterResponce,error){
	resUser := pb.RegisterResponce{}
	err := u.DB.QueryRow(
		`INSER INTO Users(
		user_name,
		email,
		password,
		full_name
		)
		VALUES(
		$1,
		$2,
		$3,
		$4)
		Returning
		id,
		user_name,
		email,
		full_name,
		created_at`,
	user.Username,user.Email,user.Password,user.FullName).Scan(
		&resUser.Id,&resUser.Username,&resUser.Email,&resUser.FullName,&resUser.CreatedAt,
	)
	if err != nil{
		return nil,err
	}

	return &resUser,nil
}

// func (u *UserRepo) Login(login *pb.LoginRequest)(*pb.LoginResponce,error){
// 	resLogin := pb.LoginResponce{}

// 	u.DB.QueryRow()
// }

func (u *UserRepo) GetProfile(id *pb.GetProfileRequest)(*pb.GetProfileResponce,error){
	resProfile := pb.GetProfileResponce{}
	err := u.DB.QueryRow(`
	SELECT
		id,
		user_name,
		email,
		full_name,
		bio,
		countries_visited,
		created_at,
		updated_at
	FROM 
		User
	WHERE
		id = $1`,id.Id).Scan(
			&resProfile.Id,
			&resProfile.Username,
			&resProfile.Email,
			&resProfile.FullName,
			&resProfile.Bio,
			&resProfile.CountriesVisited,
			&resProfile.CreatedAt,
			&resProfile.UpdatedAt,
		)
	if err != nil{
		return nil,err
	}
	return &resProfile,nil
}

func (u *UserRepo) UpdateProfile(id string,update *pb.UpdateProfileRequest)(*pb.UpdateProfileResponce,error){
	resUpdate := pb.UpdateProfileResponce{}
	err := u.DB.QueryRow(`
	UPDATE
		Users
	SET
		full_name = $1,
		bio = $2,
		countries_visited = $3,
		updated_at = $4
	WHERE
		id = $5
	Returning
		id,
		user_name,
		email,
		full_name,
		bio,
		countries_visited
		updated_at`,
	update.FullName,update.Bio,update.CountriesVisited,time.Now(),id).Scan(
		&resUpdate.Id,
		&resUpdate.Username,
		&resUpdate.Email,
		&resUpdate.FullName,
		&resUpdate.Bio,
		&resUpdate.CountriesVisited,
		&resUpdate.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}
	return &resUpdate,nil
}

func (u *UserRepo) ListProfile(reqList *pb.ListProfileRequest)(*pb.ListProfileResponce,error){

	rows,err := u.DB.Query(`
	SELECT 
		id,
		user_name,
		full_name,
		countries_visited
	FROM
		Users
			OFFSET = $1,
			LIMIT = $2`,
		(reqList.Page - 1)*reqList.Limit,reqList.Limit)

	if err != nil{
		return nil,err
	}
	var users []*pb.User
	for rows.Next(){
		var user pb.User
		err := rows.Scan(&user.Id,&user.Username,&user.FullName,&user.CountriesVisited)

		if err != nil{
			return nil,err
		}
		users = append(users, &user)
	}
	var total int

	err = u.DB.QueryRow(`
	SELECT
		COUNT(*)
	FROM	
		Users
		`).Scan(&total)

	if err != nil{
		return nil,err
	}
	resList := pb.ListProfileResponce{
		Users: users,
		Total: int32(total),
		Page: reqList.Page,
		Limit: reqList.Limit,
	}

	return &resList,nil
}

func (u *UserRepo) DeleteProfile(id *pb.DeleteProfileRequest)(*pb.DeleteProfileResponce,error){
	_,err := u.DB.Exec(`
	DELETE FROM
		Users
	WHERE
		id = $1`,id.Id)
	if err != nil{
		return &pb.DeleteProfileResponce{
			Messsage: "Failed deleted to user",
		},err
	}

	return &pb.DeleteProfileResponce{
		Messsage: "User successfully deleted",
	},nil
}

// func (u *UserRepo) ResetPassword(id string,email *pb.ResetPasswordRequest)(*pb.ResetPasswordResponce,error){
// 	u.DB.Exec(`
// 	Update
// 		Users
// 	SET
// 		Password_hash = $1
// 	WHERE
// 		id = $2`,)
// }

