-- name: ListTasks :many
select tasks.name, tasks.title, tasks.description, status.name as status
from tasks
inner join status on tasks.status_id = status.id;

-- name: ListStatus :many
select name
from status
order by id asc;

-- name: GetTask :one
select tasks.name as id, tasks.title, tasks.description, status.name as status
from tasks
inner join status on tasks.status_id = status.id
where tasks.name = ?;

-- name: InsertTask :exec
insert into tasks(name, title, description, status_id)
values (?, ?, ?, 1);

-- name: UpdateTaskStatus :exec
update tasks
set status_id = ?
where name = ?;

-- name: DeleteTask :exec
delete from tasks
where name = ?;

-- name: UpdateTask :exec
update tasks, status
set status_id = status.id
where tasks.name = ? and status.name = ?;

-- name: GetStatusId :one
select id
from status
where name = ?;
