package domain

type Repository interface {
	// customers
	GetCustomers(limit int, page int) ([]Customer, int, error)

	// BTMUser
	GetBTMUserByAccount(account string) (*BTMUser, error)
}
