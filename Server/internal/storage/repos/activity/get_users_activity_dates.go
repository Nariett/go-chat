package activity

import (
	"Server/internal/storage/dbo"
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *store) GetUsersActivityDates(_ *proto.Empty) (*proto.UserActivityDates, error) {
	var activityDates []dbo.UserActivity
	query := `SELECT users.name AS name, date AS date
			  FROM activity 
    		  JOIN users on users.id = activity.idUser`
	err := s.db.Select(&activityDates, query)
	if err != nil {
		return nil, err
	}

	activityMap := make(map[string]*timestamppb.Timestamp)
	for _, item := range activityDates {
		activityMap[item.Name] = timestamppb.New(item.Date)
	}

	userActivityDates := &proto.UserActivityDates{
		ActivityDate: activityMap,
	}

	return userActivityDates, nil
}
