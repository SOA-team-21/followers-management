package handler

import (
	"context"
	"fmt"
	"net/http"

	"followers.xws.com/model"
	follower "followers.xws.com/proto/followers"
	"followers.xws.com/service"
)

type FollowersHandler struct {
	follower.UnimplementedFollowersServiceServer
	FollowersService *service.PersonService
}

func (handler *FollowersHandler) GetProfile(ctx context.Context, request *follower.UserIdRequset) (*follower.PersonResponse, error) {

	userId := fmt.Sprint(request.UserId)

	var profile, err = handler.FollowersService.GetProfile(userId)
	if err != nil {
		return &follower.PersonResponse{}, err
	}
	return personToRpc(profile), nil
}

func (handler *FollowersHandler) GetRecommended(ctx context.Context, request *follower.UserIdRequset) (*follower.FollowersResponse, error) {

	userId := fmt.Sprint(request.UserId)

	var recommended, err = handler.FollowersService.GetRecommended(userId)
	if err != nil {
		return &follower.FollowersResponse{}, err
	}
	return followersToRpc(&recommended), nil
}

func (handler *FollowersHandler) GetFollowing(ctx context.Context, request *follower.UserIdRequset) (*follower.FollowersResponse, error) {

	userId := fmt.Sprint(request.UserId)

	var userFollowers, err = handler.FollowersService.GetFollowing(userId)
	if err != nil {
		return &follower.FollowersResponse{}, err
	}
	return followersToRpc(&userFollowers), nil
}

func (handler *FollowersHandler) GetFollowers(ctx context.Context, request *follower.UserIdRequset) (*follower.FollowersResponse, error) {

	userId := fmt.Sprint(request.UserId)

	var userFollowers, err = handler.FollowersService.GetFollowers(userId)
	if err != nil {
		return &follower.FollowersResponse{}, err
	}
	return followersToRpc(&userFollowers), nil
}

func (handler *FollowersHandler) IsFollowing(ctx context.Context, request *follower.TwoUserIdRequest) (*follower.StatusCodeResponse, error) {
	userId := fmt.Sprint(request.UserId1)
	isFollowingUserId := fmt.Sprint(request.UserId2)

	var isFollowing, err = handler.FollowersService.IsFollowing(userId, isFollowingUserId)
	if err != nil {
		return &follower.StatusCodeResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}, err
	}
	return &follower.StatusCodeResponse{
		StatusCode: http.StatusOK,
		Message:    "Following status: " + fmt.Sprint(isFollowing),
	}, nil
}

func (handler *FollowersHandler) Unfollow(ctx context.Context, request *follower.TwoUserIdRequest) (*follower.StatusCodeResponse, error) {
	userFollower := fmt.Sprint(request.UserId1)
	toUnfollow := fmt.Sprint(request.UserId2)

	var err = handler.FollowersService.Unfollow(toUnfollow, userFollower)
	if err != nil {
		return &follower.StatusCodeResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}, err
	}
	return &follower.StatusCodeResponse{
		StatusCode: http.StatusOK,
		Message:    "Unfollowing successful",
	}, nil
}

func (handler *FollowersHandler) Follow(ctx context.Context, request *follower.TwoUserIdRequest) (*follower.StatusCodeResponse, error) {
	userFollower := fmt.Sprint(request.UserId1)
	toFollow := fmt.Sprint(request.UserId2)

	var err = handler.FollowersService.Follow(toFollow, userFollower)
	if err != nil {
		return &follower.StatusCodeResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}, err
	}
	return &follower.StatusCodeResponse{
		StatusCode: http.StatusOK,
		Message:    "Following successful",
	}, nil
}

func personToRpc(person *model.Person) *follower.PersonResponse {
	return &follower.PersonResponse{
		Id:      person.Id,
		UserId:  fmt.Sprint(person.UserId),
		Name:    person.Name,
		Surname: person.Surname,
		Picture: person.Picture,
		Bio:     person.Bio,
		Quote:   person.Quote,
		Email:   person.Email,
	}
}

func followerToRpc(person *model.Follower) *follower.Follower {
	return &follower.Follower{
		UserId:  person.UserId,
		Name:    person.Name,
		Surname: person.Surname,
		Quote:   person.Quote,
		Email:   person.Email,
	}
}

func followersToRpc(people *model.Followers) *follower.FollowersResponse {
	result := make([]*follower.Follower, len(*people))
	for i, e := range *people {
		result[i] = followerToRpc(e)
	}
	return &follower.FollowersResponse{Followers: result}
}
