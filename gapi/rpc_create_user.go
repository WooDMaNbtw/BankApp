package gapi

import (
	"context"
	db "github.com/WooDMaNbtw/BankApp/db/sqlc"
	"github.com/WooDMaNbtw/BankApp/pb"
	"github.com/WooDMaNbtw/BankApp/utils"
	"github.com/WooDMaNbtw/BankApp/validators"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// validate request arguments
	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validators.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validators.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validators.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
