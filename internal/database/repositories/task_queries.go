package repositories

const create = `
INSERT INTO scheduler (date, title, comment, repeat) 
VALUES (:date, :title, :comment, :repeat);`

const update = `
UPDATE
	scheduler
SET
	date = :date,
	title = :title,
	comment = :comment,
	repeat = :repeat
WHERE
	id = :id;`

const get = `
SELECT 
	id, 
	date, 
	title, 
	comment, 
	repeat 
FROM 
	scheduler 
WHERE 
	id = :id;`

const delete = `
DELETE FROM 
	scheduler 
WHERE 
	id = :id;`
