package db

import (
	"context"
	"database/sql"
	"log"
)

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DeleteTaskRequest struct {
	Name string `json:"name"`
}

type UpdateTaskRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func execWithTx(ctx context.Context, dbp *sql.DB, handler func(*Queries) error) error {
	tx, err := dbp.BeginTx(ctx, nil)
	if err != nil {
		log.Println("db::ExecWithTx(): failed to start tx: ", err)
		return err
	}
	queries := New(dbp)
	txQueries := queries.WithTx(tx)
	err = handler(txQueries)
	if err != nil {
		log.Println("db::ExecWithTx(): failed to exec tx: ", err)
		if terr := tx.Rollback(); terr != nil {
			log.Println("db::ExecWithTx(): failed to rollback tx: ", terr)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Println("db::ExecWithTx(): failed to commit tx: ", err)
	}
	return err
}

func ListTasks(ctx context.Context, dbp *sql.DB) ([]ListTasksRow, error) {
	queries := New(dbp)
	listTasks, err := queries.ListTasks(ctx)
	if err != nil {
		log.Println("db::ListTasks(): Failed to get list task: ", err)
		return nil, err
	}
	return listTasks, nil
}

func CreateTask(ctx context.Context, dbp *sql.DB, createRequest *CreateTaskRequest) error {
	err := execWithTx(ctx, dbp, func(q *Queries) error {
		param := InsertTaskParams{
			Name:        createRequest.Name,
			Title:       createRequest.Title,
			Description: createRequest.Description,
		}
		err := q.InsertTask(ctx, param)
		if err != nil {
			log.Println("db::CreateTask(): Failed to insert task: ", err)
			return err
		}
		status_id, err := q.GetStatusId(ctx, createRequest.Status)
		if err != nil {
			log.Println("db::CreateTask(): Failed to get status id: ", err)
			return err
		}
		err = q.UpdateTaskStatus(ctx, UpdateTaskStatusParams{Name: createRequest.Name, StatusID: status_id})
		if err != nil {
			log.Println("db::CreateTask(): Failed to update status id: ", err)
		}
		return err
	})
	if err != nil {
		log.Println("db::CreateTask(): Failed to create task: ", err)
	}
	return err
}

func DeleteTask(ctx context.Context, dbp *sql.DB, deleteRequest *DeleteTaskRequest) error {
	queries := New(dbp)
	err := queries.DeleteTask(ctx, deleteRequest.Name)
	if err != nil {
		log.Println("db::DeleteTask(): Failed to delete task: ", err)
	}
	return err
}

func UpdateTask(ctx context.Context, dbp *sql.DB, updateRequest *UpdateTaskRequest) error {
	queries := New(dbp)
	err := queries.UpdateTask(ctx, UpdateTaskParams{Name: updateRequest.Name, Name_2: updateRequest.Status})
	if err != nil {
		log.Println("db::UpdateTask(): Failed to update task: ", err)
	}
	return err
}

func ListStatus(ctx context.Context, dbp *sql.DB) ([]string, error) {
	queries := New(dbp)
	listStatus, err := queries.ListStatus(ctx)
	if err != nil {
		log.Println("db::ListStatus(): failed to get list status - ", err)
	}
	return listStatus, err
}
