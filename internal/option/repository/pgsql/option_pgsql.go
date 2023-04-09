package pgsql

import (
	"context"
	"go-store/internal/entity"

	"go-store/utils/database"
	errorStatus "go-store/utils/errors"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxOptionRepository(pgx *database.PgxAccess) entity.OptionRepository {
	return &PgxAccess{pgx}
}

func (d *PgxAccess) GetSkuValue(ctx context.Context, skuId int) (res []*entity.SkuValue, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetSkuValue"})
	var skuValue []*entity.SkuValue
	rows, err := d.Pool.Query(ctx, GetSkuValueProducts, skuId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetSkuValue - Query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.SkuValue{}
		if err := rows.Scan(&tmp.Id, &tmp.OptionId, &tmp.OptionValueId); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetSkuValue - Scan")
			return nil, err
		}
		skuValue = append(skuValue, tmp)

	}

	return skuValue, nil
}

func (d *PgxAccess) GetOption(ctx context.Context, optionId int) (res *entity.Option, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOption"})
	opt := &entity.Option{}
	row := d.Pool.QueryRow(ctx, GetProductOptions, optionId)
	err = row.Scan(&opt.Id, &opt.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOption - QueryRow")
		return nil, err
	}
	return opt, nil
}

func (d *PgxAccess) GetOptionByCat(ctx context.Context, categoryId int) (res []*entity.Option, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionByCat"})
	opt := []*entity.Option{}
	rows, err := d.Pool.Query(ctx, GetProductOptionsByCat, categoryId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.Option{}
		if err := rows.Scan(&tmp.Id, &tmp.Name); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
			return nil, err
		}
		opt = append(opt, tmp)

	}

	return opt, nil
}

func (d *PgxAccess) GetOptionValue(ctx context.Context, optionValueId int) (res *entity.OptionValue, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionValue"})
	optVal := &entity.OptionValue{}
	row := d.Pool.QueryRow(ctx, GetOptionValues, optionValueId)
	err = row.Scan(&optVal.Id, &optVal.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionValue - QueryRow")
		return nil, err
	}
	return optVal, nil
}

func (d *PgxAccess) GetOptionValueByOptId(ctx context.Context, optionId int) (res []*entity.OptionValueJson, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionValueByOptId"})
	optVal := []*entity.OptionValueJson{}
	rows, err := d.Pool.Query(ctx, GetOptionValuesByOptId, optionId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.OptionValueJson{}
		if err := rows.Scan(&tmp.Id, &tmp.Name); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
			return nil, err
		}
		optVal = append(optVal, tmp)

	}
	return optVal, nil
}

func (d *PgxAccess) GetOptionBySkuValue(ctx context.Context, skuValueId int) (res *entity.OptionJson, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOption"})
	opt := &entity.Option{}
	optValue := &entity.OptionValue{}
	row := d.Pool.QueryRow(ctx, GetOptionBySkuValueId, skuValueId)
	err = row.Scan(&opt.Id, &opt.Name, &optValue.Id, &optValue.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionBySkuValue - QueryRow")
		return nil, err
	}
	option := &entity.OptionJson{
		Id:   opt.Id,
		Name: opt.Name,
	}

	optionValue := []*entity.OptionValueJson{{
		Id:   optValue.Id,
		Name: optValue.Name,
	},
	}
	option.OptionValueJson = optionValue
	return option, nil
}

func (d *PgxAccess) CreateOption(ctx context.Context, option entity.Option) (optionID *int, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOption"})
	row := d.Pool.QueryRow(ctx, CreateOption, option.CategoryId, option.Name, option.CreateTs, option.UpdateTs, option.State, option.Version)
	err = row.Scan(&optionID)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateOption - QueryRow")
		return nil, err
	}
	return optionID, nil
}

func (d *PgxAccess) CreateOptionValue(ctx context.Context, optionValue entity.OptionValue) (optionValueID *int, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOptionValue"})
	row := d.Pool.QueryRow(ctx, CreateOptionValue, optionValue.OptionId, optionValue.Name, optionValue.CreateTs, optionValue.UpdateTs, optionValue.State, optionValue.Version)
	err = row.Scan(&optionValueID)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateOptionValue - QueryRow")
		return nil, err
	}
	return optionValueID, nil
}

func (d *PgxAccess) CreateSkuValue(ctx context.Context, skuID int, optionId int, optionValueId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateSkuValue"})
	_, err := d.Pool.Exec(ctx, CreateSkuValue, skuID, optionId, optionValueId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateSkuValue - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateOption(ctx context.Context, option entity.Option) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOption"})
	_, err := d.Pool.Exec(ctx, UpdateOptionName, option.Id, option.Name, option.CategoryId, option.UpdateTs)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOption - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateOptionValue(ctx context.Context, optionValue entity.OptionValue) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOptionValue"})
	_, err := d.Pool.Exec(ctx, UpdateOptionValueName, optionValue.Id, optionValue.Name, optionValue.OptionId, optionValue.UpdateTs)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOptionValue - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveOption(ctx context.Context, optionId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveOption"})
	_, err := d.Pool.Exec(ctx, RemoveOption, optionId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - RemoveOption - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveOptionValue(ctx context.Context, optionValueId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveOptionValue"})
	_, err := d.Pool.Exec(ctx, RemoveOptionValue, optionValueId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - RemoveOptionValue - Exec")
		return err
	}
	return nil
}
func (d *PgxAccess) RemoveSkuValue(ctx context.Context, skuValueId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveSkuValue"})
	_, err := d.Pool.Exec(ctx, RemoveProductSkuValue, skuValueId)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - RemoveSkuValue - Exec")
		return err
	}
	return nil
}
