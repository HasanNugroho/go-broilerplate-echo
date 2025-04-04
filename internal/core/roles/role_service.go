package roles

type RoleService struct {
	repo IRoleRepository
}

func NewRoleService(repo IRoleRepository) *RoleService {
	return &RoleService{
		repo: repo,
	}
}
