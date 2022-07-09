package usecase

import (
	"github.com/diogoalbuquerque/migration-customers/internal/entity"

	"context"
	"errors"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/diogoalbuquerque/migration-customers/pkg/logger"
)

type MigrationUseCase struct {
	sl ServiceLegacy
	sp ServicePerson
	l  logger.Interface
}

func NewMigrationUseCase(sl ServiceLegacy, sp ServicePerson, l logger.Interface) *MigrationUseCase {
	return &MigrationUseCase{
		sl: sl,
		sp: sp,
		l:  l,
	}
}

func (uc *MigrationUseCase) Migration(ctx context.Context, datasourceLimit int) error {

	var seg *xray.Segment

	var lp []entity.LegacyPerson = nil
	var err error = nil

	err = xray.Capture(ctx, "LoadDatasource", func(sc context.Context) error {
		lp, err = uc.sl.LoadDatasource(ctx, datasourceLimit)
		seg = xray.GetSegment(sc)
		seg.AddMetadata("Result", lp)
		seg.AddMetadata("Error", err)
		return err
	})

	if err != nil {
		seg.AddError(err)
		seg.Close(err)
		uc.l.Error(err)
		return err
	}

	s := len(lp)
	uc.l.Info("%d records were found", s)

	if s > 0 {

		var errs []error = nil

		err = xray.Capture(ctx, "Migrate", func(sc context.Context) error {
			lp, errs = uc.sp.Migrate(ctx, lp)
			seg = xray.GetSegment(sc)
			seg.AddMetadata("Result", lp)
			seg.AddMetadata("Error", errs)

			if errs != nil {
				err = errors.New("run - Migrate")
				seg.AddError(err)
				uc.l.Error(errs)
				return err
			}

			return nil
		})

		s = len(lp)
		uc.l.Info("%d records were migrated", s)

		if s > 0 {

			err = xray.Capture(ctx, "Reconciliation", func(sc context.Context) error {
				lp, errs = uc.sl.Reconciliation(ctx, lp)
				seg = xray.GetSegment(sc)
				seg.AddMetadata("Result", lp)
				seg.AddMetadata("Error", errs)

				uc.l.Info("%d records were reconciled", len(lp))

				if errs != nil {
					err = errors.New("run - Reconciliation")
					seg.AddError(err)
					seg.Close(err)
					uc.l.Error(err)
					return err
				}

				return nil
			})

		} else {
			seg.Close(err)
			return err

		}
	}

	return err
}
