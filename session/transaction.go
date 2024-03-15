package session

import "GolangORM/log"

func (s *Session) Begin() error {
	log.Info("Begin transaction")
	var err error
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error("Begin transaction failed", err)
		return err
	} else {
		return nil
	}
}

func (s *Session) Commit() error {
	log.Info("Commit transaction")
	var err error
	if err = s.tx.Commit(); err != nil {
		log.Error("Commit transaction failed", err)
		return err
	} else {
		return nil
	}
}

func (s *Session) Rollback() error {
	log.Info("Rollback transaction")
	var err error
	if err = s.tx.Rollback(); err != nil {
		log.Error("Rollback transaction failed", err)
		return err
	} else {
		return nil
	}
}
