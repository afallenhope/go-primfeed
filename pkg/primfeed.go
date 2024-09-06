package primfeed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Notification struct {
	Type          string            `json:"type"`
	GroupID       interface{}       `json:"groupId,omitempty"`
	CreatedAt     string            `json:"createdAt"`
	Notifications []SubNotification `json:"notifications"`
}

type SubNotification struct {
	Origin Origin `json:"origin"`
	ID     string `json:"id"`
	Read   bool   `json:"read"`
}

type Origin struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Handle             string `json:"handle"`
	Picture            string `json:"picture"`
	PictureUuid        string `json:"pictureUuid"`
	ProfilePictureUuid string `json:"profilePictureUuid"`
	ProfileMedia       string `json:"profileMedia"`
	BannerMedia        string `json:"bannerMedia,omitempty"`
	Verified           bool   `json:"verified"`
	Type               string `json:"type"`
	IsUser             bool   `json:"isUser"`
}

type NotificationsResponse struct {
	UnreadCount   int            `json:"unreadCount"`
	Notifications []Notification `json:"notifications"`
}

type User struct {
	ID           string     `json:"id"`
	Picture      string     `json:"picture"`
	ProfileMedia any        `json:"profileMedia"`
	BannerMedia  string     `json:"bannerMedia"`
	Name         string     `json:"name"`
	About        string     `json:"about"`
	Handle       string     `json:"handle"`
	IsUser       bool       `json:"isUser"`
	Verified     bool       `json:"verified"`
	Title        any        `json:"title,omitempty"`
	Socials      Socials    `json:"socials"`
	Registered   Registered `json:"registered"`
}

type UserProfile struct {
	User
	ShowFollowButton bool `json:"showFollowButton"`
	Followers        int  `json:"followers"`
	Follow           int  `json:"follow"`
	CanFollow        bool `json:"canFollow"`
}

type Follower struct {
	User
	Owner User `json:"owner,omitempty"`
	Rules any  `json:"rules,omitempty"`
}

type Followers []Follower

type Socials struct {
	XURL          any `json:"xUrl"`
	DeviantArtURL any `json:"deviantArtUrl"`
	BlueskyURL    any `json:"blueskyUrl"`
	InstagramURL  any `json:"instagramUrl"`
	FacebookURL   any `json:"facebookUrl"`
	FlickrURL     any `json:"flickrUrl"`
	PersonalURL   any `json:"personalUrl"`
}

type Registered struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type Subscription struct {
	Type                string `json:"type"`
	MaximumMbUploadSize int    `json:"maximumMbUploadSize"`
}

type Profile struct {
	Version         string       `json:"version"`
	AvailableStores []any        `json:"availableStores"`
	Subscription    Subscription `json:"subscription"`
	SelectedStore   any          `json:"selectedStore"`
	SelectedEntity  User         `json:"selectedEntity"`
	User            User         `json:"user"`
	Token           string       `json:"token"`
	Perms           any          `json:"canAddProducts"`
}

type MyProfile struct {
	About string `json:"about"`
	Socials
}

type Perms struct {
	CanDelete bool `json:"canDelete,omitempty"`
	CanEdit   bool `json:"canEdit,omitempty"`
	CanReport bool `json:"canReport,omitempty"`
}

type Media struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	Height  int    `json:"height,omitempty"`
	Width   int    `json:"width,omitempty"`
	Version int    `json:"version,omitempty"`
}

type Feed struct {
	CommentsCount int   `json:"commentsCount,omitempty"`
	Likes         int   `json:"likes,omitempty"`
	Liked         bool  `json:"liked,omitempty"`
	Perms         Perms `json:"perms,omitempty"`
	Data          struct {
		ID            string  `json:"id,omitempty"`
		Owner         User    `json:"owner,omitempty"`
		QuotedPost    any     `json:"quotedPost,omitempty"`
		CreatedAt     int     `json:"createdAt,omitempty"`
		UpdatedAt     any     `json:"updatedAt,omitempty"`
		Content       string  `json:"content,omitempty"`
		Rating        string  `json:"rating,omitempty"`
		IsAi          bool    `json:"isAi,omitempty"`
		IsRender      bool    `json:"isRender,omitempty"`
		PublicGallery bool    `json:"publicGallery,omitempty"`
		Media         []Media `json:"media,omitempty"`
	} `json:"data,omitempty"`
}

type FeedResponse struct {
	Feed []Feed `json:"feed"`
}

type Primfeed struct {
	Token   string
	BaseURL string
	Me      struct {
		Profile       Profile
		Notifications NotificationsResponse
		Followers     Followers
		Following     Followers
	}
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CompanyID any    `json:"companyID"`
	Redirect  string `json:"redirect"`
}

type LoginInworldRequest struct {
	Username string `json:"username"`
}

type LoginInworldResponse struct {
	RequestID string `json:"requestId"`
}

type LoginInworldCodeRequest struct {
	RequestID string `json:"requestId"`
	Username  string `json:"username"`
	OTP       string `json:"otp"`
	CompanyID string `json:"companyID"`
	Redirect  string `json:"redirect"`
}

type LoginResponse struct {
	User             string `json:"user,omitempty"`
	SelectedStore    string `json:"selectedStore,omitempty"`
	Token            string `json:"token,omitempty"`
	Redirect         string `json:"redirect,omitempty"`
	ProfilePictureID string `json:"profilePictureUuid,omitempty"`
	Error            string `json:"error,omitempty"`
}

const (
	APIURL string = "api.primfeed.com"
	URL    string = "www.primfeed.com"
)

func NewPrimfeed(baseUrl string) *Primfeed {
	if !strings.HasPrefix(baseUrl, "http") {
		baseUrl = fmt.Sprintf("https://%s", baseUrl)
	}

	return &Primfeed{
		BaseURL: fmt.Sprintf("%s", baseUrl),
	}
}

func (p *Primfeed) SetToken(token string) {
	p.Token = token
}

func (p *Primfeed) Login(username string, password string, company any) (LoginResponse, error) {
	loginRequest := LoginRequest{
		Username:  username,
		Password:  password,
		CompanyID: company,
		Redirect:  "/",
	}

	var loginResponse LoginResponse
	url := fmt.Sprintf("%s/login", p.BaseURL)

	err := p.Request("POST", url, loginRequest, nil, &loginResponse)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("login failed: %v\n\n", err)
	}

	p.SetToken(loginResponse.Token)

	return loginResponse, nil

}

func (p *Primfeed) GetLoginCode(username string) (LoginInworldResponse, error) {

	var loginInworldResponse LoginInworldResponse
	url := fmt.Sprintf("%s/login/create-inworld-request", p.BaseURL)

	loginInWorldRequest := LoginInworldRequest{
		Username: username,
	}

	err := p.Request("POST", url, loginInWorldRequest, nil, &loginInworldResponse)
	if err != nil {
		return LoginInworldResponse{}, fmt.Errorf("could not send inworld request: %v", err)
	}

	// This is hacky.
	p.Me.Profile.User.Handle = username

	return loginInworldResponse, nil
}

func (p *Primfeed) LoginWithCode(requestId string, code string, company string) (LoginResponse, error) {

	var loginResponse LoginResponse

	loginInworldCode := LoginInworldCodeRequest{
		RequestID: requestId,
		Username:  p.Me.Profile.User.Handle,
		OTP:       code,
		CompanyID: company,
		Redirect:  "/",
	}

	url := fmt.Sprintf("%s/login/inworld-code", p.BaseURL)

	err := p.Request("POST", url, loginInworldCode, nil, &loginResponse)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("could not login with code: %v", err)
	}

	p.SetToken(loginResponse.Token)

	return loginResponse, nil
}

func (p *Primfeed) Request(method string, path string, data interface{}, headers map[string]string, target interface{}) error {
	var body *bytes.Buffer

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(jsonData)
	} else {
		body = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respBody) == 0 {
		return nil
	}

	if target != nil {
		err := json.Unmarshal(respBody, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Primfeed) GetUserFollowers(username string) (Followers, error) {
	url := fmt.Sprintf("%s/entity/%s/followers", p.BaseURL, username)
	var followers Followers

	if err := p.Request("GET", url, nil, nil, &followers); err != nil {
		return nil, err
	}

	return followers, nil
}

func (p *Primfeed) GetUserFollows(username string) (Followers, error) {
	url := fmt.Sprintf("%s/entity/%s/followed", p.BaseURL, username)
	var following Followers

	if err := p.Request("GET", url, nil, nil, &following); err != nil {
		return nil, err
	}

	return following, nil
}

func (p *Primfeed) IsFollowingUser(username string, user string) (bool, error) {
	var followers Followers

	followers, err := p.GetUserFollows(username)

	if err != nil {
		return false, err
	}

	is_following := false

	for _, f := range followers {
		if strings.ToLower(f.Handle) == strings.ToLower(user) {
			is_following = true
			break
		}
	}

	return is_following, nil
}

func (p *Primfeed) GetMe() error {
	var profile Profile
	url := fmt.Sprintf("%s/me", p.BaseURL)

	err := p.Request("GET", url, nil, nil, &profile)

	if err != nil {
		return fmt.Errorf("could not get profile %v", err)
	}

	followers, err := p.GetUserFollowers(profile.User.Handle)
	if err != nil {
		return fmt.Errorf("could not get followers %v", err)
	}

	follows, err := p.GetUserFollows(profile.User.Handle)
	if err != nil {
		return fmt.Errorf("could not get follows %v", err)
	}

	p.Me.Profile = profile
	p.Me.Followers = followers
	p.Me.Following = follows

	return nil
}

func (p *Primfeed) GetUserProfile(username string) (UserProfile, error) {
	var profile UserProfile
	url := fmt.Sprintf("%s/entity/%s", p.BaseURL, username)

	err := p.Request("GET", url, nil, nil, &profile)

	if err != nil {
		return profile, fmt.Errorf("could not get profile %v", err)
	}

	return profile, nil
}

func (p *Primfeed) FollowUser(username string) error {
	profile, err := p.GetUserProfile(username)
	if err != nil {
		return fmt.Errorf("error getting profile to follow: %v", err)
	}

	return p.FollowById(profile.ID)
}

func (p *Primfeed) FollowById(id string) error {
	url := fmt.Sprintf("%s/follow/%s", p.BaseURL, id)

	err := p.Request("POST", url, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error following user: %v", err)
	}

	return nil
}

func (p *Primfeed) UnfollowUser(username string) error {
	profile, err := p.GetUserProfile(username)
	if err != nil {
		return fmt.Errorf("error unfollowing user: %v", err)
	}

	return p.UnfollowById(profile.ID)
}

func (p *Primfeed) UnfollowById(id string) error {
	url := fmt.Sprintf("%s/follow/%s", p.BaseURL, id)

	err := p.Request("DELETE", url, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Primfeed) UpdateProfile(profile interface{}) error {
	url := fmt.Sprintf("%s/entity/%s", p.BaseURL, p.Me.Profile.User.Handle)

	err := p.Request("PATCH", url, profile, nil, nil)
	if err != nil {
		return fmt.Errorf("could not update profile: %v", err)
	}

	return nil
}

func (p *Primfeed) GetNotifications() (NotificationsResponse, error) {
	var notificationResponse NotificationsResponse
	url := fmt.Sprintf("%s/notifications", p.BaseURL)

	err := p.Request("GET", url, nil, nil, &notificationResponse)
	if err != nil {
		return NotificationsResponse{}, fmt.Errorf("error could not get notifications: %v", err)
	}

	p.Me.Notifications = notificationResponse
	return notificationResponse, nil
}

func (p *Primfeed) GetNotificationCount() (int, error) {
	url := fmt.Sprintf("%s/notifications/count", p.BaseURL)

	var count int
	err := p.Request("GET", url, nil, nil, &count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *Primfeed) Like(post string) error {
	url := fmt.Sprintf("%s/pf/post/%s/like", p.BaseURL, post)

	err := p.Request("POST", url, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("could not like post: %v", err)
	}

	return nil
}

func (p *Primfeed) UnLike(post string) error {
	return p.Like(post)
}

func (p *Primfeed) GetFeed(id string, page int) (FeedResponse, error) {
	url := fmt.Sprintf("%s/pf/%s/feed?page=%d", p.BaseURL, id, page)

	var feedResponse FeedResponse

	err := p.Request("GET", url, nil, nil, &feedResponse)
	if err != nil {
		return FeedResponse{}, fmt.Errorf("could not load feed: %v", err)
	}

	return feedResponse, nil
}
