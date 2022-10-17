-- name: CreatePost :one
INSERT INTO posts (
  id, 
  title,
  post_description,
  author_name,
  post_date
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY post_date;
