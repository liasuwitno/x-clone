package models

import "time"

// =====================
// USERS
// =====================
type User struct {
	ID         int64     `gorm:"primaryKey"`
	Username   string    `gorm:"type:varchar(30);uniqueIndex;not null"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password   string    `gorm:"type:varchar(255);not null"`
	Bio        *string   `gorm:"type:text"`
	ProfilePic *string   `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (User) TableName() string { return "users" }

// =====================
// TWEETS
// =====================
type Tweet struct {
	ID        int64     `gorm:"primaryKey"`
	Title     *string   `gorm:"type:varchar(255)"`
	Body      *string   `gorm:"type:text"`
	UserID    int64     `gorm:"index;not null"`
	Status    *string   `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Tweet) TableName() string { return "tweets" }

// =====================
// FOLLOWS
// following_user_id follows followed_user_id
// =====================
type Follow struct {
	ID              int64     `gorm:"primaryKey"`
	FollowingUserID int64     `gorm:"index:idx_follows_following,priority:1;not null"`
	FollowedUserID  int64     `gorm:"index:idx_follows_followed,priority:1;not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

// Unique pair (following_user_id, followed_user_id)
func (Follow) TableName() string { return "follows" }

// =====================
// LIKES
// =====================
type Like struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"index;not null"`
	TweetID   int64     `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Like) TableName() string { return "likes" }

// =====================
// EDIT HISTORY
// =====================
type EditHistory struct {
	ID           int64     `gorm:"primaryKey"`
	TweetID      int64     `gorm:"index;not null"`
	PreviousBody *string   `gorm:"type:text"`
	EditedAt     time.Time `gorm:"autoCreateTime"`
}

func (EditHistory) TableName() string { return "edit_history" }
