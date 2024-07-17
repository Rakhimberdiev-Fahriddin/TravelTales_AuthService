package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	pb "my_module/generated/auth_service"
	"my_module/logs"
	"my_module/storage/postgres"
)

type UserService struct {
	pb.UnimplementedAuthServiceServer
	User   postgres.UserRepo
	Logger *slog.Logger
}

func NewUserService(db *sql.DB) *UserService {
	logs.InitLogger()
	return &UserService{
		User: *postgres.NewUserRepo(db),
		Logger: logs.Logger,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponce, error) {
	u.Logger.Info("Register rpc method started")
	res, err := u.User.RegisterUser(req)
	if err != nil {
		u.Logger.Error(err.Error())
		return nil, err
	}
	u.Logger.Info("Register rpc method finished")
	return res, nil
}

func (u *UserService) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponce, error) {
	u.Logger.Info("Login rpc method started")
	res, err := u.User.GetUserByEmail(req.Email)
	if err != nil {
		u.Logger.Error(err.Error())
		return nil, err
	}
	if res.Password != req.Password {
		u.Logger.Error("Password is incorrect")
		return nil, fmt.Errorf("password is incorrect")
	}

	u.Logger.Info("Login rpc method finished")
	return res, nil
}

func (u *UserService) GetUserProfile(ctx context.Context,req *pb.GetProfileRequest)(*pb.GetProfileResponce,error){
	u.Logger.Info("GetProfile rpc method started")
	res, err := u.User.GetProfile(req)
	if err != nil {
		u.Logger.Error(err.Error())
		return nil, err
	}
	u.Logger.Info("GetProfile rpc method finished")
	return res, nil
}

func (u *UserService) UpdateUserProfile(ctx context.Context,req *pb.UpdateProfileRequest)(*pb.UpdateProfileResponce,error){
	u.Logger.Info("UpdateProfile rpc method started")
	res, err := u.User.UpdateProfile(req)
	if err != nil {
		u.Logger.Error(err.Error())
		return nil, err
	}
	u.Logger.Info("UpdateProfile rpc method finished")
	return res, nil
}

func (u *UserService) ListUsersProfile(ctx context.Context,req *pb.ListProfileRequest)(*pb.ListProfileResponce,error){
	u.Logger.Info("List Profile rpc method started")
	res, err := u.User.ListProfile(req)
	if err != nil {
		u.Logger.Error(err.Error())
		return nil, err
	}
	u.Logger.Info("GetUsers rpc method finished")
	return res, nil
}

func (u *UserService) DeleteUserProfile(ctx context.Context,req *pb.DeleteProfileRequest)(*pb.DeleteProfileResponce,error){
	u.Logger.Info("DeleteUser rpc method started")
	res,err := u.User.DeleteProfile(req)
	if err != nil {
		u.Logger.Error(err.Error())
		return res, err
	}
	u.Logger.Info("DeleteUser rpc method finished")
	return res, nil
}

func (u *UserService) ResetUserPassword(ctx context.Context,req *pb.ResetPasswordRequest)(*pb.ResetPasswordResponce,error){
	u.Logger.Info("Reset User Password")
	res,err := u.User.ResetUserPassword(req)
	if err != nil{
		u.Logger.Error("Failed reset user password","error",err)
		return nil,err
	}
	u.Logger.Info("Finished")
	return res,nil
}

func (u *UserService) ActivityUser(ctx context.Context,req *pb.ActivityRequest)(*pb.ActivityResponce,error){
	u.Logger.Info("Activity user")
	res,err := u.User.ActivityUser(req)
	if err != nil{
		return nil,err
	}
	return res,nil
}

func (u *UserService) FollowUser(ctx context.Context,req *pb.FollowRequest)(*pb.FollowResponce,error){
	u.Logger.Info("Follow user")
	res,err := u.User.Follow(req)
	if err != nil{
		u.Logger.Error("Failed user to follow","error",err.Error())
		return nil,err
	}
	u.Logger.Info("Follow rpc method finished")
	return res,nil
}

func (u *UserService) FollowersUsers(ctx context.Context,req *pb.FollowersRequest)(*pb.FollowersResponce,error){
	u.Logger.Info("Follower to user")
	res,err := u.User.FollowersUsers(req)
	if err != nil{
		u.Logger.Error("Failed followers to users gets","error",err.Error())
		return nil,err
	}
	return res,nil
}
