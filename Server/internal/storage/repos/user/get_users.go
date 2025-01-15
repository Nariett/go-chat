package user

func (s *store) GetUsers() ([]string, error) {
	var usernames []string
	err := s.db.Select(&usernames, "SELECT name FROM users")
	if err != nil {
		return nil, err
	}

	return usernames, nil
}
