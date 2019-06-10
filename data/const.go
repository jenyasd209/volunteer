package data

//Constant status for table "Status"
const (
	OrderStatusAvailable = 1
	OrderStatusPerformed = 2
	OrderStatusDone      = 3
)

const (
	UserRoleFreelancer = 1
	UserRoleCustomer   = 2
	UserRoleModerator  = 3
)

const (
	ModeratorFirstName  = "Moderator"
	ModeratorLastName  = "Moderator"
	ModeratorEmail  = "moderator@gmail.com"
	ModeratorPassword  = "moderator"
	ModeratorPhone  = ""
	ModeratorFacebook  = ""
	ModeratorSkype  = ""
	ModeratorAbout  = `This user follows the order on the site. His duties include deleting and editing users,
					   requests, orders. He can also create, update and delete specialties.`
)
