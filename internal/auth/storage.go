package auth

import (
	"context"
	"errors"
	"time"

	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthStorage struct {
	database *mongo.Collection
}

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists in database")
	ErrOperationTimeout      = errors.New("database operation timed out")
	ErrInternalDatabase      = errors.New("internal database error")
	ErrGenerateUniqueIndexes = errors.New("failure to generate unique indexes for the MongoDB database.")
)

func (u *AuthStorage) capture(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return domain.ErrTimeout

	case errors.Is(err, ErrUserNotFound):
		// case errors.Is(err, mongo.ErrNoDocuments):
		return domain.ErrUserNotFound

	case errors.Is(err, ErrUserAlreadyExists):
		// case mongo.IsDuplicateKeyError(err):
		return domain.ErrUserAlreadyExists

	case errors.Is(err, ErrInternalDatabase):
	case errors.Is(err, ErrGenerateUniqueIndexes):
		return domain.ErrInternal
	default:
		return domain.ErrInternal
	}

	return domain.ErrInternal
}

type UserSchema struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

func NewUserSchema(email, username, password string) *UserSchema {
	return &UserSchema{
		Email:    email,
		Username: username,
		Password: password,
	}
}

func (u *UserSchema) ToDomain() *domain.User {
	return &domain.User{
		ID:       u.ID.Hex(),
		Email:    u.Email,
		Username: u.Username,
		Password: u.Password,
	}
}

func ToSchema(user *domain.User) *UserSchema {
	return &UserSchema{
		ID:       primitive.NewObjectID(),
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
}

func NewAuthStorage(database *mongo.Client) (AuthStorage, error) {
	auth := AuthStorage{
		database: database.Database("gongo").Collection("users"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := auth.database.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return AuthStorage{}, auth.capture(ErrGenerateUniqueIndexes)
	}

	return auth, nil
}

func (u *AuthStorage) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	schema := ToSchema(&user)
	if _, err := u.database.InsertOne(ctx, schema); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, u.capture(ErrUserAlreadyExists)
		}

		return nil, u.capture(err)
	}

	return schema.ToDomain(), nil
}

func (u *AuthStorage) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var schema UserSchema

	if err := u.database.FindOne(ctx, bson.M{"email": email}).Decode(&schema); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, u.capture(ErrUserNotFound)
		}

		return nil, u.capture(err)
	}

	return schema.ToDomain(), nil
}
