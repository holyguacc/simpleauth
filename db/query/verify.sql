-- name: SetVeriftyStatues :exec
UPDATE  verification SET is_verified = $2 , verefied_on =$3 WHERE username = $1;

-- name: GetVerifyKey :one
SELECT * FROM verification
WHERE username = $1 LIMIT 1;

-- name: SetVerifyCode :exec
INSERT INTO verification (verify_key,username) VALUES ($1,$2);

-- name: GetverificationStatus :one
SELECT is_verified FROM verification
WHERE username =$1 LIMIT 1;

-- name: SetResetKeyAndState :exec
UPDATE verification SET reset_key = $2 , is_reset = $3 WHERE username = $1;


-- name: GetResetKey :one
SELECT * FROM verification
WHERE username = $1 LIMIT 1;


-- name: UpdatePassword :exec
UPDATE users SET hashed_password = $2 , password_changed_at = $3 WHERE username = $1;