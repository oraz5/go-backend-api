package pgsql

var GetCategory = `SELECT 
category.id,
category.name,
category.parent,
category.image,
category.icon,
category.create_ts,
category.update_ts,
category.state,
category.version
FROM public."category"
WHERE category.state = 'enabled'`

var GetCategoryById = `SELECT 
category.id,
category.name,
category.parent,
category.image,
category.icon,
category.create_ts,
category.update_ts,
category.state,
category.version
FROM public."category"
WHERE category.id = $1;`

var CreateCategory = `INSERT 
INTO public."category"(name, parent, icon, image, create_ts, update_ts, state, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

var UpdateCategory = `UPDATE public."category"
SET name = $2,
	parent = $3,
	icon = $4,
	image = $5,
	update_ts = $6,
	version = version + 1
WHERE category.id = $1;`

var DeleteCategory = `UPDATE public."category"
SET state = 'deleted'
WHERE category.id = $1;
`
