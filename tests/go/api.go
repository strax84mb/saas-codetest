package tests

import (
	"encoding/hex"
	"fmt"
)

type UserApi struct {
	// I have changed this from []*User to map.
	// I doubt I will ever encounter such permanent storage, but I have encountered such
	// temporary storage. Map works better here. Multi-level map works best but it seems
	// like too much overengineering for an example this simple.
	storage map[int64]*User
}

/*
   BTW, multi-level map is generally like this:

   type level3Storage map[int]*User

   type level2Storage map[int]level3Storage

   type storage map[int]level2Storage

   Add some getter/setter logic and I can tell you from personal experience
   this is one of the most performant structures in existence
*/

// I chandeg this to return error as well as I want to be able
// to initialize storage without exposing it
func NewUserApi(users []*User) (*UserApi, error) {
	var (
		storage map[int64]*User
		err     error
	)
	if users == nil {
		storage = make(map[int64]*User)
	} else {
		storage, err = convertUserListToMap(users)
		if err != nil {
			return nil, fmt.Errorf("error while initializing storage of new Api: %w", err)
		}
	}
	return &UserApi{
		storage: storage,
	}, nil
}

type User struct {
	Id       string
	Email    string
	FullName string
}

type UpdateUserRequest struct {
	Id       string
	FullName *string
	Email    *string
}

type Error string

func (e Error) Error() string { return string(e) }

var (
	UserNotFound Error = "not_found"
)

// Converts user list into a user map that allows us quicker access
func convertUserListToMap(users []*User) (map[int64]*User, error) {
	var (
		id  int64
		err error
	)
	userMap := make(map[int64]*User, len(users))
	for _, user := range users {
		id, err = decodeKey(user.Id)
		if err != nil {
			return nil, err
		}
		userMap[id] = user
	}
	return userMap, nil
}

// Decided to convert hex code in form of string into int64
// Since I'm changing UserApi.storage to map, I should use
// int64 as key instead of string as it ensures better performance
func decodeKey(strId string) (int64, error) {
	if len(strId)%2 == 1 {
		strId = "0" + strId
	}
	bytes, err := hex.DecodeString(strId)
	if err != nil {
		return 0, fmt.Errorf("failed to decode hex string: %w", err)
	}
	var id int64
	for _, b := range bytes {
		id <<= 8       // shift 8 bits
		id += int64(b) // write byte at the end
	}
	return id, nil
}

func (api *UserApi) Update(request UpdateUserRequest) (*User, error) {
	id, err := decodeKey(request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	// check if user exists
	user, ok := api.storage[id]
	if !ok {
		return nil, UserNotFound
	}
	// overwhite full name if provided
	if request.FullName != nil {
		user.FullName = *request.FullName
	}
	// overwrite email if provided
	if request.Email != nil {
		user.Email = *request.Email
	}
	return user, nil
}
