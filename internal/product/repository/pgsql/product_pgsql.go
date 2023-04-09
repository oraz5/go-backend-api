package pgsql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	"go-store/utils/database"
	errorStatus "go-store/utils/errors"
)

type PgxAccess struct {
	*database.PgxAccess
}

// NewPgxProductRepository will create an object that represent the product.Repository interface
func NewPgxProductRepository(pgx *database.PgxAccess) entity.ProductRepository {
	return &PgxAccess{pgx}
}

func (d *PgxAccess) GetSku(ctx context.Context, limit int, offset int, filterMap map[string]string) (result []*entity.Sku, products []*entity.Product, err error) {
	var rows pgx.Rows
	dbLog := log.WithFields(log.Fields{"func": "pg.GetSku"})
	baseQuery := d.Builder.
		Select("sku.id",
			"sku.product_id",
			"sku.sku",
			"sku.price",
			"sku.quantity",
			"sku.large_name",
			"sku.small_name",
			"sku.thumb_name",
			"sku.count_viewed",
			"sku.create_ts",
			"sku.update_ts",
			"sku.state",
			"sku.version",
			"product.product_name",
			"product.description",
			"product.category_id",
			"product.create_ts",
			"count(*) OVER() AS total_count").
		From("sku").
		InnerJoin("product ON sku.product_id = product.id").
		Where("product.state = 'enabled'").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	for k, v := range filterMap {
		switch k {
		case "priceStart":
			baseQuery = baseQuery.Where("sku.price >= ?", filterMap["priceStart"])
		case "priceEnd":
			baseQuery = baseQuery.Where("sku.price <= ?", filterMap["priceEnd"])
		case "productName":
			v = fmt.Sprint("%", v, "%")
			baseQuery = baseQuery.Where("(unaccent(product.product_name) ILIKE unaccent(?))", v)
		case "description":
			v = fmt.Sprint("%", v, "%")
			baseQuery = baseQuery.Where("(unaccent(product.description) ILIKE unaccent(?))", v)
		default:
			k = "product." + k
			baseQuery = baseQuery.Where(fmt.Sprintf("%s = ?", k), v)

		}
	}
	query, args, err := baseQuery.ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetSku - r.Builder - query")
		return nil, nil, err
	}
	rows, err = d.Pool.Query(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.Sku{}
		prod := &entity.Product{}
		if err := rows.Scan(&tmp.Id, &tmp.ProductId, &tmp.Sku, &tmp.Price, &tmp.Quantity, &tmp.LargeImage, &tmp.SmallImage, &tmp.ThumbImage, &tmp.CountViewed, &tmp.CreateTs, &tmp.UpdateTs, &tmp.State, &tmp.Version, &prod.ProductName, &prod.Description, &prod.CategoryId, &prod.CreateTs, &tmp.Total); err != nil {
			dbLog.WithFields(log.Fields{"skuId": tmp.Id}).Warning(err)
			return nil, nil, err
		}
		result = append(result, tmp)
		products = append(products, prod)
	}

	return result, products, nil
}

func (d *PgxAccess) GetProducts(ctx context.Context, limit int, offset int, categoryID int) (result []*entity.Product, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetProducts"})
	query, args, err := d.Builder.
		Select("id",
			"product_name",
			"description",
			"category_id",
			"create_ts",
			"count(*) OVER() AS total_count").
		From("product").
		Where("category_id = $1 AND state = 'enabled'", categoryID).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetProducts - r.Builder - query")
		return nil, err
	}
	rows, err := d.Pool.Query(ctx, query, args)
	if err != nil {
		dbLog.Warning(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.Product{}
		err := rows.Scan(&tmp.Id, &tmp.ProductName, &tmp.Description, &tmp.CategoryId, &tmp.CreateTs, &tmp.Total)
		if err != nil {
			dbLog.WithFields(log.Fields{"productId": tmp.Id}).Warning(err)
			return nil, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (d *PgxAccess) GetSkuByProductID(ctx context.Context, productId int) (result []*entity.Sku, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetSkuByProductID"})
	query, args, err := d.Builder.
		Select("id",
			"sku",
			"price",
			"quantity",
			"small_name").
		From("sku").
		Where("product_id = $1 AND state = 'enabled'", productId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetSkuByProductID - r.Builder - query")
		return nil, err
	}
	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		err = fmt.Errorf("pg.GetProductSku: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.Sku{}
		err := rows.Scan(&tmp.Id, &tmp.Sku, &tmp.Price, &tmp.Quantity, &tmp.SmallImage)
		if err != nil {
			dbLog.WithFields(log.Fields{"productId": tmp.Id}).Warning(err)
			return nil, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (d *PgxAccess) GetSingleProduct(ctx context.Context, skuCode string, skuId int) (res *entity.Sku, product *entity.Product, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetSingleProduct"})
	sku := &entity.Sku{}
	prod := &entity.Product{}
	baseQuery := d.Builder.
		Select("sku.id",
			"sku.product_id",
			"sku.sku",
			"sku.price",
			"sku.quantity",
			"sku.large_name",
			"sku.small_name",
			"sku.thumb_name",
			"sku.count_viewed",
			"sku.create_ts",
			"sku.update_ts",
			"sku.state",
			"sku.version",
			"product.product_name",
			"product.description",
			"product.category_id",
			"product.create_ts").
		From("sku").
		InnerJoin("product ON sku.product_id = product.id")
	if len(skuCode) > 0 {
		baseQuery = baseQuery.Where("sku.sku = $1 AND  product.state = 'enabled'", skuCode)
	} else if skuId > 0 {
		baseQuery = baseQuery.Where("sku.id = $1 AND  product.state = 'enabled'", skuId)
	}
	query, args, err := baseQuery.ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetSingleProduct - r.Builder - query")
		return nil, nil, err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&sku.Id, &sku.ProductId, &sku.Sku, &sku.Price, &sku.Quantity, &sku.LargeImage, &sku.SmallImage, &sku.ThumbImage, &sku.CountViewed, &sku.CreateTs, &sku.UpdateTs, &sku.State, &sku.Version, &prod.ProductName, &prod.Description, &prod.CategoryId, &prod.CreateTs)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return nil, nil, err
	}
	if err != nil {
		dbLog.Warning(err)
		return nil, nil, err
	}
	return sku, prod, nil
}

func (d *PgxAccess) CreateProduct(ctx context.Context, prod *entity.Product) (prodId *int, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateProduct"})
	query, args, err := d.Builder.
		Insert("product").
		Columns("id",
			"user_id",
			"address",
			"phone",
			"comment",
			"status",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(prod.ProductName,
			prod.Description,
			prod.CategoryId,
			prod.BrandId,
			prod.RegionId,
			prod.CreateTs,
			prod.UpdateTs,
			prod.State,
			prod.Version).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - CreateProduct - r.Builder - query")
		return nil, err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&prodId)
	if err != nil {
		dbLog.Warning(err)
		return nil, err
	}
	return prodId, nil
}

func (d *PgxAccess) UpdateProduct(ctx context.Context, prod *entity.Product) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateProduct"})
	query, args, err := d.Builder.
		Update("product").
		SetMap(map[string]interface{}{
			"category_id":  prod.CategoryId,
			"product_name": prod.ProductName,
			"description":  prod.Description,
			"brand_id":     prod.BrandId,
			"region_id":    prod.RegionId,
			"update_ts":    prod.UpdateTs,
			"state":        prod.State,
			"version":      "version + 1"}).
		Where("product.id = $1", prod.Id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - UpdateProduct - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveProduct(ctx context.Context, productId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveProduct"})
	query, args, err := d.Builder.
		Update("product").
		SetMap(map[string]interface{}{
			"update_ts": time.Now(),
			"state":     "deleted",
			"version":   "version + 1"}).
		Where("product.id = $1", productId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - RemoveProduct - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}

func (d *PgxAccess) CreateSku(ctx context.Context, sku *entity.Sku) (err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateSku"})
	query, args, err := d.Builder.
		Insert("sku").
		Columns("product_id",
			"sku",
			"price",
			"quantity",
			"large_name",
			"small_name",
			"thumb_name",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(sku.ProductId,
			sku.Sku,
			sku.Price,
			sku.Quantity,
			sku.LargeImage,
			sku.SmallImage,
			sku.ThumbImage,
			sku.CreateTs,
			sku.UpdateTs,
			sku.State,
			sku.Version).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - CreateSku - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateSku(ctx context.Context, sku *entity.Sku) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateSku"})
	query, args, err := d.Builder.
		Update("sku").
		SetMap(map[string]interface{}{
			"price":      sku.Price,
			"quantity":   sku.Quantity,
			"large_name": sku.LargeImage,
			"small_name": sku.SmallImage,
			"thumb_name": sku.ThumbImage,
			"create_ts":  "create_ts",
			"update_ts":  sku.UpdateTs,
			"state":      sku.State,
			"version":    "version + 1"}).
		Where("sku.id = $1", sku.Id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - UpdateSku - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveSku(ctx context.Context, skuId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveSku"})
	query, args, err := d.Builder.
		Update("sku").
		SetMap(map[string]interface{}{
			"state":   "deleted",
			"version": "version + 1"}).
		Where("sku.id = $1", skuId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - RemoveSku - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}
