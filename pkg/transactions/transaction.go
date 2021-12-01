package transactions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"nftshopping-store-api/pkg/databases"
	"nftshopping-store-api/pkg/sagas"
)

type TransactionContext interface {
	mongo.Session
	sagas.Saga
	context.Context
}

type transactionContext struct {
	mongo.Session
	sagas.Saga
	context.Context
}

func NewTransactionContext(ctx context.Context, transaction *Transaction) TransactionContext {
	txnCtx := mongo.NewSessionContext(ctx, transaction.Mongo.Session)
	return &transactionContext{
		Session: transaction.Mongo.Session,
		Saga:    transaction.Saga,
		Context: txnCtx,
	}
}

type Transaction struct {
	Mongo *databases.MongoTransaction
	Saga  sagas.Saga
}

func NewTransaction(name string) (transaction *Transaction, err error) {
	mongo, err := databases.NewMongoTransaction()
	if err != nil {
		return
	}
	transaction = &Transaction{
		Mongo: mongo,
		Saga:  sagas.NewSaga(name),
	}
	return
}

func (t *Transaction) With(ctx context.Context, fn func(TransactionContext) (interface{}, error)) (res interface{}, err error) {
	if err = t.Start(); err != nil {
		return
	}

	res, err = fn(NewTransactionContext(ctx, t))
	if err != nil {
		sagaErr := t.Saga.Abort(ctx)
		if sagaErr != nil {
			return nil, sagaErr
		}
		return nil, err
	}

	if err = t.Commit(ctx); err != nil {
		trxErr := t.Abort(ctx)
		if trxErr != nil {
			return nil, trxErr
		}
		return nil, err
	}
	return
}

func (t *Transaction) Start() (err error) {
	err = t.Mongo.Start()
	return
}

func (t *Transaction) Commit(ctx context.Context) (err error) {
	err = t.Mongo.Commit(ctx)
	return
}

func (t *Transaction) Abort(ctx context.Context) (err error) {
	sessionErr := t.Mongo.Abort(ctx)
	sagaErr := t.Saga.Abort(ctx)
	if sessionErr == nil && sagaErr == nil {
		return nil
	}
	err = &TransactionError{
		SessionErr: sessionErr,
		SagaErr:    sagaErr,
	}
	return
}

func (t *Transaction) End(ctx context.Context) {
	t.Mongo.End(ctx)
}

type TransactionError struct {
	SessionErr error
	SagaErr    error
}

func (e *TransactionError) Error() (strErr string) {
	var session string
	var saga string
	if e.SessionErr != nil {
		session = e.SessionErr.Error()
	}
	if e.SagaErr != nil {
		saga = e.SagaErr.Error()
	}
	return fmt.Sprintf("SessionErr:%s ,SagaErr:%s", session, saga)
}
