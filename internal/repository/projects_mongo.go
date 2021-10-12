package repository

import (
	"context"
	"errors"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectsRepo struct {
	db *mongo.Collection
}

func NewProjectsRepo(db *mongo.Database) *ProjectsRepo {
	return &ProjectsRepo{
		db: db.Collection(projectCollection),
	}
}

func (r *ProjectsRepo) GetProjects(ctx context.Context, userId primitive.ObjectID) ([]domain.ProjectMin, error) {
	cursor, err := r.db.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	var projects []domain.ProjectMin
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectsRepo) GetSelfProjects(ctx context.Context, userId primitive.ObjectID) ([]domain.SelfProjectMin, error) {
	cursor, err := r.db.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	var projects []domain.SelfProjectMin
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectsRepo) GetProjectById(ctx context.Context, projectId primitive.ObjectID) (*domain.Project, error) {
	var project *domain.Project
	if err := r.db.FindOne(ctx, bson.M{"_id": projectId}).Decode(&project); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProjectNotFound
		}

		return nil, err
	}
	return project, nil
}

func (r *ProjectsRepo) GetSelfProjectById(ctx context.Context, projectId, userId primitive.ObjectID) (*domain.SelfProject, error) {
	var project *domain.SelfProject
	if err := r.db.FindOne(ctx, bson.M{"_id": projectId, "userId": userId}).Decode(&project); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProjectNotFound
		}

		return nil, err
	}
	return project, nil
}

func (r *ProjectsRepo) GetDrafts(ctx context.Context, userId primitive.ObjectID) ([]domain.SelfProjectMin, error) {
	cursor, err := r.db.Find(ctx, bson.M{"userId": userId, "Published": false})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	var projects []domain.SelfProjectMin
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectsRepo) CreateProject(ctx context.Context, project domain.ProjectInput) error {
	_, err := r.db.InsertOne(ctx, project)
	return err
}

func (r *ProjectsRepo) UpdateProject(ctx context.Context, projectId primitive.ObjectID, project domain.SelfProject) error {
	update := bson.M{}
	if project.Name != "" {
		update["name"] = project.Name
	}
	if project.Description != "" {
		update["description"] = project.Description
	}
	if project.Files != nil {
		update["files"] = project.Files
	}
	if project.Access != "" {
		update["access"] = project.Access
	}
	update["published"] = project.Published
	update["updatedAt"] = project.UpdatedAt

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": projectId}, bson.M{"$set": update})
	return err
}

func (r *ProjectsRepo) RemoveProject(ctx context.Context, projectId primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": projectId})
	return err
}
