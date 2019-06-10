package data

import (
	"log"
)

type Moderator struct {
	User
}

func CheckModerator()  {
	if !ExistModerator(){
		CreateModerator()
	}
}

func ExistModerator() (exist bool) {
	err := Db.QueryRow(`SELECT EXISTS(SELECT id FROM users WHERE email = $1)`,
		ModeratorEmail).Scan(&exist)
	if err != nil {
		log.Println(err)
	}
	return
}

func CreateModerator()  {
	user := User{
		FirstName:ModeratorFirstName,
		LastName:ModeratorLastName,
		Email:ModeratorEmail,
		Password:ModeratorPassword,
		Phone:ModeratorPhone,
		Facebook:ModeratorFacebook,
		Skype:ModeratorSkype,
		About:ModeratorAbout,
		RoleID:UserRoleModerator,
	}
	err := user.Create()
	if err != nil{
		log.Println(err)
	}
}

func (moderator *Moderator) CreateSpecialization (specialization *Specialization) (err error){
	err = specialization.Create()
	return
}

func (moderator *Moderator) UpdateCustomer (customer *Customer) (err error){
	err = customer.UpdateInformation()
	return
}

func (moderator *Moderator) UpdateFreelancer (freelancer *Freelancer) (err error){
	err = freelancer.Delete()
	return
}

func (moderator *Moderator) UpdateAvailableOrder (order *Order) (err error){
	err = order.UpdateInformation()
	return
}

func (moderator *Moderator) UpdateSpecialization (specialization *Specialization) (err error){
	err = specialization.Update()
	return
}

func (moderator *Moderator) DeleteCustomer (user *User) (err error){
	err = user.Delete()
	return
}

//func (moderator *Moderator) DeleteFreelancer (freelancer Freelancer){
//	err := freelancer.Delete()
//	if err != nil{
//		log.Println(err)
//	}
//}

func (moderator *Moderator) DeleteAvailableOrder (order *Order) (err error){
	err = order.Delete()
	return
}

func (moderator *Moderator) DeleteRequest (request *FreelancerRequest) (err error){
	err = request.Delete()
	return
}

func (moderator *Moderator) DeleteSpecialization (specialization *Specialization) (err error){
	err = specialization.Delete()
	return
}