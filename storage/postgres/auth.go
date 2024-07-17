package postgres

import (
	"database/sql"
	pb "my_module/generated/auth_service"
	"time"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u *UserRepo) RegisterUser(user *pb.RegisterRequest) (*pb.RegisterResponce, error) {
	resUser := pb.RegisterResponce{}
	err := u.DB.QueryRow(
		`INSERT INTO Users(
			username,
			email,
			password,
			full_name
		)
		VALUES(
			$1,
			$2,
			$3,
			$4)
		RETURNING
			id,
			username,
			email,
			full_name,
			created_at`,
		user.Username, user.Email, user.Password, user.FullName).Scan(
		&resUser.Id, &resUser.Username, &resUser.Email, &resUser.FullName, &resUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &resUser, nil
}

func (u *UserRepo) GetUserByEmail(email string) (*pb.LoginResponce, error) {
	user := pb.LoginResponce{Email: email}

	err := u.DB.QueryRow(`
	SELECT 
		id,
		username,
		password,
		full_name,
		bio,
		countries_visited
	from
		Users
	where
		email = $1
	`, email).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.FullName,
		&user.Bio,
		&user.CountriesVisited,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetProfile(id *pb.GetProfileRequest) (*pb.GetProfileResponce, error) {
	resProfile := pb.GetProfileResponce{}
	err := u.DB.QueryRow(`
	SELECT
		id,
		username,
		email,
		full_name,
		bio,
		countries_visited,
		created_at,
		updated_at
	FROM 
		Users
	WHERE
		id = $1
		`, id.Id).Scan(
		&resProfile.Id,
		&resProfile.Username,
		&resProfile.Email,
		&resProfile.FullName,
		&resProfile.Bio,
		&resProfile.CountriesVisited,
		&resProfile.CreatedAt,
		&resProfile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &resProfile, nil
}

func (u *UserRepo) UpdateProfile(update *pb.UpdateProfileRequest) (*pb.UpdateProfileResponce, error) {
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
	RETURNING
		id,
		username,
		email,
		full_name,
		bio,
		countries_visited,
		updated_at`,
		update.FullName, update.Bio, update.CountriesVisited, time.Now(), update.UserId).Scan(
		&resUpdate.Id,
		&resUpdate.Username,
		&resUpdate.Email,
		&resUpdate.FullName,
		&resUpdate.Bio,
		&resUpdate.CountriesVisited,
		&resUpdate.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &resUpdate, nil
}

func (u *UserRepo) ListProfile(reqList *pb.ListProfileRequest) (*pb.ListProfileResponce, error) {

	rows, err := u.DB.Query(`
	SELECT 
		id,
		username,
		full_name,
		countries_visited
	FROM
		Users
	OFFSET  $1
	LIMIT  $2`,
		(reqList.Page-1)*reqList.Limit, reqList.Limit)

	if err != nil {
		return nil, err
	}
	var users []*pb.User
	for rows.Next() {
		var user pb.User
		err := rows.Scan(&user.Id, &user.Username, &user.FullName, &user.CountriesVisited)

		if err != nil {
			return nil, err
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

	if err != nil {
		return nil, err
	}
	resList := pb.ListProfileResponce{
		Users: users,
		Total: int32(total),
		Page:  reqList.Page,
		Limit: reqList.Limit,
	}

	return &resList, nil
}

func (u *UserRepo) DeleteProfile(id *pb.DeleteProfileRequest) (*pb.DeleteProfileResponce, error) {
	_, err := u.DB.Exec(`
	DELETE FROM
		Users
	WHERE
		id = $1`, id.Id)
	if err != nil {
		return &pb.DeleteProfileResponce{
			Messsage: "Failed deleted to user",
		}, err
	}

	return &pb.DeleteProfileResponce{
		Messsage: "User successfully deleted",
	}, nil
}

func (u *UserRepo) ResetUserPassword(req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponce, error) {

	res := pb.ResetPasswordResponce{}
	err := u.DB.QueryRow(`
	SELECT
		password
	FROM
		Users
	WHERE
		id = $1
		AND
		email = $2
	`, req.UserId, req.Email).Scan(
		&res.Password,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u *UserRepo) ActivityUser(req *pb.ActivityRequest) (*pb.ActivityResponce, error) {
	var countComment int32
	err := u.DB.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		Comments
	WHERE
		author_id = $1
	`, req.UserId).Scan(&countComment)

	if err != nil {
		return nil, err
	}

	var countStories int32
	err = u.DB.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		Stories
	WHERE
		author_id = $1
	`, req.UserId).Scan(&countStories)
	if err != nil {
		return nil, err
	}

	var countLikes int32

	err = u.DB.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		Likes
	WHERE
		user_id = $1
	`, req.UserId).Scan(&countLikes)
	if err != nil {
		return nil, err
	}

	res := pb.ActivityResponce{}

	err = u.DB.QueryRow(`
	SELECT
		updated_at
	FROM
		Users
	WHERE
		id = $1
	`, req.UserId).Scan(&res.LastActive)

	if err != nil {
		return nil, err
	}

	res.StoriesCount = countStories
	res.CommentsCount = countComment
	res.LikesReceived = countLikes

	return &res, nil
}

func (u *UserRepo) Follow(req *pb.FollowRequest) (*pb.FollowResponce, error) {
	res := pb.FollowResponce{}
	err := u.DB.QueryRow(`
	INSERT INTO
		Followers(
			follower_id,
			following_id
		)
		VALUES(
			$1,
			$2	
		)
		Returning
			follower_id,
			following_id,
			followed_at
	`, req.FollowerId, req.FollowingId).Scan(
		&res.FollowerId,
		&res.FollowingId,
		&res.FollowedAt,
	)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (u *UserRepo) FollowersUsers(req *pb.FollowersRequest) (*pb.FollowersResponce, error) {
	rows, err := u.DB.Query(`
	SELECT
		follower_id
	FROM
		Followers
	WHERE
		following_id = $1
	OFFSET $2
	LIMIT $3
	`, req.UserId, (req.Page-1)*req.Limit,req.Limit)

	if err != nil {
		return nil, err
	}

	var followers []*pb.Follower
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		var follower pb.Follower
		err = u.DB.QueryRow(`
		SELECT
			id,
			username,
			full_name
		FROM
			Users
		WHERE
			id = $1
		`, userId).Scan(
			&follower.Id,
			&follower.UserName,
			&follower.FullName,
		)

		if err != nil {
			return nil, err
		}

		followers = append(followers, &follower)
	}
	var total int32
	err = u.DB.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		Followers
	WHERE
		follower_id = $1
	`, req.UserId).Scan(&total)

	if err != nil {
		return nil, err
	}

	return &pb.FollowersResponce{
		Followers: followers,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}, nil
}
