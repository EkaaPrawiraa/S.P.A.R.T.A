type userRepository struct {
	tx *sql.Tx
}

func NewUserRepository(tx *sql.Tx) repository.UserRepository {
	return &userRepository{tx: tx}
}
