package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return domain.ErrUserAlreadyExists
	}

	return err
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) Verify(ctx context.Context, userId primitive.ObjectID, code string) error {
	res, err := r.db.UpdateOne(ctx,
		bson.M{"verification.code": code, "_id": userId},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userId primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"lastVisitAt": time.Now()}})

	return err
}

func (r *UsersRepo) GetById(ctx context.Context, userId primitive.ObjectID) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"_id": userId}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}
	return user, nil
}

func (r *UsersRepo) UpdateById(ctx context.Context, userId primitive.ObjectID, user domain.UserUpdate) error {
	update := bson.M{}
	if user.Name != "" {
		update["name"] = user.Name
	}
	if user.Email != "" {
		update["email"] = user.Email
	}
	if user.Password != "" {
		update["password"] = user.Password
	}
	if user.Role != "" {
		update["role"] = user.Role
	}
	if user.AvatarUrl != "" {
		update["avatarUrl"] = user.AvatarUrl
	}

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": update})
	return err
}

func (r *UsersRepo) RemoveById(ctx context.Context, userId primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": userId})
	return err
}
