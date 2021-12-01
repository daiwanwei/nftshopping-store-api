package repositories

var repositoryInstance *repository

func GetRepository() (instance *repository, err error) {
	if repositoryInstance == nil {
		instance, err = newRepository()
		if err != nil {
			return nil, err
		}
		repositoryInstance = instance
	}
	return repositoryInstance, nil
}

type repository struct {
	Auth        AuthDao
	User        UserDao
	Creation    CreationDao
	Transaction TransactionDao
	Brand       BrandDao
	Item        ItemDao
	Collection  CollectionDao
	Stock       StockDao
}

func newRepository() (instance *repository, err error) {
	auth, err := NewAuthDao()
	if err != nil {
		return nil, err
	}
	user, err := NewUserDao()
	if err != nil {
		return nil, err
	}
	creation, err := NewCreationDao()
	if err != nil {
		return nil, err
	}
	transaction, err := NewTransactionDao()
	if err != nil {
		return nil, err
	}
	brand, err := NewBrandDao()
	if err != nil {
		return nil, err
	}
	item, err := NewItemDao()
	if err != nil {
		return nil, err
	}
	stock, err := NewStockDao()
	if err != nil {
		return nil, err
	}
	collection, err := NewCollectionDao()
	if err != nil {
		return nil, err
	}
	return &repository{
		Auth:        auth,
		User:        user,
		Creation:    creation,
		Transaction: transaction,
		Brand:       brand,
		Item:        item,
		Collection:  collection,
		Stock:       stock,
	}, nil
}
