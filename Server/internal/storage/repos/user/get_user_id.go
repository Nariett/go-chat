package user

func (s *store) GetUserId(name string) (int32, error) {
	var userId int32
	err := s.db.Get(&userId, "SELECT id FROM users WHERE name = $1", name)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
