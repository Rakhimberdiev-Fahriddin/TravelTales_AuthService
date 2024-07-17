package postgres

import (
	"fmt"
	"testing"
	"time"

	pb "my_module/generated/auth_service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		FullName: "Test User",
	}

	mock.ExpectQuery(`INSERT INTO Users`).
		WithArgs(req.Username, req.Email, req.Password, req.FullName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "email", "full_name", "created_at"}).
			AddRow("1", req.Username, req.Email, req.FullName, time.Now().Format(time.RFC3339)))

	resp, err := repo.RegisterUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Id)
	assert.Equal(t, req.Username, resp.Username)
}

func TestGetProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := pb.GetProfileRequest{
		Id: "2d194fd1-3fc9-4490-bd7a-4ad794acc207",
	}

	res, err := repo.GetProfile(&req)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.NotEmpty(t,res.Id)
}

func TestUpdateProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.UpdateProfileRequest{
		UserId:           "2d194fd1-3fc9-4490-bd7a-4ad794acc207",
		FullName:         "Updated User",
		Bio:              "Updated Bio",
		CountriesVisited: 35,
	}

	resp, err := repo.UpdateProfile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.FullName, resp.FullName)
}

func TestListProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.ListProfileRequest{
		Page:  1,
		Limit: 10,
	}

	
	resp, err := repo.ListProfile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := pb.DeleteProfileRequest{Id: "2d194fd1-3fc9-4490-bd7a-4ad794acc207"}

	resp,_ := repo.DeleteProfile(&req)

	assert.NotNil(t, resp)
}

func TestResetUserPassword(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.ResetPasswordRequest{
		UserId: "2d194fd1-3fc9-4490-bd7a-4ad794acc207",
		Email:  "test@example.com",
	}

	resp, err := repo.ResetUserPassword(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestActivityUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.ActivityRequest{UserId: "1"}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Comments WHERE author_id = \$1`).
		WithArgs(req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Stories WHERE author_id = \$1`).
		WithArgs(req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Likes WHERE user_id = \$1`).
		WithArgs(req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))

	mock.ExpectQuery(`SELECT updated_at FROM Users WHERE id = \$1`).
		WithArgs(req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(time.Now().Format(time.RFC3339)))

	resp, err := repo.ActivityUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(5), resp.CommentsCount)
	assert.Equal(t, int32(10), resp.StoriesCount)
	assert.Equal(t, int32(15), resp.LikesReceived)
}

func TestFollow(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.FollowRequest{
		FollowerId:  "2d194fd1-3fc9-4490-bd7a-4ad794acc207",
		FollowingId: "ade85737-0bd0-40b1-a11d-6d8515a41879",
	}

	resp, err := repo.Follow(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.FollowerId, resp.FollowerId)
}

func TestFollowersUsers(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Errorf("failed to setup test database: %v", err)
		return
	}
	defer db.Close()

	repo := NewUserRepo(db)

	req := &pb.FollowersRequest{
		UserId: "2d194fd1-3fc9-4490-bd7a-4ad794acc207",
		Page:   1,
		Limit:  10,
	}

	
	resp, err := repo.FollowersUsers(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
